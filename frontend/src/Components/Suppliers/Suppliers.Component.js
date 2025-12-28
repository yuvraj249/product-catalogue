import React, { useState, useEffect, useMemo } from 'react';
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
  Layout,
  Sidebar,
  SidebarHeader,
  Main,
  TopBar,
  MenuButton,
  UserRole,
  ContentArea,
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
import { useNavigate } from 'react-router-dom';

import packageIcon from '../Dashboard/assets1/package.svg'
import menuIcon from '../Dashboard/assets1/menu.svg'
import plusIcon from '../Products/assets2/plus.svg'
import editIcon from '../Products/assets2/edit.svg'
import trashIcon from '../Products/assets2/trash.svg'
import xIcon from '../Products/assets2/cross.svg'
import folderIcon from '../Dashboard/assets1/folder-tree.svg'
import dashboardIcon from '../Dashboard/assets1/layout-dashboard.svg'
import usersIcon from '../Dashboard/assets1/users.svg'
import truckIcon from '../Dashboard/assets1/truck.svg'
import trendingIcon from '../Dashboard/assets1/trending-up.svg'
import logoutIcon from '../Dashboard/assets1/log-out.svg'
import searchIcon from '../Products/assets2/search.svg'

import { Nav, NavItem, LogoutWrapper, LogoutButton } from '../Dashboard/Dashboard.styles'

const Suppliers = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [activeItem, setActiveItem] = useState('Suppliers')

  const [suppliers, setSuppliers] = useState([])
  const [loading, setLoading] = useState(true)

  const [name, setName] = useState('')
  const [company, setCompany] = useState('')
  const [email, setEmail] = useState('')
  const [contact, setContact] = useState('')

  const [updatingId, setUpdatingId] = useState(null)
  const isUpdating = updatingId !== null

  const [globalFilter, setGlobalFilter] = useState('');

  const navigate = useNavigate()
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
      setShowModal(false)
      fetchSuppliers()
      setName('')
      setCompany('')
      setEmail('')
      setContact('')
    } catch (err) {
      console.error('Create supplier error:', err)
      alert(err.response?.data?.error || 'Failed to create supplier')
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
      setShowModal(false)
      setUpdatingId(null)
      fetchSuppliers()
    } catch (err) {
      console.error('Update supplier error:', err)
      alert(err.response?.data?.error || 'Failed to update supplier')
    }
  }

  const deleteSupplier = useCallback(async (id) => {
  if (!window.confirm('Are you sure you want to delete this supplier?')) return
  try {
    await api.delete(`suppliers/${id}`)
    fetchSuppliers()
  } catch (err) {
    alert(err.response?.data?.error || 'Failed to delete supplier')
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
              <IconButton $danger onClick={() => deleteSupplier(row.original.supplier_id)}>
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

  const menuItems = [
    { icon: dashboardIcon, label: 'Dashboard' },
    { icon: packageIcon, label: 'Products' },
    { icon: folderIcon, label: 'Categories' },
    { icon: truckIcon, label: 'Suppliers' },
    { icon: trendingIcon, label: 'Stock Movements' },
    { icon: usersIcon, label: 'Users' },
  ];

  const logout = () => {
    navigate('/');
  };

  return (
    <Layout>
      <Sidebar $isOpen={sidebarOpen}>
        <SidebarHeader>
          <Icon src={packageIcon} alt="logo" />
          Product Catalogue
        </SidebarHeader>
        <Nav>
          {menuItems.map((item) => (
            <NavItem
              key={item.label}
              $active={activeItem === item.label}
              onClick={() => setActiveItem(item.label)}
            >
              <Icon src={item.icon} alt={item.label} />
              <span>{item.label}</span>
            </NavItem>
          ))}
        </Nav>
        <LogoutWrapper>
          <LogoutButton>
            <Icon src={logoutIcon} alt="logout" />
            <span onClick={logout}>Logout</span>
          </LogoutButton>
        </LogoutWrapper>
      </Sidebar>

      <Main $sidebarOpen={sidebarOpen}>
        <TopBar>
          <MenuButton onClick={() => setSidebarOpen(!sidebarOpen)}>
            <Icon src={menuIcon} alt="menu" />
          </MenuButton>
          <UserRole>{user?.role || 'User'}</UserRole>
        </TopBar>

        <ContentArea>
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
        </ContentArea>
      </Main>

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
    </Layout>
  )
}

export default Suppliers