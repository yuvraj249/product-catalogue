package functions

import (
	"database/sql"
	"net/http"
	"product-catalogue/config"
	"product-catalogue/models"
	"product-catalogue/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	_, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin is allowed to create categories"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	query := "insert into categories( category_name, category_description) values(?, ?)"
	result, err := config.DB.ExecContext(ctx, query, catInput.CategoryName, NullStringVal(catInput.CategoryDescription))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create category"})
		c.Abort()
		return

	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "category created", "category_id": id})

}

func GetCategory(c *gin.Context) {
	_, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role := c.GetString("role")
	ctx, cancel := CtxTimeout(c)
	defer cancel()
	if role != "system_admin" && role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}
	rows, err := config.DB.QueryContext(ctx, "select category_id, category_name, category_description from categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching categories"})
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

	c.JSON(http.StatusOK, gin.H{"categories": categories})

}

func GetCategoryByID(c *gin.Context) {
	c_id := c.Param("id")
	id, err := strconv.Atoi(c_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		c.Abort()
		return
	}
	_, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role := c.GetString("role")

	if role != "system_admin" && role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	row := config.DB.QueryRowContext(ctx, "select category_id, category_name, category_description from categories where category_id = ?", id)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		c.Abort()
		return
	}
	_, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can delete category"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
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
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})

}

func UpdateCategory(c *gin.Context) {
	c_id := c.Param("id")
	id, err := strconv.Atoi(c_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		c.Abort()
		return
	}
	_, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can update category"})
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

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var exists models.Category
	var catDesp sql.NullString
	row := config.DB.QueryRowContext(ctx, "select category_id, category_name, category_description from categories where category_id = ?", id)
	if err := row.Scan(&exists.CategoryID, &exists.CategoryName, &catDesp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "update categories set category_name = ifnull(?, category_name), category_description = ifnull(?, category_description) where category_id = ?"
	_, err = config.DB.ExecContext(ctx, query, NullStringPtr(&newName), NullStringPtr(&newDesp), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating category"})
		c.Abort()
		return
	}

	var catUpdated models.Category
	var catDesp2 sql.NullString
	err = config.DB.QueryRowContext(ctx, "select category_id, category_name, category_description from categories where category_id= ?", id).Scan(&catUpdated.CategoryID, &catUpdated.CategoryName, &catDesp2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "category updated", "category_id": id, "error": "category updated but failed to fetch row"})
		c.Abort()
		return
	}

	if catDesp2.Valid {
		catUpdated.CategoryDescription = catDesp2.String
	} else {
		catUpdated.CategoryDescription = ""
	}
	c.JSON(http.StatusOK, gin.H{"message": "category updated", "category": catUpdated})

}
