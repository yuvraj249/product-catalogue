import {
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  Input,
  ModalActions,
  Icon,
  SupplierSelect,
  Button
} from "./Styles";
import xIcon from "../../Images/cross.svg";
import ModalBox from "../ModalBox";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState } from "react";
import FormField from "../FormField";
import { validateUserForm } from "../../utils/validation";

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

  const createUser = async () => {

    if (!validateUserForm(userObj, updatingId, form)) return

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

    if (!validateUserForm(userObj, updatingId, form)) return

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

  <FormField label="Name" required>
    <Input
      value={userObj.name}
      onChange={(e) =>
        setUserObj((p) => ({ ...p, name: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Email" required>
    <Input
      value={userObj.email}
      onChange={(e) =>
        setUserObj((p) => ({ ...p, email: e.target.value }))
      }
    />
  </FormField>

  {form.role !== "system_admin" && (
    <FormField label="Supplier" required>
      <SupplierSelect
        classNamePrefix="react-select"
        options={suppliers.map((s) => ({
          value: s.supplier_id,
          label: s.name || s.company || `Supplier #${s.supplier_id}`,
        }))}
        value={
          suppliers
            .map((s) => ({
              value: s.supplier_id,
              label: s.name || s.company || `Supplier #${s.supplier_id}`,
            }))
            .find(
              (opt) => opt.value === Number(userObj.supplier_id)
            ) || null
        }
        onChange={(opt) =>
          setUserObj((prev) => ({
            ...prev,
            supplier_id: opt.value,
          }))
        }
        placeholder="Select supplier"
        isSearchable
      />
    </FormField>
  )}

  <FormField
    label={
      updatingId
        ? "Password (leave blank to keep same)"
        : "Password"
    }
    required={!updatingId}
  >
    <Input
      type="password"
      value={userObj.password}
      onChange={(e) =>
        setUserObj((p) => ({ ...p, password: e.target.value }))
      }
      {...(!updatingId && { required: true })}
    />
  </FormField>

  <ModalActions>
    <Button
      variant="secondary"
      onClick={() => {
        onClose();
        setUserObj(initialUser);
      }}
    >
      Cancel
    </Button>

    <Button variant="primary">
      {updatingId ? "Update" : "Create"}
    </Button>
  </ModalActions>

</Form>

    </ModalBox>
  );
};

export default AddOrEditUserModal
