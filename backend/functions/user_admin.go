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

func CreateSuppAdmin(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can create supplier admin users"})
		c.Abort()
		return
	}

	var admin struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		SupplierID int    `json:"supplier_id"`
	}

	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}

	admin.Name = strings.TrimSpace(admin.Name)
	admin.Email = strings.TrimSpace(admin.Email)

	if admin.Name == "" || admin.Email == "" || admin.Password == "" || admin.SupplierID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, password and valid supplier_id are required"})
		c.Abort()
		return
	}

	if errEmail := utils.IsValidEmail(admin.Email); errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errEmail.Error()})
		c.Abort()
		return
	}

	if errPwd := utils.IsValidPassword(admin.Password); errPwd != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errPwd.Error()})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	exists, err := SupplierExists(ctx, admin.SupplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking supplier"})
		c.Abort()
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supplier_id does not exist"})
		c.Abort()
		return
	}

	var emailExists bool
	err = config.DB.QueryRowContext(ctx, "select exists(select 1 from users where lower(email)=?)", strings.ToLower(admin.Email)).Scan(&emailExists)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking existing user"})
		c.Abort()
		return
	}

	if emailExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		c.Abort()
		return
	}

	hashed, err := utils.HashPwd(admin.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		c.Abort()
		return
	}
	query := "insert into users(name, email, password_hash, role, supplier_id) values(?,?,?,'supplier_admin',?)"
	result, err := config.DB.ExecContext(ctx, query, admin.Name, admin.Email, hashed, admin.SupplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create supplier_admin user"})
		c.Abort()
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"message":     "supplier_admin user created successfully",
		"user_id":     id,
		"name":        admin.Name,
		"email":       admin.Email,
		"supplier_id": admin.SupplierID,
	})

}

func GetsuppAdmin(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can list users"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	rows, err := config.DB.QueryContext(ctx, "select user_id, name, email, role, supplier_id from users order by user_id asc")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching users"})
		c.Abort()
		return
	}
	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var u models.User
		var supp sql.NullInt64
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Role, &supp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning users"})
			c.Abort()
			return
		}
		if supp.Valid {
			v := int(supp.Int64)
			u.SupplierID = &v
		} else {
			u.SupplierID = nil
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetsuppAdminByID(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can view user details"})
		c.Abort()
		return
	}

	u_id := c.Param("id")
	userID, err := strconv.Atoi(u_id)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var user models.User
	var supplierID sql.NullInt64
	err = config.DB.QueryRowContext(ctx, "SELECT user_id, name, email, role, supplier_id FROM users WHERE user_id = ?", userID).Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &supplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching user"})
		c.Abort()
		return
	}
	if supplierID.Valid {
		s := int(supplierID.Int64)
		user.SupplierID = &s
	} else {
		user.SupplierID = nil
	}

	resp := gin.H{
		"user_id": user.UserID,
		"name":    user.Name,
		"email":   user.Email,
		"role":    user.Role,
	}
	if user.SupplierID != nil {
		resp["supplier_id"] = *user.SupplierID
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteSuppAdmin(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can delete users"})
		c.Abort()
		return
	}

	u_id := c.Param("id")
	uid, err := strconv.Atoi(u_id)
	if err != nil || uid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	result, err := config.DB.ExecContext(ctx, "delete from users where user_id = ?", uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		c.Abort()
		return
	}
	r, _ := result.RowsAffected()
	if r == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func UpdateSuppAdmin(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can delete users"})
		c.Abort()
		return
	}

	u_id := c.Param("id")
	uid, err := strconv.Atoi(u_id)
	if err != nil || uid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	var body struct {
		Name       string `json:"name,omitempty"`
		Email      string `json:"email,omitempty"`
		Password   string `json:"password,omitempty"`
		SupplierID int    `json:"supplier_id,omitempty"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}

	if body.Name == "" && body.Email == "" && body.Password == "" && body.SupplierID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var existing models.User
	var supp sql.NullInt64
	err = config.DB.QueryRowContext(ctx, "select user_id, name, email, password_hash, role, supplier_id from users where user_id=?", uid).Scan(&existing.UserID, &existing.Name, &existing.Email, &existing.PasswordHash, &existing.Role, &supp)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching user"})
		c.Abort()
		return
	}

	if supp.Valid {
		tmp := int(supp.Int64)
		existing.SupplierID = &tmp
	} else {
		existing.SupplierID = nil
	}

	newName := existing.Name
	newEmail := existing.Email
	newPwdHash := existing.PasswordHash


	if strings.TrimSpace(body.Name) != "" {
		newName = strings.TrimSpace(body.Name)
	}

	if strings.TrimSpace(body.Email) != "" {

		email := strings.TrimSpace(body.Email)

		if err := utils.IsValidEmail(email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var exists bool
		err := config.DB.QueryRowContext(ctx, "select exists(select 1 from users where lower(email)=lower(?) and user_id<>?)", email, uid).Scan(&exists)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while checking email"})
			c.Abort()
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			c.Abort()
			return
		}

		newEmail = email
	}

	if strings.TrimSpace(body.Password) != "" {

		if err := utils.IsValidPassword(body.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		hash, err := utils.HashPwd(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			c.Abort()
			return
		}

		newPwdHash = hash
	}

	_, err = config.DB.ExecContext(ctx, "update users set name=?, email=?,password_hash=?, WHERE user_id=?", newName, newEmail, newPwdHash, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "user updated successfully",
		"user_id":     uid,
		"name":        newName,
		"email":       newEmail,
		"supplier_id": existing.SupplierID,
		"role":        existing.Role,
	})
}
