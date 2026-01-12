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

const AddOrEditStockModal = ({ open, onClose, refetch, products, editingItem }) => {
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



  const addStockMovement = async () => {
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
  }
} 


const updateStockMovement = async () => {
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
    <ModalBox open={open} onClose={onClose}>
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
            required
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
          <CancelButton onClick={onClose}>Cancel</CancelButton>
          <SubmitButton onClick={submit}>Save</SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditStockModal
