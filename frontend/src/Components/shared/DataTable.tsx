import React, { useState } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  ColumnDef,
  RowSelectionState,
} from '@tanstack/react-table';
import { useDebounce } from '../../hooks/useDebounce';

interface DataTableProps<TData> {
  columns: ColumnDef<TData, any>[];
  data: TData[];
  pageCount: number;
  onPaginationChange: (page: number, pageSize: number) => void;
  onSearchChange: (search: string) => void;
  onBulkScrap?: (selectedRows: TData[]) => void;
}

export function DataTable<TData>({
  columns,
  data,
  pageCount,
  onPaginationChange,
  onSearchChange,
  onBulkScrap,
}: DataTableProps<TData>) {
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});
  const [globalFilter, setGlobalFilter] = useState('');

  const debouncedSearch = useDebounce(globalFilter, 400);

  React.useEffect(() => {
    onSearchChange(debouncedSearch);
  }, [debouncedSearch, onSearchChange]);

  const table = useReactTable({
    data,
    columns,
    pageCount,
    state: { rowSelection },
    enableRowSelection: true,
    onRowSelectionChange: setRowSelection,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    manualFiltering: true,
  });

  const selectedItems = table.getSelectedRowModel().rows.map((r) => r.original);

  return (
    <div className="w-full space-y-4">
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-slate-900 p-4 rounded-xl border border-slate-800 shadow-sm">
        <div className="relative w-full sm:w-96">
          <input
            type="text"
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Filter SKU, Name, Barcode..."
            className="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all"
          />
        </div>

        {selectedItems.length > 0 && (
          <div className="flex items-center gap-3 w-full sm:w-auto justify-end">
            <span className="text-xs font-medium text-slate-300 bg-slate-800 px-3 py-1.5 rounded-full border border-slate-700">
              {selectedItems.length} items selected
            </span>
            <button
              onClick={() => onBulkScrap && onBulkScrap(selectedItems)}
              className="px-3.5 py-2 bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white text-xs font-semibold rounded-lg shadow-sm transition-all focus:outline-none focus:ring-2 focus:ring-rose-500"
            >
              Bulk Move to Scrap
            </button>
          </div>
        )}
      </div>

      <div className="rounded-xl border border-slate-800 bg-slate-950 overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-900/90 border-b border-slate-800 text-xs font-semibold text-slate-400 uppercase tracking-wider">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <th key={header.id} className="px-4 py-3.5">
                      {header.isPlaceholder
                        ? null
                        : flexRender(header.column.columnDef.header, header.getContext())}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody className="divide-y divide-slate-800/60 bg-slate-950">
              {table.getRowModel().rows.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="px-4 py-12 text-center text-slate-500 text-sm">
                    No inventory records found matching your filter criteria.
                  </td>
                </tr>
              ) : (
                table.getRowModel().rows.map((row) => (
                  <tr key={row.id} className="hover:bg-indigo-950/20 transition-colors">
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id} className="px-4 py-3.5 text-slate-200 font-normal">
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
