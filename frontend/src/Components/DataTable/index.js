import { TableWrapper , Table, Th, Td, SortWrapper, Arrow, HeaderContent } from './Styles';
import {useReactTable, getCoreRowModel, getSortedRowModel, flexRender} from "@tanstack/react-table";


const initialProps = {
  data: [],
  columns: [],
}

const Datatable = (props = initialProps) => {

    const {data, columns} = props
    const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
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
                                 <HeaderContent $center={header.id === "actions"}>
                                {flexRender(
                                    header.column.columnDef.header,
                                    header.getContext()
                                )}
                                {
                                    header.column.getCanSort() && (
                                        <SortWrapper>
                                            <Arrow $active={header.column.getIsSorted() === "asc"} onClick={() => header.column.toggleSorting(false)}>▲</Arrow>
                                            <Arrow $active={header.column.getIsSorted() === "desc"} onClick={() => header.column.toggleSorting(true)}>▼</Arrow>
                                        </SortWrapper>
                                    )
                                }
                                </HeaderContent>
                            </Th>
                        ))}
                    </tr>
                ))}
            </thead>

            <tbody>
                {table.getRowModel().rows.length > 0 ? (
                    table.getRowModel().rows.map((row) => (
                      <tr key={row.id}>
                        {row.getVisibleCells().map((cell) => (
                          <Td
                            key={cell.id}
                            align={cell.column.id === 'actions' ? 'center' : 'left'}
                          >
                            {flexRender(
                              cell.column.columnDef.cell,
                              cell.getContext()
                            )}
                          </Td>
                        ))}
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <Td colSpan={columns.length} style={{ textAlign: 'center' }}>
                        No records found
                      </Td>
                    </tr>
                  )}
            </tbody>
        </Table>
    </TableWrapper>
  )
}

export default Datatable

