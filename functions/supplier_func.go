package functions

import (
	"context"
	"database/sql"
	"net/http"
	"product-catalogue/config"
	"product-catalogue/models"
	"product-catalogue/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SupplierExists(ctx context.Context, supplierID int) (bool, error) {
	var exist bool
	err := config.DB.QueryRowContext(ctx, "select exists(select 1 from suppliers where supplier_id=?)", supplierID).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err

	}
	return exist, nil

}

func GetClaims2(c *gin.Context) (jwt.MapClaims, bool) {
	claimsExist, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claims, ok := claimsExist.(jwt.MapClaims)
	return claims, ok

}

func CtxTimeout2(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 4*time.Second)
}

func NullStringVal2(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func CreateSupplier(c *gin.Context) {
	claims, ok := GetClaims2(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role, _ := claims["role"].(string)
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system amdin is allowed to create suppliers"})
		c.Abort()
		return
	}

	var supplier models.Supplier
	if err := c.BindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}

	supplier.Name = strings.TrimSpace(supplier.Name)
	supplier.ContactInfo = strings.TrimSpace(supplier.ContactInfo)
	supplier.Email = strings.TrimSpace(supplier.Email)
	supplier.Comapany = strings.TrimSpace(supplier.Comapany)

	if err := utils.SupplierValidate(supplier.Name, supplier.ContactInfo, supplier.Email, supplier.Comapany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while validating input"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout2(c)
	defer cancel()

	var exists bool
	if err := config.DB.QueryRowContext(ctx, "select exists(select 1 from suppliers where email =?)", supplier.Email).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while reading suppliers"})
		c.Abort()
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier with this email already exists"})
		c.Abort()
		return
	}

	result, err := config.DB.ExecContext(ctx, "insert into suppliers(name, contact_info, email, company) values(?,?,?,?)", supplier.Name, NullStringVal2(supplier.ContactInfo), supplier.Email, supplier.Comapany)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot not create supplier"})
		c.Abort()
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusBadRequest, gin.H{"error": "supplier created", "supplier_id": id})

}

func GetSupplier(c *gin.Context) {
	claims, ok := GetClaims2(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role, _ := claims["role"].(string)
	ctx, cancel := CtxTimeout2(c)
	defer cancel()
	var rows *sql.Rows
	var err error

	if role == "system_admin" {
		rows, err = config.DB.QueryContext(ctx, "select * from suppliers")
	} else {
		supplierValue, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier id missing in token"})
			c.Abort()
			return
		}
		supplierId := int(supplierValue.(float64))
		rows, err = config.DB.QueryContext(ctx, "select * from suppliers where supplier_id=?", supplierId)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch suppliers"})
		c.Abort()
		return
	}
	defer rows.Close()

	suppliers := []models.Supplier{}
	for rows.Next() {
		var sp models.Supplier
		var contact_info sql.NullString
		if err := rows.Scan(&sp.SupplierID, &contact_info, &sp.Email, &sp.Comapany); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning supppliers"})
			c.Abort()
			return
		}
		if contact_info.Valid {
			sp.ContactInfo = contact_info.String

		} else {
			sp.ContactInfo = ""
		}
		suppliers = append(suppliers, sp)
	}

	c.JSON(http.StatusOK, gin.H{"suppliers": suppliers})

}

func GetSupplierByID(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		c.Abort()
		return

	}

	claims, ok := GetClaims2(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication reuqired"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	ctx, cancel := CtxTimeout2(c)
	defer cancel()

	if role != "system_admin" {
		supplierValue, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier id missing in token"})
			c.Abort()
			return
		}
		if int(supplierValue.(float64)) != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

	}
	row := config.DB.QueryRowContext(ctx, "select * from suppliers where supplier_id=?", id)
	var sp models.Supplier
	var contact_info sql.NullString
	if err := row.Scan(&sp.SupplierID, &contact_info, &sp.Email, &sp.Comapany); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
			c.Abort()
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch supplier"})
		c.Abort()
		return
	}
	if contact_info.Valid {
		sp.ContactInfo = contact_info.String
	} else {
		sp.ContactInfo = ""
	}

	c.JSON(http.StatusOK, gin.H{"supplier": sp})

}

func DeleteSupplier(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		c.Abort()
		return

	}

	claims, ok := GetClaims2(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication reuqired"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can delete supplier"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout2(c)
	defer cancel()

	result, err := config.DB.ExecContext(ctx, "delete from suppliers where supplier_id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete supplier"})
		c.Abort()
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "suppplier deleted"})

}
