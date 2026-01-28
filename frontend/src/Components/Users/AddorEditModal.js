import {
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  FormGroup,
  Label,
  Input,
  ModalActions,
  CancelButton,
  SubmitButton,
  Icon,
  SupplierSelect
} from "./Styles";
import xIcon from "../../Images/cross.svg";
import ModalBox from "../ModalBox";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState } from "react";

const initialUser = { name: "", email: "", password: "", supplier_id: "" }

const AddOrEditUserModal = ({ open, onClose, updatingId, form, refetch, suppliers , setLoading}) => {
  const [userObj, setUserObj] = useState(initialUser)

  useEffect(() => {
    if (!updatingId) {
      setUserObj(initialUser);
      return;
    }

    setUserObj(form)
  }, [updatingId, form])

  const validateUserForm = () => {
  const name = userObj.name.trim()
  const email = userObj.email.trim()
  const password = userObj.password
  const supplierId = userObj.supplier_id

  if (!name) {
    toast.error("User name is required")
    return false
  }

  if (name.length < 2) {
    toast.error("User name must be at least 2 characters")
    return false
  }

  if (!email) {
    toast.error("Email cannot be empty")
    return false
  }

  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(email)) {
    toast.error("Invalid email format (e.g., user@example.com)")
    return false
  }

  if (!updatingId) {
    if (!password) {
      toast.error("Password is required")
      return false
    }

    if (password.length < 8) {
      toast.error("Password must be at least 8 characters long")
      return false
    }

    if (!/^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*(),.?":{}|<>_\-+=~`])/.test(password)) {
          toast.error("Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
          return false
    }

  }

  if (form.role !== "system_admin") {
    if (!supplierId || Number(supplierId) <= 0) {
      toast.error("Please select a valid supplier")
      return false
    }
  }

  return true
}


  const createUser = async () => {

    if (!validateUserForm()) return

    const payload = {
    name: userObj.name.trim(),
    email: userObj.email.trim(),
    supplier_id: Number(userObj.supplier_id),
   }

   if (userObj.password && userObj.password.trim() !== "") {
    payload.password = userObj.password
  }

   setLoading(true)


    try {
      await api.post("users/supplier-admin", payload);

      toast.success("User created")
      setUserObj(initialUser)
      onClose()
      refetch()
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to create user")
    } finally{
      setLoading(false)
    }
  };

  const updateUser = async () => {

    if (!validateUserForm()) return

    const payload = {
    name: userObj.name.trim(),
    email: userObj.email.trim(),
    password: userObj.password,
    supplier_id: Number(userObj.supplier_id),
   }

   setLoading(true)
    try {
      await api.put(`users/supplier-admin/${updatingId}`, payload)

      toast.success("User updated");
      setUserObj(initialUser);
      onClose()
      refetch()
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to update user")
    } finally{
      setLoading(false)
    }
  };

  const onSubmitHandler = (e) => {
    e.preventDefault()
    updatingId ? updateUser() : createUser()
  };

  return (
    <ModalBox open={open} onClose={onClose}>
      <ModalHeader>
        <ModalTitle>{updatingId ? "Edit User" : "Add User"}</ModalTitle>
        <CloseButton onClick={() => {onClose(); setUserObj(initialUser)}}>
          <Icon src={xIcon} alt="close" />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={onSubmitHandler}>
        <FormGroup>
          <Label>Name *</Label>
          <Input
            value={userObj.name}
            onChange={(e) =>
              setUserObj((p) => ({ ...p, name: e.target.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Email *</Label>
          <Input
            value={userObj.email}
            onChange={(e) =>
              setUserObj((p) => ({ ...p, email: e.target.value }))
            }
          />
        </FormGroup>
       
       {form.role !== "system_admin" && (
        <FormGroup>
            <Label>Supplier *</Label>
            <SupplierSelect
                classNamePrefix="react-select"
                options={suppliers.map(s => ({
                value: s.supplier_id,
                label: s.name || s.company || `Supplier #${s.supplier_id}`
            }))}

                value={suppliers
                .map(s => ({
                value: s.supplier_id,
                label: s.name || s.company || `Supplier #${s.supplier_id}`
            }))
            .find(opt => opt.value === Number(userObj.supplier_id)) || null}

            onChange={(opt) =>
            setUserObj(prev => ({
            ...prev,
            supplier_id: opt.value
            }))
           }
            placeholder="Select supplier"
            isSearchable
           />
        </FormGroup>
       )}

        <FormGroup>
          <Label>
            {updatingId ? "Password (leave blank to keep same)" : "Password *"}
          </Label>
          <Input
            type="password"
            value={userObj.password}
            onChange={(e) =>
              setUserObj((p) => ({ ...p, password: e.target.value }))
            }
            {...(!updatingId && { required: true })}
          />
        </FormGroup>
        <ModalActions>
          <CancelButton onClick={() => {onClose(); setUserObj(initialUser)}}>Cancel</CancelButton>
          <SubmitButton>{updatingId ? "Update" : "Create"}</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditUserModal
