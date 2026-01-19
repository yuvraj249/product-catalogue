import { ModalHeader, ModalTitle, CloseButton, Form, FormGroup, Label, Input, Textarea, ModalActions, CancelButton, SubmitButton, Icon } from './Styles'
import xIcon from '../../Images/cross.svg'
import ModalBox from '../ModalBox'
import api from '../../Api/axios'
import { toast } from 'react-toastify'
import { useState } from 'react'
import { useEffect } from 'react'

const initialCategoryObject = {
  name: "",
  description: ""
}

const AddorEditModal = ({open, onClose, updatingId, refetch, setLoading}) => {

const [categoryObj, setCategoryObj] = useState(initialCategoryObject)


useEffect(() => {
  if (!updatingId) {
    setCategoryObj(initialCategoryObject);
    return;
  }

  const fetchCategory = async () => {
    try {
      const res = await api.get(`/categories/${updatingId}`);
      const cat = res.data.category

      setCategoryObj({
          name: cat.category_name,
          description: cat.category_description,
      });
    } catch (err) {
      toast.error("Failed to load category");
    }
  };

  fetchCategory();
}, [updatingId])

const validateForm = async () => {
  const name = categoryObj.name.trim()
  const desc = categoryObj.description.trim()

  if (!name) {
    toast.error("Category name is required")
    return false
  }

  if (name.length < 2) {
    toast.error("Category name must be at least 2 characters")
    return false
  }

  if (name.length > 50) {
    toast.error("Category name must be under 50 characters")
    return false
  }

  if (desc.length > 200) {
    toast.error("Description must be under 200 characters")
    return false
  }

  try {
    const res = await api.get("/categories", {
      params: { q: name }
    })

    const exists = res.data.categories?.some(
      c => c.category_name.toLowerCase() === name.toLowerCase()
    )

    if (exists && !updatingId) {
      toast.error("Category with this name already exists")
      return false
    }
  } catch {
    toast.error("Failed to validate category name")
    return false
  }

  return true
}



const createCategory = async () => {
    const valid = await validateForm()
    if (!valid) return

    setLoading(true)
    const payload = {
      category_name: categoryObj.name,
      category_description: categoryObj.description
    }

    try{
        await api.post('/categories', payload)
        toast.success("Category created successfully")
        setCategoryObj(initialCategoryObject)
        onClose()
        refetch()
    }catch (err){
        const msg = err.response.data.error || ""
        toast.error(msg || "Failed to create or edit category")

    }finally{
      setLoading(false)
    }
   }

const updateCategory = async () => {
    const valid = await validateForm()
    if (!valid) return

    setLoading(true)
    const payload = {
      category_name: categoryObj.name,
      category_description: categoryObj.description
    }
    try {
        await api.put(`/categories/${updatingId}`, payload)
        toast.success("Category updated successfully")
        setCategoryObj(initialCategoryObject)
        onClose()
        refetch()
    } catch (err) {
        const msg = err.response.data.error || "";
        toast.error(msg || "Failed to update category");
    }finally{
      setLoading(false)
    }

   }

const handleSubmit = () => {
    updatingId ? updateCategory() : createCategory();
  }
  
const onSubmitHandlder = (e) => {
    e.preventDefault()
    handleSubmit()
  }

  
const onChangeName = (e) => {
     setCategoryObj(prev => ({
      ...prev,
      name: e.target.value
    }))
  }
  
  
const onChangeDesc = (e) => {
    setCategoryObj(prev => ({
      ...prev,
      description: e.target.value
    }))
  }

  return (
            <ModalBox
                open={open}
                onClose={onClose}
                >
                <ModalHeader>
                        <ModalTitle data-cy="category-modal-title">
                            {updatingId ? "Edit Category" : "Add Category"}
                        </ModalTitle>
                        <CloseButton
                        data-cy="close-category-modal"
                        onClick={() => {
                          onClose()
                          setCategoryObj(initialCategoryObject)
                        }}>
                            <Icon src={xIcon} alt="close" />
                        </CloseButton>
                </ModalHeader>
                <Form onSubmit={onSubmitHandlder}>
                    <FormGroup>
                        <Label>Name *</Label>
                        <Input  data-cy="category-name" value={categoryObj.name}  onChange={onChangeName}/>
                    </FormGroup>
                    <FormGroup>
                        <Label>Description *</Label>
                        <Textarea  data-cy="category-description" value={categoryObj.description} onChange={onChangeDesc}/>
                    </FormGroup>
                    <ModalActions>
                        <CancelButton data-cy="category-cancel" onClick={() => {onClose(); setCategoryObj(initialCategoryObject)}}>Cancel</CancelButton>
                        <SubmitButton data-cy="category-submit">
                            {updatingId ? "Update" : "Create"}
                        </SubmitButton>
                    </ModalActions>
                </Form>
            </ModalBox>
  )
}

export default AddorEditModal