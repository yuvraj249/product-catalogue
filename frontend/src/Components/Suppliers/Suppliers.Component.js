import React, { useState, useEffect, useMemo } from 'react';
import { ToastContainer, Toast } from './Suppliers.styles';
import {
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  flexRender,
} from '@tanstack/react-table';
import api from '../../Api/axios';
import { useCallback } from 'react';
import { getUserInfo } from '../../utils/auth';
import {
  PageHeader,
  Title,
  AddButton,
  Modal,
  ModalContent,
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  FormGroup,
  Label,
  Input,
  ModalActions,
  CancelButton,
  SubmitButton,
  Icon,
  SearchWrap,
  SearchIpt,
  TableWrapper,
  Table,
  Th,
  Td,
  ActionButtons,
  IconButton,
} from './Suppliers.styles';



import plusIcon from '../Products/assets2/plus.svg'
import editIcon from '../Products/assets2/edit.svg'
import trashIcon from '../Products/assets2/trash.svg'
import xIcon from '../Products/assets2/cross.svg'
import searchIcon from '../Products/assets2/search.svg'

const Suppliers = () => {
  const [showModal, setShowModal] = useState(false)

  const [suppliers, setSuppliers] = useState([])
  const [loading, setLoading] = useState(true)

  const [name, setName] = useState('')
  const [company, setCompany] = useState('')
  const [email, setEmail] = useState('')
  const [contact, setContact] = useState('')

  const [updatingId, setUpdatingId] = useState(null)
  const isUpdating = updatingId !== null

  const [globalFilter, setGlobalFilter] = useState('');

  const [toast, setToast] = useState({ message: '', type: '' })

  const showToast = (message, type = 'success') => {
  setToast({ message, type });

  setTimeout(() => {
    setToast({ message: '', type: '' });
  }, 3000);
}

  const user = getUserInfo()

  const fetchSuppliers = useCallback(async () => {
  try {
    setLoading(true)
    const res = await api.get('/suppliers')
    setSuppliers(res.data.suppliers || [])
  } catch (err) {
    console.log('failed to fetch suppliers: ', err)
    alert(err.response?.data?.error || 'Failed to load suppliers')
  } finally {
    setLoading(false)
  }
}, [])

  useEffect(() => {
    fetchSuppliers()
  }, [fetchSuppliers])
 
  const createSupplier = async (e) => {
    e.preventDefault()
    try {
      await api.post('/suppliers', {
        name,
        contact_info: contact,
        email,
        company,
      })
      showToast('Supplier created successfully', 'success')
      setShowModal(false)
      fetchSuppliers()
      setName('')
      setCompany('')
      setEmail('')
      setContact('')
    } catch (err) {
      console.error('Create supplier error:', err)
      showToast('Failed to create supplier', 'error');
    }
  }

  const startUpdating = useCallback((supplier) => {
  setUpdatingId(supplier.supplier_id)
  setName(supplier.name || '')
  setEmail(supplier.email || '')
  setCompany(supplier.company || '')
  setContact(supplier.contact_info || '')
  setShowModal(true)
}, [])

  const updateSupplier = async (e) => {
    e.preventDefault()
    try {
      await api.put(`/suppliers/${updatingId}`, {
        name,
        email,
        company,
        contact_info: contact,
      });
      showToast('Supplier updated successfully', 'success');
      setShowModal(false)
      setUpdatingId(null)
      fetchSuppliers()
    } catch (err) {
      console.error('Update supplier error:', err)
      showToast('Failed to update supplier', 'error');
    }
  }

  const deleteSupplier = useCallback(async (id) => {
  try {
    await api.delete(`suppliers/${id}`)
    showToast('Supplier deleted successfully', 'success');
    fetchSuppliers()
  } catch (err) {
    showToast('Failed to delete supplier', 'error');
  }
}, [fetchSuppliers])

  const columns = useMemo(
    () => [
      {
        accessorKey: 'supplier_id',
        header: 'ID',
        cell: (info) => `#${info.getValue()}`,
        size: 80,
      },
      {
        accessorKey: 'name',
        header: 'Name',
        cell: (info) => info.getValue(),
      },
      {
        accessorKey: 'company',
        header: 'Company',
        cell: (info) => info.getValue(),
      },
      {
        accessorKey: 'email',
        header: 'Email',
        cell: (info) => info.getValue(),
      },
      {
        accessorKey: 'contact_info',
        header: 'Contact',
        cell: (info) => info.getValue(),
      },
      {
        id: 'actions',
        header: 'Actions',
        cell: ({ row }) => (
          <ActionButtons>
            <IconButton onClick={() => startUpdating(row.original)}>
              <Icon src={editIcon} alt="edit" />
            </IconButton>
            {user?.role === 'system_admin' && (
              <IconButton $danger onClick={() => deleteSupplier(row.original.supplier_id) }>
                <Icon src={trashIcon} alt="delete" />
              </IconButton>
            )}
          </ActionButtons>
        ),
        size: 120,
      },
    ],
    [user, startUpdating, deleteSupplier]
  );

  const table = useReactTable({
    data: suppliers,
    columns,
    state: {
      globalFilter,
    },
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
        <>
          <PageHeader>
            <Title>Suppliers</Title>
            <SearchWrap>
              <Icon src={searchIcon} alt="search" />
              <SearchIpt
                type="text"
                placeholder="Search Suppliers..."
                value={globalFilter ?? ''}
                onChange={(e) => setGlobalFilter(e.target.value)}
              />
            </SearchWrap>
            {user?.role === 'system_admin' && (
              <AddButton
                onClick={() => {
                  setUpdatingId(null)
                  setName('')
                  setEmail('')
                  setCompany('')
                  setContact('')
                  setShowModal(true)
                }}
              >
                <Icon src={plusIcon} alt="add" />
                Add Supplier
              </AddButton>
            )}
          </PageHeader>

          {loading ? (
            <p>Loading suppliers...</p>
          ) : (
            <TableWrapper>
              <Table>
                <thead>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <tr key={headerGroup.id}>
                      {headerGroup.headers.map((header) => (
                        <Th
                          key={header.id}
                          align={header.id === 'actions' ? 'center' : 'left'}
                        >
                          {header.isPlaceholder
                            ? null
                            : flexRender(
                                header.column.columnDef.header,
                                header.getContext()
                              )}
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
                        No suppliers found
                      </Td>
                    </tr>
                  )}
                </tbody>
              </Table>
            </TableWrapper>
          )}

      {showModal && (
        <Modal onClick={() => setShowModal(false)}>
          <ModalContent onClick={(e) => e.stopPropagation()}>
            <ModalHeader>
              <ModalTitle>
                {isUpdating ? 'Edit Supplier' : 'Add New Supplier'}
              </ModalTitle>
              <CloseButton onClick={() => setShowModal(false)}>
                <Icon src={xIcon} alt="close" />
              </CloseButton>
            </ModalHeader>

            <Form onSubmit={isUpdating ? updateSupplier : createSupplier}>
              <FormGroup>
                <Label>Name *</Label>
                <Input
                  type="text"
                  placeholder="Enter supplier name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Email *</Label>
                <Input
                  type="email"
                  placeholder="supplier@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Company *</Label>
                <Input
                  type="text"
                  placeholder="Company name"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Contact Info *</Label>
                <Input
                  type="text"
                  placeholder="+91-"
                  value={contact}
                  onChange={(e) => setContact(e.target.value)}
                  required
                />
              </FormGroup>

              <ModalActions>
                <CancelButton type="button" onClick={() => setShowModal(false)}>
                  Cancel
                </CancelButton>
                <SubmitButton type="submit">
                  {isUpdating ? 'Update' : 'Create'}
                </SubmitButton>
              </ModalActions>
            </Form>
          </ModalContent>
        </Modal>
      )}
      <ToastContainer>
        {
            toast.message && (
                <Toast $type={toast.type}>
                    {toast.message}
                </Toast>
            )
        }
      </ToastContainer>
  </>

  )
}

export default Suppliers