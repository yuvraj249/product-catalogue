import { ModalHeader, ModalTitle, CloseButton, Form, FormGroup, Label, Input, Textarea, ModalActions, CancelButton, SubmitButton, Icon } from './Styles'
import xIcon from '../../Images/cross.svg'
import ModalBox from '../ModalBox'
import api from '../../Api/axios'
import { toast } from 'react-toastify'
import { useState } from 'react'
import { useEffect } from 'react'
const AddorEditModal = ({open, onClose, updatingId, refetch}) => {

const [object , setObject] = useState({
        fields: {name: "", description: ""}
})

useEffect(() => {
  if (!updatingId) {
    setObject({ fields: { name: "", description: "" } });
    return;
  }

  const fetchCategory = async () => {
    try {
      const res = await api.get(`/categories/${updatingId}`);
      const cat = res.data.category

      setObject({
        fields: {
          name: cat.category_name,
          description: cat.category_description,
        },
      });
    } catch (err) {
      toast.error("Failed to load category");
    }
  };

  fetchCategory();
}, [updatingId])



const createCategory = async () => {
    try{
        await api.post('/categories',{
            category_name: object.fields.name,
            category_description: object.fields.description
        })
        toast.success("Category created successfully")
        onClose()
        setObject({
            fields: {name: "", description: ""}
        })
        refetch()
    }catch (err){
        const msg = err.response.data.error || ""
        toast.error(msg || "Failed to create or edit category")

    }
   }

const updateCategory = async () => {
    try {
        await api.put(`/categories/${updatingId}`, {
        category_name: object.fields.name,
        category_description: object.fields.description
        })
        toast.success("Category updated successfully")
        onClose()
        setObject({
        fields: {name: "", description: ""}
        })
        refetch()
    } catch (err) {
        const msg = err.response.data.error || "";
        toast.error(msg || "Failed to update category");
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
    setObject((prev) => ({
        ...prev,
        fields: { ...prev.fields, name: e.target.value }
    }))
  }
  
  
const onChangeDesc = (e) => {
    setObject((prev)=> ({
        ...prev,
        fields: { ...prev.fields, description: e.target.value }
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
                        onClick={onClose}>
                            <Icon src={xIcon} alt="close" />
                        </CloseButton>
                </ModalHeader>
                <Form onSubmit={onSubmitHandlder}>
                    <FormGroup>
                        <Label>Name *</Label>
                        <Input value={object.fields.name}  onChange={onChangeName} required/>
                    </FormGroup>
                    <FormGroup>
                        <Label>Description *</Label>
                        <Textarea value={object.fields.description} onChange={onChangeDesc}/>
                    </FormGroup>
                    <ModalActions>
                        <CancelButton onClick={onClose}>Cancel</CancelButton>
                        <SubmitButton>
                            {updatingId ? "Update" : "Create"}
                        </SubmitButton>
                    </ModalActions>
                </Form>
            </ModalBox>
  )
}

export default AddorEditModal