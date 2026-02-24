import {
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  Input,
  ModalActions,
  Icon,
  Button
} from "../CommonUIStyles/Styles";
import { SupplierSelect } from "./Styles";
import xIcon from "../../Images/cross.svg";
import ModalBox from "../ModalBox";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState, useMemo } from "react";
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

  const supplierOptions = useMemo(
    () =>
      (suppliers || []).map((s) => ({
        value: s.supplier_id,
        label: s.name || s.company || `Supplier #${s.supplier_id}`,
      })),
    [suppliers]
  );

  const selectedSupplier = useMemo(
    () =>
      supplierOptions.find(
        (o) => o.value === Number(userObj.supplier_id)
      ) || null,
    [supplierOptions, userObj.supplier_id]
  );

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
    supplier_id: Number(userObj.supplier_id),
  }
  
  if (userObj.password && userObj.password.trim() !== "") {
    payload.password = userObj.password
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
        <ModalTitle data-cy="user-modal-title">{updatingId ? "Edit User" : "Add User"}</ModalTitle>
        <CloseButton onClick={() => {onClose(); setUserObj(initialUser)}}>
          <Icon src={xIcon} alt="close" />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={onSubmitHandler}>

  <FormField label="Name" required>
    <Input
      data-cy="user-name"
      value={userObj.name}
      onChange={(e) =>
        setUserObj((p) => ({ ...p, name: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Email" required>
    <Input
      data-cy="user-email"
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
        options={supplierOptions}
        value={selectedSupplier}
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
      data-cy="user-password"
      type="password"
      value={userObj.password}
      onChange={(e) =>
        setUserObj((p) => ({ ...p, password: e.target.value }))
      }
    />
  </FormField>

  <ModalActions>
    <Button
      data-cy="user-cancel-btn"
      variant="secondary"
      onClick={() => {
        onClose();
        setUserObj(initialUser);
      }}
    >
      Cancel
    </Button>

    <Button data-cy="user-submit-btn" variant="primary">
      {updatingId ? "Update" : "Create"}
    </Button>
  </ModalActions>

</Form>

    </ModalBox>
  );
};

export default AddOrEditUserModal
