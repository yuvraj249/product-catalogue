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
  

const fetchCategories = useCallback(async () => {
    try{
       setLoading(true)
        const res = await api.get("/categories")
        setState((prev) => ({...prev, data: res.data.categories || []}))     
    }catch(err){
        console.log('failed to fetch categories: ', err)
        toast.error("Failed to load categories")
    } finally{
        setLoading(false)
    }
   }, [])

   useEffect(() => {
    fetchCategories()
   }, [fetchCategories])



const location = useLocation()
useEffect(() => {
    const getFromBrowser = new URLSearchParams(window.location.search)
    const q = getFromBrowser.get("q")
    if (q !== null) {
        setGlobalFilter(q)
    }
}, [location.search])

useEffect(() => {
  const handler = setTimeout(async () => {
    const search = globalFilter.trim()

    switch (true) {
        case search === "":
            fetchCategories()
            return

        case /^\d+$/.test(search):

            try {
                const res = await api.get(`/categories/${search}`)
                setState(prev => ({
                ...prev,
                data: [res.data.category]
            }))
          } catch {
                setState(prev => ({ ...prev, data: [] }))
        }
        return

        default:
            try {
                const res = await api.get("/categories", {
                params: { name: search, description: search }
            })

                setState(prev => ({
                    ...prev,
                    data: res.data.categories || []
               }))
              } catch (err) {
                    console.log(err.response?.data || err.message)
                    toast.error("Failed to search categories")
        }
    }
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
                globalFilter={globalFilter}
                setGlobalFilter={setGlobalFilter}          
                />
                )}
         <AddorEditModal
            open={modalOpen}
            onClose={onModalClose}
            updatingId={state.updatingId}
            refetch={fetchCategories}
        />                    
        <ToastContainer position="top-right" autoClose={3000}/>        
    </>
   )}

export default Categories

