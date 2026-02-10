package functions

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"product-catalogue/config"
	"product-catalogue/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var lowStockAlert int

func init() {
	lowStockAlert = 10
	if s := os.Getenv("LOW_STOCK_ALERT"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v >= 0 {
			lowStockAlert = v
		}
	}
}

func SupplierCompany(ctx context.Context, supplierID int) (string, error) {
	var company string
	err := config.DB.QueryRowContext(ctx, "select company from suppliers where supplier_id = ? ", supplierID).Scan(&company)
	if err != nil {
		return "", err

	}
	return company, nil
}

func SupplierIDCompany(ctx context.Context, company string) ([]int, error) {
	rows, err := config.DB.QueryContext(ctx, "select supplier_id from suppliers where company = ?")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	idStore := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		idStore = append(idStore, id)

	}
	return idStore, nil

}

func CountStock(ctx context.Context, productID int) (int, error) {
	var total sql.NullInt64
	err := config.DB.QueryRowContext(ctx, "select ifnull(sum(case when movement_type = 'IN' then quantity when movement_type = 'OUT' then -quantity end), 0) from stock_movements where product_id = ?", productID).Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return int(total.Int64), nil
}

func CreateStockMovement(c *gin.Context) {

	role := c.GetString("role")

	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier_admin can create stock movements"})
		c.Abort()
		return
	}
	var reqSt struct {
		ProductID    int    `json:"product_id"`
		Quantity     int    `json:"quantity"`
		MovementType string `json:"movement_type"`
		Reason       string `json:"reason,omitempty"`
	}
	if err := c.BindJSON(&reqSt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}

	reqSt.MovementType = strings.ToUpper(strings.TrimSpace(reqSt.MovementType))
	if reqSt.MovementType != "IN" && reqSt.MovementType != "OUT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stock movement type must be IN or OUT"})
		c.Abort()
		return

	}

	if reqSt.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stock qunatity must be greater than 0"})
		c.Abort()
		return
	}

	supplierID := c.GetInt("supplier_id")
	userID := c.GetInt("user_id")
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var prodSuppID int
	err := config.DB.QueryRowContext(ctx, "select product_supplier_id from products where product_id=?", reqSt.ProductID).Scan(&prodSuppID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching product"})
		c.Abort()
		return
	}

	if supplierID != prodSuppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to change stock for this product"})
		c.Abort()
		return
	}

	if reqSt.MovementType == "OUT" {
		current, err := CountStock(ctx, reqSt.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count current stock"})
			c.Abort()
			return
		}
		if current-reqSt.Quantity < lowStockAlert {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "warning.stock going below threshhold value",
				"current_stock": current,
				"min_quantity":  lowStockAlert,
			})
			c.Abort()
			return
		}
	}

	var reason sql.NullString
	if reqSt.Reason != "" && strings.TrimSpace(reqSt.Reason) != "" {
		reason = sql.NullString{String: strings.TrimSpace(reqSt.Reason), Valid: true}
	}

	query := "insert into stock_movements (product_id, quantity, movement_type, reason, performed_by) values(?, ?, ?, ?, ?)"
	result, err := config.DB.ExecContext(ctx, query, reqSt.ProductID, reqSt.Quantity, reqSt.MovementType, reason, int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stock movement"})
		c.Abort()
		return
	}

	id, _ := result.LastInsertId()
	newStock, _ := CountStock(ctx, reqSt.ProductID)
	c.JSON(http.StatusCreated, gin.H{
		"message":       "stock movement created",
		"stock_id":      id,
		"product_id":    reqSt.ProductID,
		"current_stock": newStock,
		"changed_by":    int(userID),
	})

}

func GetStockMovements(c *gin.Context) {
	role := c.GetString("role")
	search := strings.ToLower(strings.TrimSpace(c.Query("q")))
	stockID := "%" + search + "%"
	productName := "%" + search + "%"
	quantity := "%" + search + "%"
	movementType := "%" + search + "%"
	reason := "%" + search + "%"
	username := "%" + search + "%"
	ctx, cancel := CtxTimeout(c)
	defer cancel()
	var rows *sql.Rows
	var err error
	switch role {
	case "system_admin":
		rows, err = config.DB.QueryContext(ctx, "select sm.stock_id, sm.product_id, p.product_name, sm.quantity, sm.movement_type, sm.reason, sm.performed_by, coalesce(u.name,'') as username from stock_movements sm join products p on sm.product_id = p.product_id left join users u on sm.performed_by = u.user_id where ?='' or cast(sm.stock_id as char) like ? or lower(p.product_name) like ? or cast(sm.quantity as char) like ? or lower(sm.movement_type) like ? or lower(coalesce(sm.reason,'')) like ? or lower(coalesce(u.name,'')) like ? order by sm.stock_id desc", search, stockID, productName, quantity, movementType, reason, username)

	case "supplier_admin":
		supplierID := c.GetInt("supplier_id")
		rows, err = config.DB.QueryContext(
			ctx,
			"select sm.stock_id, sm.product_id, p.product_name, sm.quantity, sm.movement_type, sm.reason, sm.performed_by, '' as username from stock_movements sm join products p on sm.product_id = p.product_id where p.product_supplier_id=? and (?='' or cast(sm.stock_id as char) like ? or lower(p.product_name) like ? or cast(sm.quantity as char) like ? or lower(sm.movement_type) like ? or lower(coalesce(sm.reason,'')) like ?) order by sm.stock_id desc",
			supplierID, search, stockID, productName, quantity, movementType, reason, username,
		)
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "db query failed",
			"details": err.Error(),
		})
		c.Abort()
		return
	}

	defer rows.Close()

	movements := []models.Stock_Movement{}

	for rows.Next() {
		var sm models.Stock_Movement
		var reason sql.NullString
		var username sql.NullString
		if err := rows.Scan(&sm.StockID, &sm.ProductID, &sm.ProductName, &sm.Quantity, &sm.MovementType, &reason, &sm.PerformedBy, &username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning stock movements"})
			c.Abort()
			return
		}
		if reason.Valid {
			sm.Reason = reason.String
		} else {
			sm.Reason = ""
		}
		if username.Valid {
			sm.Username = username.String
		} else {
			sm.Username = ""
		}

		movements = append(movements, sm)

	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading stock movements"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"movements": movements})

}

func UpdateStockMovement(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		c.Abort()
		return

	}
	role := c.GetString("role")
	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier_admin can update stock movements"})
		c.Abort()
		return
	}

	var req struct {
		Quantity     int    `json:"quantity"`
		MovementType string `json:"movement_type"`
		Reason       string `json:"reason"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
	}

	req.MovementType = strings.ToUpper(strings.TrimSpace(req.MovementType))
	if req.MovementType != "IN" && req.MovementType != "OUT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movement_type must be IN or OUT"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be > 0"})
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var productID int
	var oldQty int
	var oldType string

	err = config.DB.QueryRowContext(ctx, "select product_id, quantity, movement_type from stock_movements where stock_id = ?", id).Scan(&productID, &oldQty, &oldType)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock movement not found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error fetching movement"})
		c.Abort()
		return
	}

	supplierID := c.GetInt("supplier_id")

	var prodSuppID int
	err = config.DB.QueryRowContext(ctx, "select product_supplier_id from products where product_id = ?", productID).Scan(&prodSuppID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error fetching product"})
		return
	}

	if supplierID != prodSuppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this product"})
		return
	}

	currentStock, err := CountStock(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count stock"})
		return
	}

	if oldType == "IN" {
		currentStock -= oldQty
	} else {
		currentStock += oldQty
	}

	if req.MovementType == "IN" {
		currentStock += req.Quantity
	} else {
		currentStock -= req.Quantity
	}

	if currentStock < lowStockAlert {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "stock would go below threshold",
			"current_stock": currentStock,
			"min_quantity":  lowStockAlert,
		})
		return
	}

	var reason sql.NullString
	if strings.TrimSpace(req.Reason) != "" {
		reason = sql.NullString{String: strings.TrimSpace(req.Reason), Valid: true}
	}

	_, err = config.DB.ExecContext(ctx, "update stock_movements set quantity = ? , movement_type = ? , reason = ? where stock_id = ?", req.Quantity, req.MovementType, reason, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stock movement"})
		return
	}

	newStock, err := CountStock(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count stock movement"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "stock movement updated",
		"stock_id":      id,
		"product_id":    productID,
		"current_stock": newStock,
	})
}

func DeleteStockMovement(c *gin.Context) {
	p_id := c.Param("id")
	id, err := strconv.Atoi(p_id)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		c.Abort()
		return

	}
	role := c.GetString("role")
	if role != "supplier_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only supplier_admin can delete stock movements"})
		c.Abort()
		return
	}

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var productID, qty int
	var mtype string

	err = config.DB.QueryRowContext(ctx, "select product_id, quantity, movement_type from stock_movements where stock_id = ?", id).Scan(&productID, &qty, &mtype)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock movement not found"})
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error fetching movement"})
		c.Abort()
		return
	}

	supplierID := c.GetInt("supplier_id")
	var prodSuppID int
	err = config.DB.QueryRowContext(ctx, "select product_supplier_id from products where product_id = ?", productID).Scan(&prodSuppID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error fetching product"})
		c.Abort()
		return
	}

	if supplierID != prodSuppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this product"})
		c.Abort()
		return
	}

	currentStock, err := CountStock(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count stock movement"})
		c.Abort()
		return
	}

	if mtype == "IN" {
		currentStock -= qty
	} else {
		currentStock += qty
	}

	_, err = config.DB.ExecContext(ctx, "delete from stock_movements where stock_id = ?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete stock movement"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock movement deleted"})

}
