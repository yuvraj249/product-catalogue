CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
CREATE TYPE user_role AS ENUM ('SuperAdmin', 'TenantAdmin', 'WarehouseManager', 'Cashier', 'Auditor');
CREATE TYPE location_type AS ENUM ('Internal', 'Customer', 'Vendor', 'InventoryLoss', 'Production');
CREATE TYPE invoice_status AS ENUM ('Draft', 'Unpaid', 'Paid', 'Cancelled');

-- Tenants
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'Cashier',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_tenant_user_email UNIQUE (tenant_id, email)
);

-- Locations (Physical warehouses, Virtual partner drops, Scrap/Loss accounts)
CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    type location_type NOT NULL,
    is_scrap BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_tenant_location_name UNIQUE (tenant_id, name)
);

-- Inventory Items (Master SKU Record)
CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    barcode VARCHAR(100),
    cost_price NUMERIC(12, 4) NOT NULL CHECK (cost_price >= 0),
    sale_price NUMERIC(12, 4) NOT NULL CHECK (sale_price >= 0),
    min_stock_threshold INT NOT NULL DEFAULT 0,
    max_stock_threshold INT NOT NULL DEFAULT 1000,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_tenant_sku UNIQUE (tenant_id, sku),
    CONSTRAINT uq_tenant_barcode UNIQUE (tenant_id, barcode)
);

-- Point-In-Time Stock Balances
CREATE TABLE IF NOT EXISTS stock_quantities (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, location_id, item_id)
);

-- Immutable Double-Entry Ledger for Stock Movements
CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    source_location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    destination_location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    reference_document VARCHAR(100) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Invoices
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    customer_name VARCHAR(255) NOT NULL,
    total_amount NUMERIC(12, 4) NOT NULL CHECK (total_amount >= 0),
    status invoice_status NOT NULL DEFAULT 'Draft',
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    idempotency_key UUID UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12, 4) NOT NULL CHECK (unit_price >= 0)
);

-- Performance Indexing
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_quantities_lookup ON stock_quantities(tenant_id, location_id, item_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_audit ON stock_movements(tenant_id, item_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_inventory_items_barcode ON inventory_items(tenant_id, barcode);
CREATE INDEX IF NOT EXISTS idx_invoices_stripe ON invoices(stripe_payment_intent_id);
