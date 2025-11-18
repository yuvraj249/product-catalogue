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
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
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

	supplierID, ok := claims["supplier_id"]
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
		c.Abort()
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "user_id missing in token"})
		c.Abort()
		return
	}
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
	suppID := int(supplierID.(float64))
	if suppID != prodSuppID {
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
	claims, ok := GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}
	role := c.GetString("role")
	prodFilter := 0
	if p := c.Query("product_id"); p != "" {
		pid, err := strconv.Atoi(p)
		if err != nil || pid <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
			c.Abort()
			return
		}
		prodFilter = pid
	}
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var args []interface{}
	var query string
	var rows *sql.Rows
	var err error
	switch role {
	case "system_admin":
		if prodFilter != 0 {
			query = "select stock_id, product_id,quantity,  movement_type, reason, performed_by from stock_movements where product_id = ? order by stock_id desc "
			args = []interface{}{prodFilter}
		} else {
			query = "select stock_id, product_id,quantity,  movement_type, reason, performed_by from stock_movements order by stock_id desc "
			args = nil
		}
		rows, err = config.DB.QueryContext(ctx, query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching stock movements"})
			c.Abort()
			return
		}
	case "supplier_admin":
		supplierID, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
			c.Abort()
			return
		}
		suppID := int(supplierID.(float64))
		company, err := SupplierCompany(ctx, suppID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusForbidden, gin.H{"error": "supplier not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch supplier company"})
			}
			c.Abort()
			return
		}
		if prodFilter != 0 {
			query = "select sm.stock_id, sm.product_id,sm.quantity,  sm.movement_type, sm.reason, sm.performed_by from stock_movements sm join products p on sm.product_id = p.product_id join suppliers s on p.product_supplier_id = s.supplier_id where s.company = ? and sm.product_id = ? order by sm.stock_id desc "
			args = []interface{}{company, prodFilter}
		} else {
			query = "select sm.stock_id, sm.product_id,sm.quantity,  sm.movement_type, sm.reason, sm.performed_by from stock_movements sm join products p on sm.product_id = p.product_id join suppliers s on p.product_supplier_id = s.supplier_id where s.company = ? order by sm.stock_id desc "
			args = []interface{}{company}

		}
		rows, err = config.DB.QueryContext(ctx, query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error while fetching stock movements"})
			c.Abort()
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
		return
	}

	defer rows.Close()

	movements := []models.Stock_Movement{}

	for rows.Next() {
		var sm models.Stock_Movement
		var reason sql.NullString
		var performed_by sql.NullInt64
		if err := rows.Scan(&sm.StockID, &sm.ProductID, &sm.Quantity, &sm.MovementType, &reason, &performed_by); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning stock movements"})
			c.Abort()
			return
		}
		if reason.Valid {
			sm.Reason = reason.String
		} else {
			sm.Reason = ""
		}
		if performed_by.Valid {
			sm.PerformedBy = int(performed_by.Int64)
		} else {
			sm.PerformedBy = 0
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
