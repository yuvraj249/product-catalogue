package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type InventoryRepository struct {
	DB *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) TransferStock(ctx context.Context, movement entity.StockMovement) error {
	if movement.SourceLocationID == movement.DestinationLocationID {
		return entity.ErrIdenticalLocations
	}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Identify Source Location Type
	var srcType string
	err = tx.QueryRowContext(ctx,
		`SELECT type FROM locations WHERE id = $1 AND tenant_id = $2`,
		movement.SourceLocationID, movement.TenantID,
	).Scan(&srcType)
	if err != nil {
		return fmt.Errorf("failed to verify source location: %w", err)
	}

	// 2. Lock and Validate internal locations
	if srcType == "Internal" {
		var currentStock int
		queryLock := `
			SELECT quantity 
			FROM stock_quantities 
			WHERE tenant_id = $1 AND location_id = $2 AND item_id = $3 
			FOR UPDATE`

		err = tx.QueryRowContext(ctx, queryLock, movement.TenantID, movement.SourceLocationID, movement.ItemID).Scan(&currentStock)
		if err == sql.ErrNoRows || currentStock < movement.Quantity {
			return entity.ErrInsufficientStock
		} else if err != nil {
			return fmt.Errorf("error acquiring pessimistic stock lock: %w", err)
		}

		// Decrement
		_, err = tx.ExecContext(ctx, `
			UPDATE stock_quantities 
			SET quantity = quantity - $1, updated_at = NOW() 
			WHERE tenant_id = $2 AND location_id = $3 AND item_id = $4`,
			movement.Quantity, movement.TenantID, movement.SourceLocationID, movement.ItemID)
		if err != nil {
			return fmt.Errorf("failed to decrement source stock: %w", err)
		}
	}

	// 3. Increment Destination Location
	upsertDest := `
		INSERT INTO stock_quantities (tenant_id, location_id, item_id, quantity, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (tenant_id, location_id, item_id) 
		DO UPDATE SET quantity = stock_quantities.quantity + EXCLUDED.quantity, updated_at = NOW()`

	_, err = tx.ExecContext(ctx, upsertDest, movement.TenantID, movement.DestinationLocationID, movement.ItemID, movement.Quantity)
	if err != nil {
		// Fallback for SQLite upsert if ON CONFLICT syntax differs
		upsertFallback := `
			INSERT INTO stock_quantities (tenant_id, location_id, item_id, quantity, updated_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
			ON CONFLICT(tenant_id, location_id, item_id) 
			DO UPDATE SET quantity = stock_quantities.quantity + $4, updated_at = CURRENT_TIMESTAMP`
		_, err = tx.ExecContext(ctx, upsertFallback, movement.TenantID, movement.DestinationLocationID, movement.ItemID, movement.Quantity)
		if err != nil {
			return fmt.Errorf("failed to increment destination stock: %w", err)
		}
	}

	// 4. Record Immutable Movement
	movementQuery := `
		INSERT INTO stock_movements (id, tenant_id, item_id, source_location_id, destination_location_id, quantity, reference_document, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	if movement.ID == uuid.Nil {
		movement.ID = uuid.New()
	}

	_, err = tx.ExecContext(ctx, movementQuery,
		movement.ID, movement.TenantID, movement.ItemID,
		movement.SourceLocationID, movement.DestinationLocationID,
		movement.Quantity, movement.ReferenceDocument, movement.CreatedBy)
	if err != nil {
		movementQueryFallback := `
			INSERT INTO stock_movements (id, tenant_id, item_id, source_location_id, destination_location_id, quantity, reference_document, created_by, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)`
		_, err = tx.ExecContext(ctx, movementQueryFallback,
			movement.ID, movement.TenantID, movement.ItemID,
			movement.SourceLocationID, movement.DestinationLocationID,
			movement.Quantity, movement.ReferenceDocument, movement.CreatedBy)
		if err != nil {
			return fmt.Errorf("failed to append double-entry stock movement: %w", err)
		}
	}

	return tx.Commit()
}

func (r *InventoryRepository) CreateItem(ctx context.Context, item *entity.InventoryItem) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	query := `
		INSERT INTO inventory_items (id, tenant_id, sku, name, barcode, cost_price, sale_price, min_stock_threshold, max_stock_threshold, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.DB.ExecContext(ctx, query,
		item.ID, item.TenantID, item.SKU, item.Name, item.Barcode,
		item.CostPrice, item.SalePrice, item.MinStockThreshold, item.MaxStockThreshold,
		item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory item: %w", err)
	}
	return nil
}

func (r *InventoryRepository) GetItemByID(ctx context.Context, tenantID, itemID uuid.UUID) (*entity.InventoryItem, error) {
	query := `
		SELECT i.id, i.tenant_id, i.sku, i.name, COALESCE(i.barcode, ''), i.cost_price, i.sale_price, i.min_stock_threshold, i.max_stock_threshold, COALESCE(SUM(sq.quantity), 0) as current_stock, i.created_at, i.updated_at
		FROM inventory_items i
		LEFT JOIN stock_quantities sq ON sq.tenant_id = i.tenant_id AND sq.item_id = i.id
		WHERE i.tenant_id = $1 AND i.id = $2
		GROUP BY i.id`

	var item entity.InventoryItem
	err := r.DB.QueryRowContext(ctx, query, tenantID, itemID).Scan(
		&item.ID, &item.TenantID, &item.SKU, &item.Name, &item.Barcode,
		&item.CostPrice, &item.SalePrice, &item.MinStockThreshold, &item.MaxStockThreshold,
		&item.CurrentStock, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch item by ID: %w", err)
	}
	return &item, nil
}

func (r *InventoryRepository) GetItemByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*entity.InventoryItem, error) {
	query := `
		SELECT i.id, i.tenant_id, i.sku, i.name, COALESCE(i.barcode, ''), i.cost_price, i.sale_price, i.min_stock_threshold, i.max_stock_threshold, COALESCE(SUM(sq.quantity), 0) as current_stock, i.created_at, i.updated_at
		FROM inventory_items i
		LEFT JOIN stock_quantities sq ON sq.tenant_id = i.tenant_id AND sq.item_id = i.id
		WHERE i.tenant_id = $1 AND i.barcode = $2
		GROUP BY i.id`

	var item entity.InventoryItem
	err := r.DB.QueryRowContext(ctx, query, tenantID, barcode).Scan(
		&item.ID, &item.TenantID, &item.SKU, &item.Name, &item.Barcode,
		&item.CostPrice, &item.SalePrice, &item.MinStockThreshold, &item.MaxStockThreshold,
		&item.CurrentStock, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch item by barcode: %w", err)
	}
	return &item, nil
}

func (r *InventoryRepository) ListItems(ctx context.Context, tenantID uuid.UUID, search string, limit, offset int) ([]entity.InventoryItem, int, error) {
	searchPattern := "%" + search + "%"
	countQuery := `SELECT COUNT(*) FROM inventory_items WHERE tenant_id = $1 AND ($2 = '' OR sku LIKE $3 OR name LIKE $3 OR barcode LIKE $3)`

	var total int
	err := r.DB.QueryRowContext(ctx, countQuery, tenantID, search, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory items: %w", err)
	}

	query := `
		SELECT i.id, i.tenant_id, i.sku, i.name, COALESCE(i.barcode, ''), i.cost_price, i.sale_price, i.min_stock_threshold, i.max_stock_threshold, COALESCE(SUM(sq.quantity), 0) as current_stock, i.created_at, i.updated_at
		FROM inventory_items i
		LEFT JOIN stock_quantities sq ON sq.tenant_id = i.tenant_id AND sq.item_id = i.id
		WHERE i.tenant_id = $1 AND ($2 = '' OR i.sku LIKE $3 OR i.name LIKE $3 OR i.barcode LIKE $3)
		GROUP BY i.id
		ORDER BY i.created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.DB.QueryContext(ctx, query, tenantID, search, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query inventory items: %w", err)
	}
	defer rows.Close()

	items := []entity.InventoryItem{}
	for rows.Next() {
		var item entity.InventoryItem
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.SKU, &item.Name, &item.Barcode,
			&item.CostPrice, &item.SalePrice, &item.MinStockThreshold, &item.MaxStockThreshold,
			&item.CurrentStock, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *InventoryRepository) GetLocationByID(ctx context.Context, tenantID, locationID uuid.UUID) (*entity.Location, error) {
	query := `SELECT id, tenant_id, name, type, is_scrap, created_at FROM locations WHERE tenant_id = $1 AND id = $2`
	var loc entity.Location
	err := r.DB.QueryRowContext(ctx, query, tenantID, locationID).Scan(&loc.ID, &loc.TenantID, &loc.Name, &loc.Type, &loc.IsScrap, &loc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location: %w", err)
	}
	return &loc, nil
}

func (r *InventoryRepository) GetOrCreateScrapLocation(ctx context.Context, tenantID uuid.UUID) (*entity.Location, error) {
	query := `SELECT id, tenant_id, name, type, is_scrap, created_at FROM locations WHERE tenant_id = $1 AND is_scrap = TRUE LIMIT 1`
	var loc entity.Location
	err := r.DB.QueryRowContext(ctx, query, tenantID).Scan(&loc.ID, &loc.TenantID, &loc.Name, &loc.Type, &loc.IsScrap, &loc.CreatedAt)
	if err == nil {
		return &loc, nil
	}

	loc = entity.Location{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      "Scrap Location",
		Type:      entity.LocationInventoryLoss,
		IsScrap:   true,
		CreatedAt: time.Now(),
	}

	insert := `INSERT INTO locations (id, tenant_id, name, type, is_scrap, created_at) VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err = r.DB.ExecContext(ctx, insert, loc.ID, loc.TenantID, loc.Name, loc.Type, loc.IsScrap)
	if err != nil {
		insertFallback := `INSERT INTO locations (id, tenant_id, name, type, is_scrap, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)`
		_, err = r.DB.ExecContext(ctx, insertFallback, loc.ID, loc.TenantID, loc.Name, loc.Type, loc.IsScrap)
		if err != nil {
			return nil, fmt.Errorf("failed to create scrap location: %w", err)
		}
	}

	return &loc, nil
}

func (r *InventoryRepository) GetOrCreateDefaultLocations(ctx context.Context, tenantID uuid.UUID) (mainWhID, customerLocID uuid.UUID, err error) {
	// Main Warehouse
	var mainLoc entity.Location
	err = r.DB.QueryRowContext(ctx, `SELECT id FROM locations WHERE tenant_id = $1 AND type = 'Internal' LIMIT 1`, tenantID).Scan(&mainLoc.ID)
	if err != nil {
		mainLoc = entity.Location{
			ID:        uuid.New(),
			TenantID:  tenantID,
			Name:      "Main Warehouse",
			Type:      entity.LocationInternal,
			IsScrap:   false,
			CreatedAt: time.Now(),
		}
		_, _ = r.DB.ExecContext(ctx, `INSERT INTO locations (id, tenant_id, name, type, is_scrap, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)`,
			mainLoc.ID, mainLoc.TenantID, mainLoc.Name, mainLoc.Type, mainLoc.IsScrap)
	}

	// Virtual Customer Delivery Location
	var custLoc entity.Location
	err = r.DB.QueryRowContext(ctx, `SELECT id FROM locations WHERE tenant_id = $1 AND type = 'Customer' LIMIT 1`, tenantID).Scan(&custLoc.ID)
	if err != nil {
		custLoc = entity.Location{
			ID:        uuid.New(),
			TenantID:  tenantID,
			Name:      "Customer Delivery Location",
			Type:      entity.LocationCustomer,
			IsScrap:   false,
			CreatedAt: time.Now(),
		}
		_, _ = r.DB.ExecContext(ctx, `INSERT INTO locations (id, tenant_id, name, type, is_scrap, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)`,
			custLoc.ID, custLoc.TenantID, custLoc.Name, custLoc.Type, custLoc.IsScrap)
	}

	return mainLoc.ID, custLoc.ID, nil
}
