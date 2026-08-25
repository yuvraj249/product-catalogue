import React from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { InventoryItem } from '../../types';
import { Badge } from '../../Components/ui/badge';
import { Checkbox } from '../../Components/ui/checkbox';

export const inventoryColumns: ColumnDef<InventoryItem>[] = [
  {
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={table.getIsAllPageRowsSelected()}
        indeterminate={table.getIsSomePageRowsSelected()}
        onChange={table.getToggleAllPageRowsSelectedHandler()}
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        disabled={!row.getCanSelect()}
        onChange={row.getToggleSelectedHandler()}
      />
    ),
  },
  {
    accessorKey: 'sku',
    header: 'SKU',
    cell: ({ row }) => (
      <span className="font-mono font-medium text-indigo-700 bg-indigo-50 px-2 py-1 rounded">
        {row.original.sku}
      </span>
    ),
  },
  {
    accessorKey: 'name',
    header: 'Item Name',
    cell: ({ row }) => <span className="font-semibold text-gray-900">{row.original.name}</span>,
  },
  {
    accessorKey: 'barcode',
    header: 'Barcode',
    cell: ({ row }) => (
      <span className="font-mono text-gray-600 text-xs">
        {row.original.barcode || '—'}
      </span>
    ),
  },
  {
    accessorKey: 'sale_price',
    header: 'Price',
    cell: ({ row }) => (
      <span className="font-medium text-gray-900">
        ${row.original.sale_price.toFixed(2)}
      </span>
    ),
  },
  {
    accessorKey: 'current_stock',
    header: 'Stock Level',
    cell: ({ row }) => {
      const stock = row.original.current_stock || 0;
      const min = row.original.min_stock_threshold || 0;

      let variant: 'destructive' | 'warning' | 'success' = 'success';
      let label = 'Optimal';

      if (stock === 0) {
        variant = 'destructive';
        label = 'Out of Stock';
      } else if (stock <= min) {
        variant = 'warning';
        label = `Low Stock (${stock})`;
      } else {
        label = `Optimal (${stock})`;
      }

      return (
        <div className="flex items-center gap-2">
          <Badge variant={variant}>{label}</Badge>
        </div>
      );
    },
  },
  {
    accessorKey: 'min_stock_threshold',
    header: 'Min Threshold',
    cell: ({ row }) => <span className="text-gray-500">{row.original.min_stock_threshold}</span>,
  },
];
