import React, {useState, useEffect, useCallback} from 'react'
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
    CancelButton,
    Form, 
    FormGroup, 
    Input, 
    Label,
    CloseButton,
    ModalActions, 
    ModalHeader, 
    ModalTitle, 
    SubmitButton, 
    Textarea
} from './Styles'
import {toast,ToastContainer} from 'react-toastify'
import ModalBox from '../ModalBox'
import plusIcon from '../../Images/plus.svg'
import searchIcon from '../../Images/search.svg'
import editIcon from '../../Images/edit.svg'
import trashIcon from '../../Images/trash.svg'
import xIcon from '../../Images/cross.svg'
import Datatable from '../DataTable/index'


const Categories = () => { 
   const user = getUserInfo();
   const [loading, setLoading] = useState(true)
   const [modalOpen, setModalOpen] = useState(false)
   const [globalFilter, setGlobalFilter] = useState("")
   const [state, setState] = useState({
    data: [],
    updatingId: null,
    form: { name: "", description: "" },
  })

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



   const createCategory = async () => {
    try{
        await api.post('/categories',{
            category_name: state.form.name,
            category_description: state.form.description
        })
        toast.success("Category created successfully")
        setModalOpen(false)
        setState((prev) => ({
            ...prev,
            updatingId: null,
            form: { name: "", description: "" }
        }))
        fetchCategories()
    }catch (err){
        const msg = err.response.data.error || ""
        toast.error(msg || "Failed to create or edit category")

    }
   }


    const updateCategory = async () => {
     try {
        await api.put(`/categories/${state.updatingId}`, {
        category_name: state.form.name,
        category_description: state.form.description
    });

    toast.success("Category updated successfully");

    setModalOpen(false);

    setState(prev => ({
      ...prev,
      updatingId: null,
      form: { name: "", description: "" }
    }));

    fetchCategories();
    } catch (err) {
        const msg = err.response.data.error || "";
        toast.error(msg || "Failed to update category");
    }
   }



   const deleteCategory = async (id) => {
    try{
        await api.delete(`/categories/${id}`)
        toast.success('Category deleted successfully')
        fetchCategories()
    }catch{
       toast.error('Failed to delete category') 
    }
   }


   const handleSubmit = () => {
    state.updatingId ? updateCategory() : createCategory();
  }


   const columns = 
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
        {
           id: 'actions',
           header: 'Actions', 
           cell: ({row} ) => 
            user.role === "system_admin" && (
                <ActionButtons>
                <IconButton onClick={() => {setState((prev) => ({
                    ...prev,
                    updatingId: row.original.category_id,
                    form: {
                        name: row.original.category_name,
                        description: row.original.category_description,
                    },
                }));
                setModalOpen(true) 
                }}><Icon src={editIcon} alt="edit" /></IconButton>
                <IconButton onClick={() => deleteCategory(row.original.category_id)}>
                    <Icon src={trashIcon} alt="delete" />
                </IconButton>
                </ActionButtons>
            ) 
        },
    ]
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
                                onClick={() => {
                                    setState((prev) => ({
                                        ...prev,
                                        modalOpen: true,
                                        updatingId: null,
                                        form: { name: "", description: "" },
                                    }))
                                    setModalOpen(true)
                                 }}
                                >
                                    <Icon src={plusIcon} alt="add" /> 
                                    Add Category 
                                   </AddButton>
                            )
                        }
                </PageHeader>
                <Datatable 
                data={state.data}
                columns={columns}
                globalFilter={globalFilter}
                setGlobalFilter={setGlobalFilter}          
                /> 
                <ModalBox
                open={modalOpen}
                onClose={() => {
                    setModalOpen(false)
                    setState(prev => ({ ...prev, updatingId: null }))
                } }
                >
                    <ModalHeader>
                        <ModalTitle>
                            {state.updatingId ? "Edit Category" : "Add Category"}
                        </ModalTitle>
                        <CloseButton
                        onClick={() => {
                            setModalOpen(false)
                            setState(prev => ({ ...prev, updatingId: null }))
                            }}>
                            <Icon src={xIcon} alt="close" />
                        </CloseButton>
                    </ModalHeader>
                <Form onSubmit={(e) => {e.preventDefault(); handleSubmit()}}>
                    <FormGroup>
                        <Label>Name *</Label>
                        <Input value={state.form.name}  onChange={(e) => setState((prev) => ({
                            ...prev,
                            form: { ...prev.form, name: e.target.value }
                        }))} required/>
                    </FormGroup>
                    <FormGroup>
                        <Label>Description *</Label>
                        <Textarea value={state.form.description} onChange={(e) => setState((prev)=> ({
                            ...prev,
                            form: { ...prev.form, description: e.target.value }
                        }))}/>
                    </FormGroup>
                    <ModalActions>
                        <CancelButton
                        onClick={() => {
                            setModalOpen(false)
                            setState(prev => ({ ...prev, updatingId: null }))
                            }}>
                            Cancel
                        </CancelButton>
                        <SubmitButton>
                            {state.updatingId ? "Update" : "Create"}
                        </SubmitButton>
                    </ModalActions>
                </Form>
        </ModalBox>
                    
        <ToastContainer position="top-right" autoClose={3000} />        
    </>
   )}

export default Categories

