import {ModalHeader, ModalTitle, CloseButton, Form, FormGroup, Label, Input, ModalActions, CancelButton, SubmitButton, Icon} from "./Styles"
import xIcon from "../../Images/cross.svg"
import ModalBox from "../ModalBox"
import api from "../../Api/axios"
import { toast } from "react-toastify"
import { useEffect, useState } from "react"

const initialSupplier = {
  name: "",
  contact: "",
  email: "",
  company: "",
};

const AddOrEditSupplierModal = ({ open, onClose, updatingId, refetch, setLoading }) => {
  const [supplier, setSupplier] = useState(initialSupplier)
  useEffect(() => {
    if (!updatingId) {
      setSupplier(initialSupplier)
      return
    }
    const fetchSupplier = async () => {
      try {
        const res = await api.get(`/suppliers/${updatingId}`)

        const s = res.data.supplier;

        setSupplier({
          name: s.name,
          contact: s.contact_info,
          email: s.email,
          company: s.company,
        });
      } catch {
        toast.error("Failed to load supplier")
      }
    };
    fetchSupplier()
  }, [updatingId])

  const validateSupplierForm = () => {
  const name = supplier.name.trim()
  const contact = supplier.contact.trim()
  const email = supplier.email.trim()
  const company = supplier.company.trim()

  if (!name) {
    toast.error("Supplier name required")
    return false
  }

  if (name.length > 50) {
    toast.error("Supplier name too long")
    return false
  }

  if (name.length < 2) {
    toast.error("Supplier name too short")
    return false
  }

  const validName = /^[A-Za-z ]+$/
  if (!validName.test(name)) {
    toast.error("Supplier name should only contain alphabets and spaces")
    return false
  }

  if (!/[A-Za-z]/.test(name)) {
    toast.error("Supplier name must contain at least one letter")
    return false
  }

  if (!contact) {
    toast.error("Contact number required")
    return false
  }

  const validPhone = /^[0-9+\-() ]{7,15}$/
  if (!validPhone.test(contact)) {
    toast.error("Contact should contain only numbers, +, -, () or spaces")
    return false
  }

  if (!email) {
    toast.error("Email is required")
    return false
  }

  const validEmail =
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!validEmail.test(email)) {
    toast.error("Invalid email format")
    return false
  }

  if (!company) {
    toast.error("Company name required")
    return false
  }

  if (company.length > 50) {
    toast.error("Company name too long")
    return false
  }

  if (company.length < 2) {
    toast.error("Company name too short")
    return false
  }

  const validCompany = /^[A-Za-z0-9 ]+$/
  if (!validCompany.test(company)) {
    toast.error("Company name should contain only alphabets, numbers and spaces")
    return false
  }

  if (!/[A-Za-z]/.test(company)) {
    toast.error("Company name must contain at least one letter")
    return false
  }

  return true
}


  const createSupplier = async () => {
   if (!validateSupplierForm()) return

    setLoading(true)
    
    const payload = {
      name: supplier.name,
      contact_info: supplier.contact,
      email: supplier.email,
      company: supplier.company,
    };

    try {
      await api.post("/suppliers", payload);
      toast.success("Supplier created successfully");
      setSupplier(initialSupplier);
      onClose();
      refetch();
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to create supplier");
    }finally{
      setLoading(false)
    }
  }


  const updateSupplier = async () => {
    if (!validateSupplierForm()) return

    setLoading(true)
    const payload = {
      name: supplier.name,
      contact_info: supplier.contact,
      email: supplier.email,
      company: supplier.company,
    }

    try {
      await api.put(`/suppliers/${updatingId}`, payload)
      toast.success("Supplier updated successfully")
      setSupplier(initialSupplier)
      onClose()
      refetch()
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to update supplier")
    }finally{
      setLoading(false)
    }
  }

  const handleSubmit = () => {
    updatingId ? updateSupplier() : createSupplier()
  };

  const onSubmitHandler = (e) => {
    e.preventDefault();
    handleSubmit();
  };

  return (
    <ModalBox open={open} onClose={onClose}>
      <ModalHeader>
        <ModalTitle data-cy="supplier-modal-title">{updatingId ? "Edit Supplier" : "Add Supplier"}</ModalTitle>
        <CloseButton  data-cy="supplier-close" onClick={() => {onClose(); setSupplier(initialSupplier)}}>
          <Icon src={xIcon} alt="close" />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={onSubmitHandler}>
        <FormGroup>
          <Label>Name *</Label>
          <Input
            data-cy="supplier-name"
            value={supplier.name}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, name: e.target.value }))
            }
          />
        </FormGroup>


        <FormGroup>
          <Label>Contact *</Label>
          <Input
            data-cy="supplier-contact"
            value={supplier.contact}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, contact: e.target.value }))
            }
      
          />
        </FormGroup>


        <FormGroup>
          <Label>Email *</Label>
          <Input
            data-cy="supplier-email"
            value={supplier.email}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, email: e.target.value }))
            }
  
          />
        </FormGroup>

        <FormGroup>
          <Label>Company *</Label>
          <Input
            data-cy="supplier-company"
            value={supplier.company}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, company: e.target.value }))
            }
          />
        </FormGroup>
        <ModalActions>
          <CancelButton data-cy="supplier-cancel" onClick={() => {setSupplier(initialSupplier); onClose() }}>Cancel</CancelButton>

          <SubmitButton data-cy="supplier-submit">{updatingId ? "Update" : "Create"}</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  )
}

export default AddOrEditSupplierModal
