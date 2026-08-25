import {
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  Input,
  Button,
  ModalActions,
  Icon,
} from "../CommonUIStyles/Styles";
import { SelectBox } from "./Styles";
import ModalBox from "../ModalBox";
import xIcon from "../../Images/cross.svg";
import api from "../../api/axios";
import { toast } from "react-toastify";
import { useEffect, useState, useMemo } from "react";
import FormField from "../FormField";
import { validateStockForm } from "../../utils/validation";

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


const productOptions = useMemo(() => {
    return (products || []).map((p) => ({
      value: p.product_id,
      label: p.product_name,
    }));
  }, [products]);

  const selectedProduct = useMemo(() => {
    return (
      productOptions.find((o) => o.value === form.product_id) || null
    );
  }, [productOptions, form.product_id]);

  const movementTypeOptions = useMemo(
    () => [
      { value: "IN", label: "IN (Add Stock)" },
      { value: "OUT", label: "OUT (Remove Stock)" },
    ],
    []
  );

  const selectedMovementType = useMemo(() => {
    return (
      movementTypeOptions.find(
        (o) => o.value === form.movement_type
      ) || null
    );
  }, [movementTypeOptions, form.movement_type])


  const addStockMovement = async () => {

  if (!validateStockForm(form)) return
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

  if (!validateStockForm(form)) return
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
        <ModalTitle data-cy="stock-modal-title">Update Stock</ModalTitle>
        <CloseButton data-cy="stock-modal-close" onClick={onClose}>
          <Icon src={xIcon} />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={(e) => e.preventDefault()}>

  <FormField label="Product" required>
    <SelectBox
      data-cy="stock-product"
      classNamePrefix="react-select"
      options={productOptions}
      value={selectedProduct}
      onChange={(opt) =>
        setForm((f) => ({ ...f, product_id: opt.value }))
      }
    />
  </FormField>

  <FormField label="Quantity" required>
    <Input
      data-cy="stock-quantity"
      type="number"
      value={form.quantity}
      onChange={(e) =>
        setForm((f) => ({ ...f, quantity: e.target.value }))
      }
      min="1"
    />
  </FormField>

  <FormField label="Movement Type" required>
    <SelectBox
      data-cy="stock-movement-type"
      classNamePrefix="react-select"
      options={movementTypeOptions}
      value={selectedMovementType}
      onChange={(opt) =>
        setForm((f) => ({ ...f, movement_type: opt.value }))
      }
    />
  </FormField>

  <FormField label="Reason">
    <Input
      data-cy="stock-reason"
      value={form.reason}
      onChange={(e) =>
        setForm((f) => ({ ...f, reason: e.target.value }))
      }
      placeholder="optional"
    />
  </FormField>

  <ModalActions>
    <Button
      variant="secondary"
      onClick={() => {
        onClose();
        setForm(initialState);
      }}
    >
      Cancel
    </Button>

    <Button data-cy="stock-submit-btn" variant="primary" onClick={submit}>Save</Button>
  </ModalActions>

</Form>

    </ModalBox>
  );
};

export default AddOrEditStockModal
