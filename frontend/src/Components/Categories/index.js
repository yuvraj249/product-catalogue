import React, {useState, useEffect, useCallback, useMemo} from 'react'
import api from '../../Api/axios'
import { getUserInfo } from '../../utils/auth'
import {
    PageHeader,
    Title,
    AddButton,
    Icon, 
    SearchWrap,
    SearchIpt,
    IconButton,
    ActionButtons,
} from './Styles'
import {toast,ToastContainer} from 'react-toastify'
import plusIcon from '../../Images/plus.svg'
import searchIcon from '../../Images/search.svg'
import editIcon from '../../Images/edit.svg'
import trashIcon from '../../Images/trash.svg'
import Datatable from '../DataTable';
import AddorEditModal from './AddorEditModal'
import { useLocation } from 'react-router-dom'


const Categories = () => { 
   const user = getUserInfo();
   const [loading, setLoading] = useState(true)
   const [globalFilter, setGlobalFilter] = useState("")
   const [modalOpen, setModalOpen] = useState(false)
   const [state, setState] = useState({
    data: [],
    updatingId: null,
    form: { name: "", description: "" },
  });
  

const fetchCategories = useCallback(async (q = "") => {
    try {
      setLoading(true)
      const res = await api.get("/categories", {
        params: { q }
      })
      setState(p => ({ ...p, data: res.data.categories || [] }))
    } catch {
      toast.error("Failed to load categories")
    } finally {
      setLoading(false)
    }
  }, [])

const location = useLocation()
useEffect(() => {
    const getFromBrowser = new URLSearchParams(window.location.search)
    const q = getFromBrowser.get("q")
    if (q !== null) {
        setGlobalFilter(q)
    }
}, [location.search])

useEffect(() => {
    const handler = setTimeout(() => {
      const value = globalFilter.trim()

      fetchCategories(value)
    }, 350)

    return () => clearTimeout(handler)
  }, [globalFilter, fetchCategories])


const deleteCategory = useCallback(async (id) => {
    try{
        await api.delete(`/categories/${id}`)
        toast.success('Category deleted successfully')
        fetchCategories()
    }catch{
       toast.error('Failed to delete category') 
    }
   },[fetchCategories])


  const onClickEdit = useCallback((row) => {
    setState((prev) => ({
        ...prev,
        updatingId: row.original.category_id,
        form: {
            name: row.original.category_name,
            description: row.original.category_description,
        },
    }));
    setModalOpen(true)
  },[])

  const onClickAddCatg = () => {
        setState((prev) => ({
            ...prev,
            updatingId: null,
            form: { name: "", description: "" },
      }))
        setModalOpen(true)
  }

  const onModalClose = () => {
        setModalOpen(false)
        setState(prev => ({...prev, updatingId: null}))
   }

   const columns = useMemo(() => 
    [
        {
            accessorKey: 'category_id',
            header: 'ID',
            cell: (info) => `#${info.getValue()}`
        },
        {
            accessorKey: 'category_name',
            header: 'Name',
        },
        {
            accessorKey: 'category_description',
            header: 'Description',
        },
           ...(user.role === "system_admin" ?
            [ {
           id: 'actions',
           header: 'Actions', 
           cell: ({row} ) => 
            (
                <ActionButtons>
                <IconButton onClick={() => onClickEdit(row) }><Icon src={editIcon} alt="edit" /></IconButton>
                <IconButton onClick={() => deleteCategory(row.original.category_id)}>
                    <Icon src={trashIcon} alt="delete" />
                </IconButton>
                </ActionButtons>
            ) 
        },
    ] : [])
    ], [user.role, onClickEdit, deleteCategory])


   return (
            <>
                <PageHeader>
                    <Title>Categories</Title>
                        <SearchWrap>
                            <Icon src={searchIcon} alt='search'/>
                            <SearchIpt 
                              value={globalFilter}
                              onChange={(e) => setGlobalFilter(e.target.value)}
                              placeholder='Search Categories'
                             />
                        </SearchWrap>
                        {
                            user.role === 'system_admin' && (
                                <AddButton 
                                onClick={onClickAddCatg}>
                                    <Icon src={plusIcon} alt="add" /> 
                                    Add Category 
                                </AddButton>
                            )
                        }
                </PageHeader>
                {loading ? (
                    <p style={{ padding: "20px" }}>Loading categories…</p>
                ) : (
                <Datatable 
                data={state.data}
                columns={columns}         
                />
                )}
         <AddorEditModal
            open={modalOpen}
            onClose={onModalClose}
            updatingId={state.updatingId}
            refetch={fetchCategories}
            setLoading={setLoading}
        />                    
        <ToastContainer position="top-right" autoClose={3000}/>        
    </>
   )}

export default Categories

