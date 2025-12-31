import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useReactTable,getCoreRowModel, getFilteredRowModel, getSortedRowModel,flexRender } from '@tanstack/react-table';
import api from '../../Api/axios';
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
Select,
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
RoleBadge,
PasswordNote,
ToastContainer,
Toast

 } from '../Users/Users.styles';

import plusIcon from '../Products/assets2/plus.svg'
import trashIcon from '../Products/assets2/trash.svg'
import xIcon from '../Products/assets2/cross.svg'
import searchIcon from '../Products/assets2/search.svg'
import eyeIcon from '../Login/assets/eye-open.svg'
import eyeOffIcon from '../Login/assets/eye-closed.svg'


const Users = () => {
  const [showModal, setShowModal] = useState(false)
  const [users, setUsers] = useState([])
  const [suppliers, setSuppliers] = useState([])
  const [loading, setLoading] = useState(true)

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [supplierID, setSupplierID] = useState('')
  const [showPassword, setShowPassword] = useState(false)

  const [globalFilter, setGlobalFilter] = useState('')

  const [toast, setToast] = useState({ message: '', type: '' });

  const showToast = (message, type = 'success') => {
    setToast({ message, type })
    setTimeout(() => {
    setToast({ message: '', type: '' });
  }, 3000)
};

  const user = getUserInfo()

  const fetchUsers = useCallback(async () => {
    try{
        setLoading(true)
        const res = await api.get('/users/supplier-admin')
        setUsers(res?.data?.users || [])   
    } catch(err){
        console.log('Failed to fetch users: ', err)
        alert(err.response?.data?.error || 'Failed to load users')  
    } finally {
        setLoading(false)
    }
  }, [])

  const fetchSuppliers = useCallback(async () => {
    try{
        setLoading(true)
        const res = await api.get('/suppliers');
        setSuppliers(res?.data?.suppliers || []);
    }catch(err){
        console.log('Failed to fetch suppliers: ', err)
    }

  },[])

  useEffect(() => {
    fetchUsers()
    fetchSuppliers()
  },[fetchUsers, fetchSuppliers])

  const createUser = async (e) => {
    e.preventDefault()
    if (!supplierID){
        alert("PLease select the supplier")
        return
    }

    try{
        await api.post("/users/supplier-admin",{
            name,
            email,
            password,
            supplier_id: parseInt(supplierID)
        })

        setShowModal(false)
        fetchUsers()
        resetForm()
        showToast('Supplier admin created successfully', 'success')

    }catch(err){
        console.error('Error while creating user', err)
        showToast('Failed to create user', 'error');

    }

  }

  const deleteUser = useCallback(async (id) => {
  if (!window.confirm("Are you sure you want to delete this user?")) return;

  try {
    await api.delete(`/users/supplier-admin/${id}`);
    await fetchUsers()
    showToast('User deleted successfully', 'success')
  } catch (err) {
    showToast('Failed to delete user', 'error')
  }
}, [fetchUsers]);

  const resetForm = () => {
    setName('')
    setEmail('')
    setPassword('')
    setSupplierID('')
    setShowPassword(false)
  }

  const getSupplierName = useCallback((supplierId) => {
    const supplier = suppliers.find(s => s.supplier_id === supplierId)
    return supplier ? supplier.company : '-';
  },[suppliers])

  const columns = useMemo(
    () => [
      {
        accessorKey: 'user_id',
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
        accessorKey: 'email',
        header: 'Email',
        cell: (info) => info.getValue(),
      },
      {
        accessorKey: 'role',
        header: 'Role',
        cell: (info) => {
          const role = info.getValue()
          return (
            <RoleBadge $isSystemAdmin={role === 'system_admin'}>
              {role === 'system_admin' ? 'System Admin' : 'Supplier Admin'}
            </RoleBadge>
          );
        },
        size: 150,
      },
      {
        accessorKey: 'supplier_id',
        header: 'Supplier',
        cell: (info) => {
          const supplierId = info.getValue();
          return supplierId ? getSupplierName(supplierId) : '-';
        },
      },
      {
        id: 'actions',
        header: 'Actions',
        cell: ({ row }) => (
          <ActionButtons>
            {row.original.role === 'supplier_admin' && user?.role === 'system_admin' && (
              <IconButton $danger onClick={() => deleteUser(row.original.user_id)}>
                <Icon src={trashIcon} alt="delete" />
              </IconButton>
            )}
          </ActionButtons>
        ),
        size: 100,
      },
    ],
    [user, deleteUser, getSupplierName]
  )

  const table = useReactTable({
    data: users,
    columns,
    state: {
      globalFilter,
    },
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  if (user?.role !== 'system_admin') {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <h2>Access Denied</h2>
        <p>Only system admin can access this page.</p>
      </div>
    );
  }

  return (
      <>
          <PageHeader>
            <Title>Users</Title>
            <SearchWrap>
              <Icon src={searchIcon} alt="search" />
              <SearchIpt
                type="text"
                placeholder="Search Users..."
                value={globalFilter ?? ''}
                onChange={(e) => setGlobalFilter(e.target.value)}
              />
            </SearchWrap>
            <AddButton
              onClick={() => {
                resetForm();
                setShowModal(true);
              }}
            >
              <Icon src={plusIcon} alt="add" />
              Add Supplier Admin
            </AddButton>
          </PageHeader>

          {loading ? (
            <p>Loading users...</p>
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
                        No users found
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
              <ModalTitle>Add New Supplier Admin</ModalTitle>
              <CloseButton onClick={() => setShowModal(false)}>
                <Icon src={xIcon} alt="close" />
              </CloseButton>
            </ModalHeader>

            <Form onSubmit={createUser}>
              <FormGroup>
                <Label>Full Name *</Label>
                <Input
                  type="text"
                  placeholder="Yuvraj Bisht"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Email Address *</Label>
                <Input
                  type="email"
                  placeholder="yuvraj@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Password *</Label>
                <div style={{ position: 'relative' }}>
                  <Input
                    type={!showPassword ? 'text' : 'password'}
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    style={{ paddingRight: '40px' }}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    style={{
                      position: 'absolute',
                      right: '12px',
                      top: '50%',
                      transform: 'translateY(-50%)',
                      background: 'none',
                      border: 'none',
                      cursor: 'pointer',
                      padding: '4px',
                      display: 'flex',
                      alignItems: 'center',
                      color: '#4f6b72',
                    }}
                  >
                    <Icon src={showPassword ? eyeOffIcon : eyeIcon} alt="toggle" />
                  </button>
                </div>
                <PasswordNote>
                  Must be 8+ characters with uppercase, lowercase, number and special character
                </PasswordNote>
              </FormGroup>

              <FormGroup>
                <Label>Assign to Supplier *</Label>
                <Select
                  value={supplierID}
                  onChange={(e) => setSupplierID(e.target.value)}
                  required
                >
                  <option value="">Select Supplier</option>
                  {suppliers.map((supplier) => (
                    <option key={supplier.supplier_id} value={supplier.supplier_id}>
                      {supplier.company} - {supplier.name}
                    </option>
                  ))}
                </Select>
              </FormGroup>

              <PasswordNote style={{ marginBottom: '16px', background: 'rgba(42,123,155,0.1)', padding: '12px', borderRadius: '6px' }}>
                <strong>Note:</strong> This user will only have access to manage products and stock for their assigned supplier company.
              </PasswordNote>

              <ModalActions>
                <CancelButton type="button" onClick={() => setShowModal(false)}>
                  Cancel
                </CancelButton>
                <SubmitButton type="submit">Create User</SubmitButton>
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

export default Users
