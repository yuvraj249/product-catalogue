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

func GetClaims3(c *gin.Context) (jwt.MapClaims, bool) {
	claimsExist, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claims, ok := claimsExist.(jwt.MapClaims)
	return claims, ok

}

func NullStringVal3(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func CtxTimeout3(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 4*time.Second)
}

func CreateCategory(c *gin.Context) {
	claims, ok := GetClaims3(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role, _ := claims["role"].(string)
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system amdin is allowed to create categories"})
		c.Abort()
		return
	}

	var catInput models.Category
	if err := c.BindJSON(&catInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}

	catInput.CategoryName = strings.TrimSpace(catInput.CategoryName)
	catInput.CategoryDescription = strings.TrimSpace(catInput.CategoryDescription)

	if err := utils.CategoryValidate(catInput.CategoryName, catInput.CategoryDescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while validating input"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout3(c)
	defer cancel()

	query := "insert into categories(category_id, category_name, category_description) values(?, ?, ?)"
	result, err := config.DB.ExecContext(ctx, query, catInput.CategoryID, catInput.CategoryName, NullStringVal3(catInput.CategoryDescription))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		c.Abort()
		return

	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"message": "product created", "product_id": id})

}

func GetCategory(c *gin.Context) {
	claims, ok := GetClaims3(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role, _ := claims["role"].(string)
	ctx, cancel := CtxTimeout3(c)
	defer cancel()
	var rows *sql.Rows
	var err error

	if role == "system_admin" {
		rows, err = config.DB.QueryContext(ctx, "select * from categories")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching categories"})
			c.Abort()
			return
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}
	defer rows.Close()
	categories := []models.Category{}
	for rows.Next() {
		var ct models.Category
		var catDesp sql.NullString

		if err := rows.Scan(&ct.CategoryID, &ct.CategoryName, &catDesp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning categories"})
			c.Abort()
			return
		}
		if catDesp.Valid {
			ct.CategoryDescription = catDesp.String
		} else {
			ct.CategoryDescription = ""
		}
		categories = append(categories, ct)

	}

	c.JSON(http.StatusOK, gin.H{"message": categories})

}

func GetCategoryByID(c *gin.Context) {
	c_id := c.Param("id")
	id, err := strconv.Atoi(c_id)
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
	if role == "system_admin" || role == "supplier_admin" {
		row = config.DB.QueryRowContext(ctx, "select * from categories where category_id=?", c_id)

	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch category"})
		c.Abort()
		return
	}

	var ct models.Category
	var catDesp sql.NullString

	if err := row.Scan(&ct.CategoryID, &ct.CategoryName, &catDesp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		c.Abort()
		return
	}

	if catDesp.Valid {
		ct.CategoryDescription = catDesp.String
	} else {
		ct.CategoryDescription = ""
	}

	c.JSON(http.StatusOK, gin.H{"category": ct})

}

func DeleteCategory(c *gin.Context) {
	c_id := c.Param("id")
	id, err := strconv.Atoi(c_id)
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
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can delete category"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout3(c)
	defer cancel()

	result, err := config.DB.ExecContext(ctx, "delete from categories where category_id= ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		c.Abort()
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})

}

func UpdateCategory(c *gin.Context) {
	c_id := c.Param("id")
	id, err := strconv.Atoi(c_id)
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
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can delete category"})
		c.Abort()
		return
	}

	var catInput struct {
		CategoryName        *string `json:"category_name"`
		CategoryDescription *string `json:"category_description"`
	}
	if err := c.BindJSON(&catInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if catInput.CategoryName == nil && catInput.CategoryDescription == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout3(c)
	defer cancel()

	var exists models.Category
	var catDesp sql.NullString
	row := config.DB.QueryRowContext(ctx, "select * from categories where category_id = ?")
	if err := row.Scan(&exists.CategoryID, &exists.CategoryName, &catDesp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusFound, gin.H{"error": "category not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch category"})
		c.Abort()
		return
	}
	if catDesp.Valid {
		exists.CategoryDescription = catDesp.String
	} else {
		exists.CategoryDescription = ""
	}

	newName := exists.CategoryName
	newDesp := exists.CategoryDescription

	if catInput.CategoryName != nil {
		newName = strings.TrimSpace(*catInput.CategoryName)
	}
	if catInput.CategoryDescription != nil {
		newDesp = strings.TrimSpace(*catInput.CategoryDescription)
	}

	if err := utils.CategoryValidate(newName, newDesp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please enter valid category details"})
		c.Abort()
		return
	}

	var nameArg interface{} = nil
	var despArg interface{} = nil

	if catInput.CategoryName != nil {
		nameArg = newName
	}
	if catInput.CategoryDescription != nil {
		despArg = newDesp
	}

	query := "update set category_id = ifnull(?, category_name), category_description = ifnull(?, category_description) where category_id = ?"
	_, err = config.DB.ExecContext(ctx, query, nameArg, despArg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating category"})
		c.Abort()
		return
	}

	var catUpdated models.Category
	var catDesp2 sql.NullString
	err = config.DB.QueryRowContext(ctx, "select * from categories where category_id= ?", id).Scan(&catUpdated.CategoryID, &catUpdated.CategoryName, &catDesp2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "category updated", "category_id": id, "error": "category updated but failed to fetch row"})
		c.Abort()
		return
	}

	if catDesp2.Valid {
		catUpdated.CategoryName = catDesp2.String
	} else {
		catUpdated.CategoryDescription = ""
	}
	c.JSON(http.StatusOK, gin.H{"message": "category updated", "category": catUpdated})

}
