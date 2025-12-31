import React, { useState, useEffect, useCallback, useMemo} from 'react';
import api from '../../Api/axios';
import {
  PageHeader,
  Title,
  StatsGrid,
  StatCard,
  StatIcon,
  StatInfo,
  StatValue,
  StatLabel,
  Section,
  SectionHeader,
  SectionTitle,
  ThresholdBadge,
  Table,
  Th,
  Td,
  TableWrapper,
  Icon
} from './Dashboard.styles'
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  flexRender,
} from '@tanstack/react-table';

import packageIcon from './assets1/package.svg'
import folderIcon from './assets1/folder-tree.svg'
import usersIcon from './assets1/users.svg'
import truckIcon from './assets1/truck.svg'
import alertIcon from './assets1/alert-triangle.svg'

export const Dashboard = ()  => {
  const [supplierCount, setSupplierCount] = useState(0)
  const [productCount, setProductCount] = useState(0)
  const [categoryCount, setCategoryCount] = useState(0)
  const [userCount, setUserCount] = useState(0)


  const [movements, setMovements] = useState([]);
  const [globalFilter, setGlobalFilter] = useState('');

  useEffect(() => {
    const fetchSupplierCount = async () => {
    try {
        const supRes = await api.get('/suppliers')
        setSupplierCount(supRes?.data?.suppliers?.length || 0)
    } catch (err) {
        console.error("Failed to fetch supplier count", err)
    }
  }
  fetchSupplierCount();
  }, [])

  useEffect(() => {
    const fetchCategoryCount = async () => {
      try{
        const catgRes = await api.get('/categories')
        setCategoryCount(catgRes?.data?.categories?.length || 0)
      }catch (err) {
        console.error("Failed to fetch category count", err)
      }
    }
    fetchCategoryCount()
  },[])

  useEffect(() => {
    const fetchProductCount = async () => {
      try{
        const prodRes = await api.get('/products')
        setProductCount(prodRes?.data?.products?.length || 0)
      } catch (err) {
        console.error("Failed to fetch product count", err)
      }
    }
    fetchProductCount()
  },[])

  useEffect(() => {
    const fetchUserCount = async () => {
      try{
        const userRes = await api.get('/users/supplier-admin')
        const userCount = userRes?.data?.users?.filter(u => u.role === 'supplier_admin').length || 0
        setUserCount(userCount)
      } catch (err) {
        console.error("Failed to fetch user count", err)
      }
    }
    fetchUserCount()
  },[])

  const fetchMovements = useCallback(async () => {
    try {
      const res = await api.get('/stock_movements');
      setMovements(res?.data?.movements || []);
    } catch (err) {
      console.error("Failed to fetch stock movements", err);
    }
  }, []);

  useEffect(() => {
    fetchMovements();
  }, [fetchMovements]);

  const stats = [
    { icon: packageIcon, label: 'Total Products', value: productCount, color: '#8b7355' },
    { icon: folderIcon, label: 'Categories', value: categoryCount, color: '#7a9b76' },
    { icon: truckIcon, label: 'Suppliers', value: supplierCount, color: '#7a8b9b' },
    { icon: usersIcon, label: 'Supplier Admins', value: userCount, color: '#c9a86a' },
  ]

  const columns = useMemo(() => [
    {
      accessorKey: 'stock_id',
      header: 'Stock ID',
      cell: info => `#${info.getValue()}`
    },
    {
      accessorKey: 'product_id',
      header: 'Product ID',
    },
    {
      accessorKey: 'quantity',
      header: 'Quantity',
    },
    {
      accessorKey: 'movement_type',
      header: 'Movement Type',
      cell: info => info.getValue() === 'IN' ? 'IN' : 'OUT'
    },
    {
      accessorKey: 'reason',
      header: 'Reason',
      cell: info => info.getValue() || '-'
    },
    {
      accessorKey: 'performed_by',
      header: 'Performed By',
      cell: info => info.getValue() || '-'
    }
  ], [])

  const table = useReactTable({
    data: movements,
    columns,
    state: {
      globalFilter
    },
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  return (
    <>
      <PageHeader>
            <Title>Dashboard</Title>
      </PageHeader>

      <StatsGrid>
            {stats.map((stat, index) => (
              <StatCard key={index}>
                <StatIcon $color={stat.color}>
                  <Icon src={stat.icon} alt={stat.label} />
                </StatIcon>
                <StatInfo>
                  <StatValue>{stat.value}</StatValue>
                  <StatLabel>{stat.label}</StatLabel>
                </StatInfo>
              </StatCard>
            ))}
        </StatsGrid>

      <Section>
          <SectionHeader>
              <SectionTitle>
                <Icon src={alertIcon} alt="alert" />
                Low Stock Alert
              </SectionTitle>
              <ThresholdBadge>Threshold: 10 units</ThresholdBadge>
          </SectionHeader>

      <TableWrapper>
        <Table>
          <thead>
            {table.getHeaderGroups().map(headerGroup => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map(header => (
                  <Th key={header.id}>
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </Th>
                ))}
              </tr>
            ))}
          </thead>

          <tbody>
            {table.getRowModel().rows.map(row => (
              <tr key={row.id}>
                {row.getVisibleCells().map(cell => (
                  <Td key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </Td>
                ))}
              </tr>
            ))}

            {table.getRowModel().rows.length === 0 && (
              <tr>
                <Td colSpan={columns.length} style={{ textAlign: 'center' }}>
                  No stock movements found
                </Td>
              </tr>
            )}
          </tbody>
        </Table>
      </TableWrapper>
    </Section>
  </>
  );
}

export default Dashboard