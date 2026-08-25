export type UserRole = 'SuperAdmin' | 'TenantAdmin' | 'WarehouseManager' | 'Cashier' | 'Auditor';
export type LocationType = 'Internal' | 'Customer' | 'Vendor' | 'InventoryLoss' | 'Production';
export type InvoiceStatus = 'Draft' | 'Unpaid' | 'Paid' | 'Cancelled';

export interface Tenant {
  id: string;
  name: string;
  subdomain: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  tenant_id: string;
  email: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}

export interface Location {
  id: string;
  tenant_id: string;
  name: string;
  type: LocationType;
  is_scrap: boolean;
  created_at: string;
}

export interface InventoryItem {
  id: string;
  tenant_id: string;
  sku: string;
  name: string;
  barcode: string;
  cost_price: number;
  sale_price: number;
  min_stock_threshold: number;
  max_stock_threshold: number;
  current_stock: number;
  created_at: string;
  updated_at: string;
}

export interface StockMovement {
  id: string;
  tenant_id: string;
  item_id: string;
  source_location_id: string;
  destination_location_id: string;
  quantity: number;
  reference_document: string;
  created_by: string;
  created_at: string;
}

export interface InvoiceItem {
  id: string;
  tenant_id: string;
  invoice_id: string;
  item_id: string;
  quantity: number;
  unit_price: number;
}

export interface Invoice {
  id: string;
  tenant_id: string;
  customer_name: string;
  total_amount: number;
  status: InvoiceStatus;
  stripe_payment_intent_id?: string;
  idempotency_key?: string;
  items?: InvoiceItem[];
  created_at: string;
  updated_at: string;
}
