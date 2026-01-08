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

const AddOrEditUserModal = ({ open, onClose, updatingId, form, refetch, suppliers }) => {
  const [userObj, setUserObj] = useState(initialUser)

  useEffect(() => {
    if (!updatingId) {
      setUserObj(initialUser);
      return;
    }

    setUserObj(form)
  }, [updatingId, form])

  const createUser = async () => {
    try {
      await api.post("users/supplier-admin", {
        name: userObj.name,
        email: userObj.email,
        password: userObj.password,
        supplier_id: Number(userObj.supplier_id),
      });

      toast.success("User created")
      setUserObj(initialUser)
      onClose()
      refetch()
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to create user")
    }
  };

  const updateUser = async () => {
    try {
      await api.put(`users/supplier-admin/${updatingId}`, {
        name: userObj.name,
        email: userObj.email,
        password: userObj.password || undefined,
        supplier_id: Number(userObj.supplier_id),
      })

      toast.success("User updated");
      setUserObj(initialUser);
      onClose()
      refetch()
    } catch (err) {
      toast.error(err.response?.data?.error || "Failed to update user")
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
        <CloseButton onClick={onClose}>
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
            required
          />
        </FormGroup>

        <FormGroup>
          <Label>Email *</Label>
          <Input
            value={userObj.email}
            onChange={(e) =>
              setUserObj((p) => ({ ...p, email: e.target.value }))
            }
            required
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
          <CancelButton onClick={onClose}>Cancel</CancelButton>
          <SubmitButton>{updatingId ? "Update" : "Create"}</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditUserModal
