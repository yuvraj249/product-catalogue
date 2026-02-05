import { ModalHeader, ModalTitle, CloseButton, Form, Input, Textarea, ModalActions, Icon, Button } from './Styles'
import xIcon from '../../Images/cross.svg'
import ModalBox from '../ModalBox'
import api from '../../Api/axios'
import { toast } from 'react-toastify'
import { useState } from 'react'
import { useEffect } from 'react'
import FormField from '../FormField'
import { validateCategoryForm } from '../../utils/validation'

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


const createCategory = async () => {
    const valid = await validateCategoryForm(categoryObj, updatingId)
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
    const valid = await validateCategoryForm(categoryObj, updatingId)
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
                        <ModalTitle>
                            {updatingId ? "Edit Category" : "Add Category"}
                        </ModalTitle>
                        <CloseButton
                        onClick={() => {
                          onClose()
                          setCategoryObj(initialCategoryObject)
                        }}>
                            <Icon src={xIcon} alt="close" />
                        </CloseButton>
                </ModalHeader>
                <Form onSubmit={onSubmitHandlder}>

                  <FormField label="Name" required>
                    <Input
                      value={categoryObj.name}
                      onChange={onChangeName}
                   />
                  </FormField>

                  <FormField label="Description" required>
                    <Textarea
                    value={categoryObj.description}
                    onChange={onChangeDesc}
                  />
                  </FormField>

                  <ModalActions>
                    <Button
                      variant="secondary"
                      onClick={() => {
                      onClose();
                      setCategoryObj(initialCategoryObject);
                     }}
                    >Cancel
                    </Button>

                    <Button  variant="primary">
                      {updatingId ? "Update" : "Create"}
                    </Button>
                    </ModalActions>

                    </Form>
            </ModalBox>
  )
}

export default AddorEditModal