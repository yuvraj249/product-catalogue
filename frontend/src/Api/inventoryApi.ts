import { apiClient } from './client';
import { InventoryItem } from '../types';

export interface ListItemsResponse {
  items: InventoryItem[];
  total: number;
}

export const inventoryApi = {
  listItems: async (query = '', limit = 25, offset = 0): Promise<ListItemsResponse> => {
    const res = await apiClient.get<ListItemsResponse>('/api/v1/inventory/items', {
      params: { q: query, limit, offset },
    });
    return res.data;
  },

  getItemByBarcode: async (barcode: string): Promise<InventoryItem> => {
    const res = await apiClient.get<InventoryItem>('/api/v1/inventory/barcode', {
      params: { barcode },
    });
    return res.data;
  },

  createItem: async (item: Partial<InventoryItem>): Promise<InventoryItem> => {
    const res = await apiClient.post<InventoryItem>('/api/v1/inventory/items', item);
    return res.data;
  },

  transferStock: async (data: {
    item_id: string;
    source_location_id: string;
    destination_location_id: string;
    quantity: number;
    reference_document: string;
  }) => {
    const res = await apiClient.post('/api/v1/inventory/transfer', data);
    return res.data;
  },

  bulkScrapStock: async (items: { item_id: string; source_location_id?: string; quantity: number }[]) => {
    const promises = items.map((it) =>
      apiClient.post('/api/v1/inventory/scrap', {
        item_id: it.item_id,
        source_location_id: it.source_location_id || '00000000-0000-0000-0000-000000000000',
        quantity: it.quantity,
        reference_document: 'BULK-SCRAP-AUDIT',
      })
    );
    return Promise.all(promises);
  },
};
