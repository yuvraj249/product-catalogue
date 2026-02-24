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

	"github.com/gin-gonic/gin"
)

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
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

func CreateSupplier(c *gin.Context) {
	role := c.GetString("role")
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
	supplier.Company = strings.TrimSpace(supplier.Company)

	if err := utils.SupplierValidate(supplier.Name, supplier.ContactInfo, supplier.Email, supplier.Company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var count int
	if err := config.DB.QueryRowContext(ctx, "select count(*) from suppliers where lower(email) = lower(?)", supplier.Email).Scan(&count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while reading suppliers"})
		c.Abort()
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier with this email already exists"})
		c.Abort()
		return
	}

	result, err := config.DB.ExecContext(ctx, "insert into suppliers(name, contact_info, email, company) values(?,?,?,?)", supplier.Name, NullStringVal(supplier.ContactInfo), supplier.Email, supplier.Company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot not create supplier"})
		c.Abort()
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "supplier created", "supplier_id": id})

}

func GetSupplier(c *gin.Context) {
	role := c.GetString("role")
	search := strings.TrimSpace(c.Query("q"))

	supplierID := "%" + search + "%"
	contactInfo := "%" + search + "%"
	name := "%" + search + "%"
	email := "%" + search + "%"
	company := "%" + search + "%"
	ctx, cancel := CtxTimeout(c)
	defer cancel()
	var rows *sql.Rows
	var err error

	switch role {
	case "system_admin":
		rows, err = config.DB.QueryContext(
			ctx, "SELECT supplier_id, name, contact_info, email, company FROM suppliers WHERE ? = '' OR CAST(supplier_id AS CHAR) LIKE ? OR contact_info LIKE ? OR LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(company) LIKE LOWER(?) ORDER BY supplier_id ASC",
			search, supplierID, contactInfo, name, email, company,
		)

	case "supplier_admin":
		supplierID := c.GetInt("supplier_id")
		rows, err = config.DB.QueryContext(
			ctx,
			"SELECT supplier_id, name, contact_info, email, company FROM suppliers WHERE supplier_id = ? AND (? = '' OR contact_info LIKE ? OR LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(company) LIKE LOWER(?))",
			search, supplierID, contactInfo, name, email, company,
		)

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}
	defer rows.Close()

	suppliers := []models.Supplier{}
	for rows.Next() {
		var sp models.Supplier
		var contact_info sql.NullString
		if err = rows.Scan(&sp.SupplierID, &sp.Name, &contact_info, &sp.Email, &sp.Company); err != nil {
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
	role := c.GetString("role")
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	if role != "system_admin" {
		supplierValue := c.GetInt("supplier_id")
		var loggedCompany string
		err := config.DB.QueryRowContext(ctx, "select company from suppliers where supplier_id= ?", supplierValue).Scan(&loggedCompany)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch logged in supplier company"})
			c.Abort()
			return
		}
		var getCompany string
		err = config.DB.QueryRowContext(ctx, "select company from suppliers where supplier_id= ?", id).Scan(&getCompany)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch requested company"})
			c.Abort()
			return
		}

		if loggedCompany != getCompany {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied - supplier belongs to a different company"})
			c.Abort()
			return
		}

	}
	row := config.DB.QueryRowContext(ctx, "select supplier_id, name, contact_info, email, company from suppliers where supplier_id = ?", id)
	var sp models.Supplier
	var contact_info sql.NullString
	if err := row.Scan(&sp.SupplierID, &sp.Name, &contact_info, &sp.Email, &sp.Company); err != nil {
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
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can delete supplier"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
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
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "suppplier deleted"})

}

func UpdateSupplier(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		c.Abort()
		return
	}
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system admin can update supplier"})
		c.Abort()
		return
	}
	var supplierInput struct {
		Name        string `json:"name,omitempty"`
		ContactInfo string `json:"contact_info,omitempty"`
		Email       string `json:"email,omitempty"`
		Company     string `json:"company,omitempty"`
	}

	if err := c.BindJSON(&supplierInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if supplierInput.Name == "" && supplierInput.ContactInfo == "" && supplierInput.Email == "" && supplierInput.Company == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
		c.Abort()
		return
	}
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var exists models.Supplier
	var contactInfo sql.NullString
	row := config.DB.QueryRowContext(ctx, "select  supplier_id, name, contact_info, email, company from suppliers where supplier_id= ?", id)
	if err := row.Scan(&exists.SupplierID, &exists.Name, &contactInfo, &exists.Email, &exists.Company); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch supplier"})
		c.Abort()
		return
	}
	if contactInfo.Valid {
		exists.ContactInfo = contactInfo.String

	} else {
		exists.ContactInfo = ""
	}

	newName := exists.Name
	newContact := exists.ContactInfo
	newEmail := exists.Email
	newCompany := exists.Company

	if supplierInput.Name != "" {
		newName = strings.TrimSpace(supplierInput.Name)
	}
	if supplierInput.ContactInfo != "" {
		newContact = strings.TrimSpace(supplierInput.ContactInfo)
	}
	if supplierInput.Email != "" {
		newEmail = strings.TrimSpace(supplierInput.Email)
	}
	if supplierInput.Company != "" {
		newCompany = strings.TrimSpace(supplierInput.Company)
	}

	nameChanged := false
	contactChanged := false
	emailChanged := false
	companyChanged := false

	if supplierInput.Name != "" {
		if normalizeString(newName) != normalizeString(exists.Name) {
			nameChanged = true
		}
	}
	if supplierInput.ContactInfo != "" {
		if normalizeString(newContact) != normalizeString(exists.ContactInfo) {
			contactChanged = true
		}
	}
	if supplierInput.Email != "" {
		if normalizeEmail(newEmail) != normalizeEmail(exists.Email) {
			emailChanged = true
		}
	}
	if supplierInput.Company != "" {
		if normalizeString(newCompany) != normalizeString(exists.Company) {
			companyChanged = true
		}
	}
	if !(nameChanged || contactChanged || emailChanged || companyChanged) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no changes provided"})
		c.Abort()
		return
	}
	if err := utils.SupplierValidate(newName, newContact, newEmail, newCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if emailChanged {
		var count int
		query := "select count(*) from suppliers where lower(email) = ? and supplier_id <> ?"
		if err := config.DB.QueryRowContext(ctx, query, normalizeEmail(newEmail), id).Scan(&count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking email"})
			c.Abort()
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already used by another supplier"})
			c.Abort()
			return
		}
	}

	var nameArg interface{} = nil
	var contactArg interface{} = nil
	var emailArg interface{} = nil
	var companyArg interface{} = nil

	if nameChanged {
		nameArg = newName
	}
	if contactChanged {
		contactArg = newContact
	}
	if emailChanged {
		emailArg = newEmail
	}
	if companyChanged {
		companyArg = newCompany
	}

	query := "update suppliers set name = ifnull(?, name), contact_info = ifnull(?, contact_info), email = ifnull(?, email), company= ifnull(?, company) where supplier_id = ?"
	_, err = config.DB.ExecContext(ctx, query, nameArg, contactArg, emailArg, companyArg, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating supplier"})
		c.Abort()
		return
	}

	var supplierUpdated models.Supplier
	var contactUpated sql.NullString
	err = config.DB.QueryRowContext(ctx, "select  supplier_id, name, contact_info, email, company from suppliers where supplier_id= ?", id).Scan(&supplierUpdated.SupplierID, &supplierUpdated.Name, &contactUpated, &supplierUpdated.Email, &supplierUpdated.Company)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "supplier updated", "supplier_id": id, "error": "supplier updated but failed to fetch row"})
		c.Abort()
		return
	}

	if contactUpated.Valid {
		supplierUpdated.ContactInfo = contactUpated.String
	} else {
		supplierUpdated.ContactInfo = ""

	}
	c.JSON(http.StatusOK, gin.H{"message": "supplier updated", "supplier": supplierUpdated})

}
