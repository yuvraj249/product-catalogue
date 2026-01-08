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

const AddOrEditSupplierModal = ({ open, onClose, updatingId, refetch }) => {
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
    }
  }


  const updateSupplier = async () => {
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
        <ModalTitle>{updatingId ? "Edit Supplier" : "Add Supplier"}</ModalTitle>
        <CloseButton onClick={onClose}>
          <Icon src={xIcon} alt="close" />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={onSubmitHandler}>
        <FormGroup>
          <Label>Name *</Label>
          <Input
            value={supplier.name}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, name: e.target.value }))
            }
            required
          />
        </FormGroup>


        <FormGroup>
          <Label>Contact *</Label>
          <Input
            value={supplier.contact}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, contact: e.target.value }))
            }
            required
          />
        </FormGroup>


        <FormGroup>
          <Label>Email *</Label>
          <Input
            value={supplier.email}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, email: e.target.value }))
            }
            required
          />
        </FormGroup>

        <FormGroup>
          <Label>Company *</Label>
          <Input
            value={supplier.company}
            onChange={(e) =>
              setSupplier((p) => ({ ...p, company: e.target.value }))
            }
            required
          />
        </FormGroup>
        <ModalActions>
          <CancelButton onClick={onClose}>Cancel</CancelButton>

          <SubmitButton>{updatingId ? "Update" : "Create"}</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  )
}

export default AddOrEditSupplierModal
