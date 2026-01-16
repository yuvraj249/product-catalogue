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
  SelectBox,
} from "./Styles";
import ModalBox from "../ModalBox";
import xIcon from "../../Images/cross.svg";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState } from "react";

const initialState = {
  product_id: null,
  quantity: "",
  movement_type: "IN",
  reason: "",
};

const AddOrEditStockModal = ({ open, onClose, refetch, products, editingItem, setLoading}) => {
  const [form, setForm] = useState(initialState)

  useEffect(() => {
  if (open && editingItem) {
    setForm({
      product_id: editingItem.product_id,
      quantity: editingItem.quantity,
      movement_type: editingItem.movement_type,
      reason: editingItem.reason || "",
    })
  }

  if (!open) {
    setForm(initialState)
  }
}, [open, editingItem])


const validateStockForm = () => {
    const qty = Number(form.quantity)

    if (!form.product_id || form.product_id <= 0) {
      toast.error("Please select a valid product")
      return false
    }

    if (!form.quantity) {
      toast.error("Quantity is required")
      return false
    }

    if (isNaN(qty) || qty <= 0) {
      toast.error("Quantity must be a positive number")
      return false
    }

    if (!Number.isInteger(qty)) {
    toast.error("Quantity must be a whole number (no decimals)")
    return false
  }

  if (qty > 100000) {
      toast.error("Quantity too large")
      return false
    }

    if (!form.movement_type || !["IN", "OUT"].includes(form.movement_type)) {
      toast.error("Invalid movement type")
      return false
    }

    if (form.reason && form.reason.length > 200) {
      toast.error("Reason too long (max 200 characters)")
      return false
    }

    return true
  }


  const addStockMovement = async () => {

  if (!validateStockForm()) return
  setLoading(true)
  if (!form.product_id || !form.quantity) {
    toast.error("product and quantity required")
    return
  }

  const payload = {
    product_id: Number(form.product_id),
    quantity: Number(form.quantity),
    movement_type: form.movement_type,
    reason: form.reason?.trim(),
  }

  try {
    await api.post("/stock_movements", payload)
    toast.success("stock added")
    onClose()
    refetch()
  } catch (e) {
    toast.error(e.response?.data?.error || "add failed")
  } finally{
    setLoading(false)
  }
} 


const updateStockMovement = async () => {

  if (!validateStockForm()) return
  setLoading(true)
  if (!form.product_id || !form.quantity) {
    toast.error("product and quantity required")
    return
  }

  const payload = {
    product_id: Number(form.product_id),
    quantity: Number(form.quantity),
    movement_type: form.movement_type,
    reason: form.reason?.trim(),
  }

  try {
    await api.put(`/stock_movements/${editingItem.stock_id}`, payload)
    toast.success("stock updated")
    onClose()
    refetch()
  } catch (e) {
    toast.error(e.response?.data?.error || "update failed")
  } finally{
    setLoading(false)
  }
}

const submit = () => {
  if (editingItem) {
    updateStockMovement()
  } else {
    addStockMovement()
  }
}

  return (
    <ModalBox open={open} onClose={() => {onClose(); setForm(initialState)}}>
      <ModalHeader>
        <ModalTitle>Update Stock</ModalTitle>
        <CloseButton onClick={onClose}>
          <Icon src={xIcon} />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={(e) => e.preventDefault()}>
        <FormGroup>
          <Label>Product *</Label>
          <SelectBox
            classNamePrefix='react-select'
            options={(products || []).map(p => ({
              value: p.product_id,
              label: p.product_name,
            }))}
            value={(products || [])
              .map(p => ({ value: p.product_id, label: p.product_name }))
              .find(o => o.value === form.product_id) || null}
            onChange={(opt) =>
              setForm(f => ({ ...f, product_id: opt.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Quantity *</Label>
          <Input
            type="number"
            value={form.quantity}
            onChange={(e) =>
              setForm(f => ({ ...f, quantity: e.target.value }))
            }
            min="1"
          />
        </FormGroup>

        <FormGroup>
            <Label>Movement Type *</Label>
            <SelectBox
                classNamePrefix='react-select'
                options={[
                    { value: "IN", label: "IN (Add Stock)" },
                    { value: "OUT", label: "OUT (Remove Stock)" },
                ]}
                value={[
                    { value: "IN", label: "IN (Add Stock)" },
                    { value: "OUT", label: "OUT (Remove Stock)" },
                    ].find(o => o.value === form.movement_type)}
                onChange={(opt) =>
                    setForm(f => ({ ...f, movement_type: opt.value }))
                }
            />
        </FormGroup>

        <FormGroup>
          <Label>Reason</Label>
          <Input
            value={form.reason}
            onChange={(e) =>
              setForm(f => ({ ...f, reason: e.target.value }))
            }
            placeholder="optional"
          />
        </FormGroup>

        <ModalActions>
          <CancelButton onClick={() => {onClose(); setForm(initialState)}}>Cancel</CancelButton>
          <SubmitButton onClick={submit}>Save</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditStockModal
