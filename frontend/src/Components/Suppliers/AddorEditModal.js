import {ModalHeader, ModalTitle, CloseButton, Form, Input, ModalActions, Icon, Button} from "../CommonUIStyles/Styles"
import xIcon from "../../Images/cross.svg"
import ModalBox from "../ModalBox"
import api from "../../Api/axios"
import { toast } from "react-toastify"
import { useEffect, useState } from "react"
import FormField from "../FormField"
import { validateSupplierForm } from "../../utils/validation"

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

  const createSupplier = async () => {
   if (!validateSupplierForm(supplier)) return

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
    if (!validateSupplierForm(supplier)) return

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

  <FormField label="Name" required>
    <Input
      data-cy="supplier-name"
      value={supplier.name}
      onChange={(e) =>
        setSupplier((p) => ({ ...p, name: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Contact" required>
    <Input
      data-cy="supplier-contact"
      value={supplier.contact}
      onChange={(e) =>
        setSupplier((p) => ({ ...p, contact: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Email" required>
    <Input
      data-cy="supplier-email"
      value={supplier.email}
      onChange={(e) =>
        setSupplier((p) => ({ ...p, email: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Company" required>
    <Input
      data-cy="supplier-company"
      value={supplier.company}
      onChange={(e) =>
        setSupplier((p) => ({ ...p, company: e.target.value }))
      }
    />
  </FormField>

  <ModalActions>
    <Button
      variant="secondary"
      data-cy="supplier-cancel"
      onClick={() => {
        setSupplier(initialSupplier);
        onClose();
      }}
    >
      Cancel
    </Button>

    <Button varaint="primary" data-cy="supplier-submit">
      {updatingId ? "Update" : "Create"}
    </Button>
  </ModalActions>

</Form>

    </ModalBox>
  )
}

export default AddOrEditSupplierModal
