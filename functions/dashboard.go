package functions

import (
	"database/sql"
	"net/http"
	"os"
	"product-catalogue/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetClaims6(c *gin.Context) (jwt.MapClaims, bool) {
	claimsExist, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claims, ok := claimsExist.(jwt.MapClaims)
	return claims, ok

}

func GetDashboard(c *gin.Context) {
	claims, ok := GetClaims6(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return
	}

	role, _ := claims["role"].(string)
	lowStock := 10
	if st := os.Getenv("LOW_STOCK_ALERT"); st != "" {
		if vl, err := strconv.Atoi(st); err == nil && vl >= 0 {
			lowStock = vl
		}
	}
	ctx, cancel := CtxTimeout5(c)
	defer cancel()

	type lowStockProd struct {
		ProductID int    `json:"product_id"`
		Name      string `json:"product_name"`
		Stock     int    `json:"current_stock"`
		Supplier  *int   `json:"product_supplier_id,omitempty"`
	}

	resp := gin.H{}

	if role == "system_admin" {
		var totalProd, totalCatg, totalSupp, totalSuppAdmin int
		err := config.DB.QueryRowContext(ctx, "select count(*) from products").Scan(&totalProd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count products"})
			c.Abort()
			return
		}
		err = config.DB.QueryRowContext(ctx, "select count(*) from categories").Scan(&totalCatg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count categories"})
			c.Abort()
			return
		}
		err = config.DB.QueryRowContext(ctx, "select count(*) from suppliers").Scan(&totalSupp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count suppliers"})
			c.Abort()
			return
		}
		err = config.DB.QueryRowContext(ctx, "select count(*) from users where role='supplier_admin'").Scan(&totalSuppAdmin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count supplier admins"})
			c.Abort()
			return
		}

		resp["total_products"] = totalProd
		resp["total_categories"] = totalCatg
		resp["total_suppliers"] = totalSupp
		resp["total_supplier_admins"] = totalSuppAdmin

		rows, err := config.DB.QueryContext(ctx, "select p.product_id , p.product_name, ifnull(sum(case when sm.movement_type = 'IN' then sm.quantity when sm.movement_type = 'OUT' then -sm.quantity else 0 end ), 0) as stock from products p left join stock_movements sm on p.product_id = sm.product_id group by p.product_id, p.product_name having stock < ? order by stock asc", lowStock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch low stock products"})
			c.Abort()
			return
		}
		defer rows.Close()

		type lowStock1 struct {
			ProductID int    `json:"product_id"`
			Name      string `json:"product_name"`
			Stock     int    `json:"current_stock"`
		}

		lowStockList := []lowStock1{}
		for rows.Next() {
			var ls lowStock1
			err := rows.Scan(&ls.ProductID, &ls.Name, &ls.Stock)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning low stock products"})
				c.Abort()
				return
			}
			lowStockList = append(lowStockList, ls)
		}
		resp["low_stock_threshold"] = lowStock
		resp["low_stock_products"] = lowStockList

		c.JSON(http.StatusOK, resp)
		return
	}

	if role == "supplier_admin" {
		supplierVal, ok := claims["supplier_id"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "supplier_id missing in token"})
			c.Abort()
			return
		}
		suppID := int(supplierVal.(float64))
		var company string
		err := config.DB.QueryRowContext(ctx, "select company from suppliers where supplier_id=?", suppID).Scan(&company)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch company"})
			c.Abort()
			return
		}
		resp["company"] = company
		var totalProd, totalCmpSupp int
		err = config.DB.QueryRowContext(ctx, "select count(*) from products p join suppliers s on p.product_supplier_id = s.supplier_id where s.company= ?", company).Scan(&totalProd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count company products"})
			c.Abort()
			return
		}

		err = config.DB.QueryRowContext(ctx, "select count(*) from suppliers where company=?", company).Scan(&totalCmpSupp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count company suppliers"})
			c.Abort()
			return
		}
		resp["total_products"] = totalProd
		resp["company_suppliers"] = totalCmpSupp

		rows, err := config.DB.QueryContext(ctx, "select p.product_id , p.product_name, ifnull(sum(case when sm.movement_type = 'IN' then sm.quantity when sm.movement_type = 'OUT' then -sm.quantity else 0 end ), 0) as stock from products p join suppliers s on p.product_supplier_id = s.supplier_id left join stock_movements sm on p.product_id = sm.product_id where s.company=? group by p.product_id, p.product_name having stock < ? order by stock asc", company, lowStock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch low stock products for company"})
			c.Abort()
			return
		}
		defer rows.Close()

		type LowStock2 struct {
			ProductID int    `json:"product_id"`
			Name      string `json:"product_name"`
			Stock     int    `json:"current_stock"`
		}
		lowStockList := []LowStock2{}

		for rows.Next() {
			var ls LowStock2
			err := rows.Scan(&ls.ProductID, &ls.Name, &ls.Stock)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning low stock products"})
				return
			}
			lowStockList = append(lowStockList, ls)
		}

		resp["low_stock_threshold"] = lowStock
		resp["low_stock_products"] = lowStockList

		mvntRows, err := config.DB.QueryContext(ctx, "select sm.stock_id, sm.product_id, p.product_name,sm.quantity, sm.movement_type, sm.reason from stock_movements sm join products p on sm.product_id = p.product_id join suppliers s on p.product_supplier_id = s.supplier_id where s.company = ? order by sm.stock_id desc", company)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recent stock movements"})
			c.Abort()
			return
		}
		defer mvntRows.Close()

		type stMove struct {
			StockID      int    `json:"stock_id"`
			ProductID    int    `json:"product_id"`
			ProductName  string `json:"product_name"`
			MovementType string `json:"movement_type"`
			Quantity     int    `json:"quantity"`
			Reason       string `json:"reason"`
		}

		recent := []stMove{}

		for mvntRows.Next() {
			var m stMove
			var reason sql.NullString
			err := mvntRows.Scan(&m.StockID, &m.ProductID, &m.ProductName, &m.MovementType, &m.Quantity, &reason)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning movements"})
				c.Abort()
				return
			}
			if reason.Valid {
				m.Reason = reason.String
			} else {
				m.Reason = ""
			}
			recent = append(recent, m)
		}
		resp["recent_stock_movements"] = recent

		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized role"})

}
