import React, {useState, useEffect, useMemo, useCallback} from 'react'
import { useReactTable, getCoreRowModel, getFilteredRowModel, getSortedRowModel, flexRender} from '@tanstack/react-table'
import api from '../../Api/axios'
import { getUserInfo } from '../../utils/auth'
import {
    PageHeader,
    Title,
    AddButton,
    Modal,
    ModalContent,
    ModalHeader,
    ModalTitle,
    CloseButton,
    FormGroup,
    Form,
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
} from './Categories.styles'

import plusIcon from '../Products/assets2/plus.svg'
import editIcon from '../Products/assets2/edit.svg'
import trashIcon from '../Products/assets2/trash.svg'
import xIcon from '../Products/assets2/cross.svg'
import searchIcon from '../Products/assets2/search.svg'
import { Toast, ToastContainer } from '../Suppliers/Suppliers.styles'

const Categories = () => {
    
    const [showModal, setShowModal] = useState(false)

    const [categories , setCategories] = useState([])
    const [loading, setLoading] = useState(true)

    const [categoryName, setCategoryName] = useState('')
    const [catetoryDesp, setCatgDesp] = useState('')

    const [updatingId, setUpdatingId] = useState(null)
    const isUpdating = updatingId !== null

    const [globalFilter, setGlobalFilter] = useState('')

    const [toast, setToast] = useState({ message: '', type: '' })

    const showToast = (message, type = 'success') => {
        setToast({ message, type })
        setTimeout(() => {
            setToast({ message: '', type: '' })
        }, 3000)
   }
   
   const user = getUserInfo()

   const fetchCategories = useCallback(async () => {
    try{
        setLoading(true)
        const res = await api.get("/categories")
        setCategories(res?.data?.categories || [])      
    }catch(err){
        console.log('failed to fetch categories: ', err)
        alert(err.response?.data?.error || 'Failed to load categories')
    }finally{
        setLoading(false)
    }
   }, [])

   useEffect(() => {
    fetchCategories()
   }, [fetchCategories])

   const createCategory = async (e) => {
    e.preventDefault()
    try{
           await api.post('/categories',{
           category_name: categoryName,
           category_description: catetoryDesp 
        })
        showToast('Category created successfully', 'success')
        setShowModal(false)
        fetchCategories()
        setCategoryName('')
        setCatgDesp('')    
    } catch(err){
        console.error('Error while creating category:', err)
        showToast('Failed to create category', 'error')

    }
   }
   
   const startUpdating = useCallback((category) => {
     setUpdatingId(category.category_id)
     setCategoryName(category.category_name || '')
     setCatgDesp(category.category_description || '')
     setShowModal(true)
   }, [])

   const updateCategory = async (e) => {
    e.preventDefault()
    try{
        await api.put(`/categories/${updatingId}`,{
        category_name:categoryName,
        category_description:catetoryDesp
    })
    setShowModal(false)
    setUpdatingId(null)
    fetchCategories()
    }catch(err){
        console.error('Error while updating category:', err)
        showToast('Failed to update category', 'error')
    }
   }

   const deleteCategory = useCallback(async (id) => {
    try{
        await api.delete(`/categories/${id}`)
        showToast('Category deleted successfully', 'success')
        fetchCategories()
    }catch{
       showToast('Failed to delete category', 'error') 
    }
   }, [fetchCategories])

   const columns = useMemo(
    () => [
        {
            accessorKey: 'category_id',
            header: 'ID',
            cell: (info) => `#${info.getValue()}`,
            size: 80,
        },
        {
            accessorKey: 'category_name',
            header: 'Name',
            cell: (info) => info.getValue(),
        },
        {
            accessorKey: 'category_description',
            header: 'Description',
            cell: (info) => info.getValue(),
        },
        {
           id: 'actions',
           header: 'Actions', 
            cell: ({ row }) => (
                <ActionButtons>
                    {user?.role === 'system_admin' && (
                        <IconButton onClick={() => startUpdating(row.original)}>
                            <Icon src={editIcon} alt="edit" />
                        </IconButton>
                    )}
                    {
                      user?.role === 'system_admin' && (
                        <IconButton onClick={() => deleteCategory(row.original.category_id)}>
                            <Icon src={trashIcon} alt="delete" />
                        </IconButton>
                      )  
                    }
                </ActionButtons>
            ),
            size: 120
        },
    ],
    [user, startUpdating,  deleteCategory]
   )

   const table = useReactTable({
    data: categories,
    columns,
    state: {
      globalFilter,
    },
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel()
   })
   return (
            <>
                <PageHeader>
                    <Title>Categories</Title>
                        <SearchWrap>
                            <Icon src={searchIcon} alt='search'/>
                            <SearchIpt 
                              type='text'
                              placeholder='Search Categories'
                              value={globalFilter ?? ''}
                              onChange={(e) => setGlobalFilter(e.target.value)}
                             />
                        </SearchWrap>
                        {
                            user?.role === 'system_admin' && (
                                <AddButton 
                                onClick={() => {
                                    setUpdatingId(null)
                                    setCategoryName('')
                                    setCatgDesp('')
                                    setShowModal(true)
                                 }}
                                >
                                    <Icon src={plusIcon} alt="add" /> 
                                    Add Category 
                                   </AddButton>
                            )
                        }
                    </PageHeader>

                    {loading ? (
                        <p>Loading categories...</p>
                    ): (
                        <TableWrapper>
                            <Table>
                                <thead>
                                    {
                                        table.getHeaderGroups().map((headerGroup)=>(
                                           <tr key={headerGroup.id}>
                                            {
                                                headerGroup.headers.map((header) => (
                                                    <Th key={header.id} align={header.id === 'actions' ? 'center' : 'left'}>
                                                        {header.isPlaceholder ? null : flexRender(header.column.columnDef.header,header.getContext())}
                                                    </Th>
                                                ))
                                            }
                                           </tr> 
                                        ))
                                    }
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
                                                        No Categories found
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
                                {isUpdating? 'Edit Category': 'Add new Category'}
                            </ModalTitle>
                            <CloseButton onClick={() => setShowModal(false)}>
                                <Icon src={xIcon} alt="close" />
                            </CloseButton>
                        </ModalHeader>

                        <Form onSubmit={isUpdating ? updateCategory: createCategory}>
                            <FormGroup>
                                <Label>Name *</Label>
                                <Input
                                  type='text'
                                  placeholder='Enter Category name'
                                  value={categoryName}
                                  onChange={(e) => setCategoryName(e.target.value)}
                                  required     
                                />
                            </FormGroup>

                            <FormGroup>
                                <Label>Description *</Label>
                                <Input 
                                  type='text'
                                  placeholder='Enter Category description'
                                  value={catetoryDesp}
                                  onChange={(e) => setCatgDesp(e.target.value)}
                                />
                            </FormGroup>

                            <ModalActions>
                                <CancelButton type="button" onClick={() => setShowModal(false)}>Cancel</CancelButton>
                                <SubmitButton type='submit'>{isUpdating ? 'Update': 'Create'}</SubmitButton>
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
   )}

export default Categories

