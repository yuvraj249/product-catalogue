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

func GetClaims(c *gin.Context) (jwt.MapClaims, bool) {
	claimsExist, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claims, ok := claimsExist.(jwt.MapClaims)
	return claims, ok

}

func NullStringVal(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func NullStringPtr(p *string) interface{} {
	if p == nil {
		return nil
	}
	v := strings.TrimSpace(*p)
	if v == "" {
		return nil

	}
	return v
}

func NullFlaotPtr(p *float64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func NullIntPtr(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func CtxTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 4*time.Second)
}

func CreateProduct(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier amdin is allowed to create product"})
		c.Abort()
		return
	}
	supplierValue, ok := claims["supplier_id"]
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token !!"})
		c.Abort()
		return
	}

	supplierID := int(supplierValue.(float64))
	var in struct {
		ProductName        string   `json:"product_name"`
		ProductDescription *string  `json:"product_description"`
		ProductCost        float64  `json:"product_cost"`
		ProductCategoryID  *int     `json:"product_category_id"`
		DiscountType       *string  `json:"discount_type"`
		DiscountValue      *float64 `json:"discount_value"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	var product models.Product
	product.ProductName = strings.TrimSpace(in.ProductName)
	if in.ProductDescription != nil {
		product.ProductDescription = strings.TrimSpace(*in.ProductDescription)
	} else {
		product.ProductDescription = ""
	}
	product.ProductCost = in.ProductCost
	product.ProductCategoryID = in.ProductCategoryID
	product.ProductSupplierID = &supplierID
	product.DiscountType = in.DiscountType
	product.DiscountValue = in.DiscountValue

	if err := utils.ProductValidate(product.ProductName, product.ProductDescription, product.ProductCost, product.ProductCategoryID, product.ProductSupplierID, product.DiscountType, product.DiscountValue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()
	if product.ProductCategoryID != nil {
		ok, err := utils.CategoryExists(*product.ProductCategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking category"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category id does not exist"})
			c.Abort()
			return
		}

	}

	query := "insert into products(product_name,product_description,product_cost,product_category_id, product_supplier_id,discount_type,discount_value) values(?,?,?,?,?,?,?)"
	result, err := config.DB.ExecContext(ctx, query, product.ProductName, NullStringVal(product.ProductDescription), product.ProductCost, NullIntPtr(product.ProductCategoryID), NullIntPtr(product.ProductSupplierID), NullStringPtr(product.DiscountType), NullFlaotPtr(product.DiscountValue))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		c.Abort()
		return

	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"message": "product created ", "product_id": id})

}

func GetProduct(c *gin.Context) {
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var rows *sql.Rows
	var err error

	switch role {
	case "system_admin":
		rows, err = config.DB.QueryContext(ctx, "select product_id, product_name, product_description, product_cost, product_category_id, product_supplier_id, discount_type, discount_value from products")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not fetch products"})
			c.Abort()
			return
		}
	case "supplier_admin":
		supplierVal, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
			c.Abort()
			return
		}
		suppID := int(supplierVal.(float64))
		var company string
		if err = config.DB.QueryRowContext(ctx, "select company from suppliers WHERE supplier_id = ?", suppID).Scan(&company); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusForbidden, gin.H{"error": "supplier not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch supplier company"})
			c.Abort()
			return
		}
		rows, err = config.DB.QueryContext(ctx, "select p.product_id, p.product_name, p.product_description, p.product_cost, p.product_category_id, p.product_supplier_id, p.discount_type,p.discount_value from products p join suppliers s on p.product_supplier_id = s.supplier_id where s.company = ? order by product_id desc", company)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch company products"})
			c.Abort()
			return
		}

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		var desp sql.NullString
		var cost sql.NullFloat64
		var catgID sql.NullInt64
		var suppId sql.NullInt64
		var discType sql.NullString
		var discValue sql.NullFloat64

		if err := rows.Scan(&p.ProductID, &p.ProductName, &desp, &cost, &catgID, &suppId, &discType, &discValue); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning products"})
			c.Abort()
			return
		}

		if cost.Valid {
			p.ProductCost = cost.Float64
		} else {
			p.ProductCost = 0
		}

		if desp.Valid {
			p.ProductDescription = desp.String
		} else {
			p.ProductDescription = ""
		}
		if catgID.Valid {
			temp := int(catgID.Int64)
			p.ProductCategoryID = &temp
		} else {
			p.ProductCategoryID = nil
		}
		if suppId.Valid {
			temp := int(suppId.Int64)
			p.ProductSupplierID = &temp
		} else {
			p.ProductSupplierID = nil
		}
		if discType.Valid {
			temp := discType.String
			p.DiscountType = &temp
		} else {
			p.DiscountType = nil
		}
		if discValue.Valid {
			temp := float64(discValue.Float64)
			p.DiscountValue = &temp
		} else {
			p.DiscountValue = nil
		}

		products = append(products, p)

	}
	c.JSON(http.StatusOK, gin.H{"products": products})

}

func GetProductByID(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		c.Abort()
		return
	}
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var row *sql.Row

	switch role {
	case "system_admin":
		row = config.DB.QueryRowContext(ctx, "select product_id, product_name, product_description, product_cost, product_category_id, product_supplier_id, discount_type, discount_value from products where product_id=?", id)
	case "supplier_admin":
		supplierVal, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
			c.Abort()
			return
		}
		suppID := int(supplierVal.(float64))
		row = config.DB.QueryRowContext(ctx, "select p.product_id, p.product_name, p.product_description, p.product_cost, p.product_category_id, p.product_supplier_id, p.discount_type,p.discount_value from products p join suppliers s on p.product_supplier_id = s.supplier_id where p.product_id = ? and s.company = (select company from suppliers where supplier_id = ?)", id, suppID)

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return

	}

	var p models.Product
	var desp sql.NullString
	var cost sql.NullFloat64
	var catgID sql.NullInt64
	var suppId sql.NullInt64
	var discType sql.NullString
	var discValue sql.NullFloat64

	if err := row.Scan(&p.ProductID, &p.ProductName, &desp, &cost, &catgID, &suppId, &discType, &discValue); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		c.Abort()
		return

	}

	if cost.Valid {
		p.ProductCost = cost.Float64
	} else {
		p.ProductCost = 0
	}

	if desp.Valid {
		p.ProductDescription = desp.String
	} else {
		p.ProductDescription = ""
	}
	if catgID.Valid {
		temp := int(catgID.Int64)
		p.ProductCategoryID = &temp

	} else {
		p.ProductCategoryID = nil
	}

	if suppId.Valid {
		temp := int(suppId.Int64)
		p.ProductSupplierID = &temp

	} else {
		p.ProductSupplierID = nil
	}

	if discType.Valid {
		temp := discType.String
		p.DiscountType = &temp
	} else {
		p.DiscountType = nil
	}
	if discValue.Valid {
		temp := float64(discValue.Float64)
		p.DiscountValue = &temp
	} else {
		p.DiscountValue = nil
	}

	c.JSON(http.StatusOK, gin.H{"product": p})

}

func UpdateProduct(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		c.Abort()
		return

	}

	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier admin can update product"})
		c.Abort()
		return

	}

	supplierValue, ok := claims["supplier_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "supplier id missing in token"})
		c.Abort()
		return
	}
	suppID := int(supplierValue.(float64))
	var product struct {
		ProductName        *string  `json:"product_name"`
		ProductDescription *string  `json:"product_description"`
		ProductCost        *float64 `json:"product_cost"`
		ProductCategoryID  *int     `json:"product_category_id"`
		ProductSupplierID  *int     `json:"product_supplier_id"`
		DiscountType       *string  `json:"discount_type"`
		DiscountValue      *float64 `json:"discount_value"`
	}
	err = c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if product.ProductSupplierID != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot change product supplier"})
		return
	}
	if product.ProductName == nil && product.ProductDescription == nil && product.ProductCost == nil && product.ProductCategoryID == nil && product.ProductSupplierID == nil && product.DiscountType == nil && product.DiscountValue == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided for update"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var exist models.Product
	var desp sql.NullString
	var catgID sql.NullInt64
	var supppID sql.NullInt64
	var discType sql.NullString
	var discValue sql.NullFloat64

	row := config.DB.QueryRowContext(ctx, "select product_id , product_name, product_description, product_cost, product_category_id, product_supplier_id, discount_type, discount_value from products where product_id = ?", id)
	if err := row.Scan(&exist.ProductID, &exist.ProductName, &desp, &exist.ProductCost, &catgID, &supppID, &discType, &discValue); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		c.Abort()
		return
	}

	if desp.Valid {
		exist.ProductDescription = desp.String
	} else {
		exist.ProductDescription = ""
	}
	if catgID.Valid {
		tmp := int(catgID.Int64)
		exist.ProductCategoryID = &tmp
	} else {
		exist.ProductCategoryID = nil
	}
	if supppID.Valid {
		tmp := int(supppID.Int64)
		exist.ProductSupplierID = &tmp
	} else {
		exist.ProductSupplierID = nil
	}
	if discType.Valid {
		tmp := discType.String
		exist.DiscountType = &tmp
	} else {
		exist.DiscountType = nil
	}
	if discValue.Valid {
		tmp := float64(discValue.Float64)
		exist.DiscountValue = &tmp
	} else {
		exist.DiscountValue = nil
	}
	if exist.ProductSupplierID == nil || *exist.ProductSupplierID != suppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own products"})
		return
	}

	newName := exist.ProductName
	newDesc := exist.ProductDescription
	newCost := exist.ProductCost
	newCat := exist.ProductCategoryID
	newDiscType := exist.DiscountType
	newDiscVal := exist.DiscountValue

	if product.ProductName != nil {
		trimmed := strings.TrimSpace(*product.ProductName)
		newName = trimmed
	}
	if product.ProductDescription != nil {
		trimmed := strings.TrimSpace(*product.ProductDescription)
		newDesc = trimmed
	}
	if product.ProductCost != nil {
		newCost = *product.ProductCost
	}
	if product.ProductCategoryID != nil {
		newCat = product.ProductCategoryID
	}
	if product.DiscountType != nil {
		trimmed := strings.TrimSpace(*product.DiscountType)
		newDiscType = &trimmed
	}
	if product.DiscountValue != nil {
		newDiscVal = product.DiscountValue
	}

	if err := utils.ProductValidate(newName, newDesc, newCost, newCat, exist.ProductSupplierID, newDiscType, newDiscVal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if newCat != nil {
		ok, err := utils.CategoryExists(*newCat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking category"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category id does not exist"})
			c.Abort()
			return
		}
	}

	_, err = config.DB.ExecContext(ctx, "update products set product_name = ifnull(?, product_name), product_description = ifnull(?, product_description), product_cost = ifnull(?, product_cost), product_category_id = ifnull(?, product_category_id), discount_type = ifnull(?, discount_type), discount_value = ifnull(?, discount_value) where product_id = ?", NullStringPtr(&newName), NullStringPtr(&newDesc), newCost, NullIntPtr(newCat), NullStringPtr(newDiscType), NullFlaotPtr(newDiscVal), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating product"})
		c.Abort()
		return
	}

	row2 := config.DB.QueryRowContext(ctx, "select product_id , product_name, product_description, product_cost, product_category_id, product_supplier_id, discount_type, discount_value from products where product_id = ?", id)
	var updated models.Product
	var desp2 sql.NullString
	var catID2 sql.NullInt64
	var suppID2 sql.NullInt64
	var discT sql.NullString
	var discV sql.NullFloat64

	if err := row2.Scan(&updated.ProductID, &updated.ProductName, &desp2, &updated.ProductCost, &catID2, &suppID2, &discT, &discV); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "product updated", "product_id": id, "warning": "updated but failed to fetch row"})
		c.Abort()
		return
	}

	if desp2.Valid {
		updated.ProductDescription = desp2.String
	} else {
		updated.ProductDescription = ""
	}
	if catID2.Valid {
		tmp := int(catID2.Int64)
		updated.ProductCategoryID = &tmp
	} else {
		updated.ProductCategoryID = nil
	}
	if suppID2.Valid {
		tmp := int(suppID2.Int64)
		updated.ProductSupplierID = &tmp
	} else {
		updated.ProductSupplierID = nil
	}
	if discT.Valid {
		tmp := discT.String
		updated.DiscountType = &tmp
	} else {
		updated.DiscountType = nil
	}
	if discV.Valid {
		tmp := float64(discV.Float64)
		updated.DiscountValue = &tmp
	} else {
		updated.DiscountValue = nil
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated", "product": updated})

}

func DeleteProduct(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		c.Abort()
		return
	}
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role, _ := claims["role"].(string)
	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier_admin can delete products"})
		c.Abort()
		return
	}

	supplierVal, ok := claims["supplier_id"]
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
		c.Abort()
		return
	}

	suppID := int(supplierVal.(float64))
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var prodSuppID sql.NullInt64
	err = config.DB.QueryRowContext(ctx, "SELECT product_supplier_id FROM products WHERE product_id = ?", id).Scan(&prodSuppID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		c.Abort()
		return
	}

	if !prodSuppID.Valid || int(prodSuppID.Int64) != suppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own products"})
		c.Abort()
		return
	}

	result, err := config.DB.ExecContext(ctx, "delete from products where product_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		c.Abort()
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})

}
