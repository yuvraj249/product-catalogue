import React from "react";
import { TableWrapper , Table, Th, Td } from './Styles';
import {useReactTable, getCoreRowModel, getFilteredRowModel, getSortedRowModel, flexRender} from "@tanstack/react-table";

const Datatable = ({data=[], columns, globalFilter, setGlobalFilter}) => {
    const table = useReactTable({
    data,
    columns,
    state: { globalFilter },
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel()
  })


  return (
    <TableWrapper>
        <Table>
            <thead>
                {table.getHeaderGroups().map((hg) => (
                    <tr key={hg.id}>
                        {hg.headers.map((header) => (
                            <Th key={header.id} align={header.id === "actions" ? "center" : "left"}>
                                {flexRender(
                                    header.column.columnDef.header,
                                    header.getContext()
                                )}
                            </Th>
                        ))}
                    </tr>
                ))}
            </thead>

            <tbody>
                {
                    table.getRowModel().rows.map((row) => (
                        <tr key={row.id}>
                            {
                                row.getVisibleCells().map((cell) => (
                                    <Td key={cell.id} align={cell.column.id === "actions" ? "center" : "left"}>
                                        {flexRender(
                                            cell.column.columnDef.cell,
                                            cell.getContext()
                                        )}
                                    </Td>
                                ))
                            }
                        </tr>
                    ))
                }
            </tbody>
        </Table>
    </TableWrapper>
  )
}

export default Datatable

