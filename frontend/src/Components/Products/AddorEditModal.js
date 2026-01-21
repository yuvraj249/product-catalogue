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
  CategorySelect,
} from "./Styles";
import ModalBox from "../ModalBox";
import xIcon from "../../Images/cross.svg";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState } from "react";

const initialProduct = {
  product_name: "",
  product_description: "",
  product_cost: "",
  product_category_id: null,
  discount_type: "",
  discount_value: "",
};

const AddOrEditProductModal = ({ open, onClose, updatingId, form, refetch, categories, setLoading }) => {
  const [product, setProduct] = useState(initialProduct);

  useEffect(() => {
    if (!updatingId) {
      setProduct(initialProduct);
      return;
    }

    setProduct({
      product_name: form.product_name ?? "",
      product_description: form.product_description ?? "",
      product_cost: String(form.product_cost ?? ""),
      product_category_id: Number(form.product_category_id) || null,
      discount_type: form.discount_type ?? "",
      discount_value: String(form.discount_value ?? ""),
    });
  }, [updatingId, form]);

  const validateProductForm = () => {
    const name = product.product_name.trim();
    const desc = product.product_description?.trim() || "";
    const cost = Number(product.product_cost);
    const categoryId = Number(product.product_category_id);
    const discountType = product.discount_type.trim().toLowerCase();
    const discountValue = Number(product.discount_value);

    if (!name) {
      toast.error("Product name required");
      return false;
    }

    if (name.length > 50) {
      toast.error("Product name too long");
      return false;
    }

    if (name.length < 2) {
      toast.error("Product name too short");
      return false;
    }

    const validName = /^[A-Za-z0-9 ]+$/;
    if (!validName.test(name)) {
      toast.error("Product name should contain only alphabets, numbers and spaces");
      return false;
    }

    if (!/[A-Za-z]/.test(name)) {
      toast.error("Product name must contain at least one letter");
      return false;
    }

    if (cost <= 0 || cost > 9999999999.5) {
      toast.error("Please enter a realistic product cost");
      return false;
    }

    if (desc && desc.length > 1500) {
      toast.error("Description too long");
      return false;
    }

    if (!categoryId || categoryId <= 0) {
      toast.error("Please select a valid category");
      return false;
    }

    if (discountValue !== 0 && !discountType) {
      toast.error("Discount type required if discount value is provided");
      return false;
    }

    if (discountType) {
      if (discountType !== "flat" && discountType !== "percent") {
        toast.error("Discount type must be flat or percent");
        return false;
      }

      if (discountValue === 0) {
        toast.error("Discount value required when discount type is selected");
        return false;
      }

      if (discountValue < 0) {
        toast.error("Discount value must be >= 0");
        return false;
      }

      if (discountType === "percent" && discountValue > 100) {
        toast.error("Percent discount cannot exceed 100");
        return false;
      }

      if (discountType === "flat" && discountValue > cost) {
        toast.error("Flat discount cannot exceed product cost");
        return false;
      }
    }

    return true;
  }

  const createProduct = async () => {
  if(!validateProductForm()) return
  setLoading(true)
  const payload = {
    product_name: product.product_name.trim(),
    product_description: product.product_description?.trim(),
    product_cost: Number(product.product_cost),
    product_category_id: Number(product.product_category_id),
    discount_type: product.discount_type?.trim(),
    discount_value: Number(product.discount_value),
  }

  try {
    await api.post("/products", payload);
    toast.success("Product created");
    onClose();
    refetch();
  } catch (e) {
    toast.error(e.response?.data?.error || "Failed");
  } finally{
    setLoading(false)
  }
}


const updateProduct = async () => {
  if(!validateProductForm()) return
  setLoading(true)
  const payload = {
    product_name: product.product_name.trim(),
    product_description: product.product_description?.trim(),
    product_cost: Number(product.product_cost),
    product_category_id: Number(product.product_category_id),
    discount_type: product.discount_type?.trim(),
    discount_value: Number(product.discount_value),
  }

  try {
    await api.put(`/products/${updatingId}`, payload);
    toast.success("Product updated");
    onClose();
    refetch();
  } catch (e) {
    toast.error(e.response?.data?.error || "Failed");
  } finally{
    setLoading(false)
  }
}


  return (
    <ModalBox open={open} onClose={onClose}>
      <ModalHeader>
        <ModalTitle data-cy="product-modal-title">{updatingId ? "Edit Product" : "Add Product"}</ModalTitle>
        <CloseButton onClick={() => {onClose(); setProduct(initialProduct)}}>
          <Icon src={xIcon} />
        </CloseButton>
      </ModalHeader>

      <Form data-cy="product-form" onSubmit={(e) => e.preventDefault()}>
        <FormGroup>
          <Label>Name *</Label>
          <Input
            data-cy="product-name"
            value={product.product_name}
            onChange={(e) =>
              setProduct((p) => ({ ...p, product_name: e.target.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Category *</Label>
          <CategorySelect
            data-cy="product-category"
            classNamePrefix="react-select"
            options={(categories || []).map((c) => ({
                value: c.category_id,
                label: c.category_name,
            }))}
            value={
              categories
                .map((c) => ({
                  value: c.category_id,
                  label: c.category_name,
                }))
                .find(
                  (o) => o.value === Number(product.product_category_id)
                ) || null
            }
            onChange={(opt) =>
              setProduct((p) => ({
                ...p,
                product_category_id: Number(opt.value),
              }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Cost *</Label>
          <Input
            data-cy="product-cost"
            type="number"
            value={product.product_cost}
            onChange={(e) =>
              setProduct((p) => ({ ...p, product_cost: e.target.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Discount Type</Label>
          <Input
            data-cy="product-discount-type"
            value={product.discount_type}
            onChange={(e) =>
              setProduct((p) => ({ ...p, discount_type: e.target.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Discount Value</Label>
          <Input
            data-cy="product-discount-value"
            type="number"
            value={product.discount_value}
            onChange={(e) =>
              setProduct((p) => ({ ...p, discount_value: e.target.value }))
            }
          />
        </FormGroup>

        <ModalActions>
          <CancelButton data-cy="product-cancel-btn" onClick={() => {onClose(); setProduct(initialProduct)}}>Cancel</CancelButton>
          <SubmitButton data-cy="product-submit-btn" onClick={updatingId ? updateProduct : createProduct}>
            {updatingId ? "Update" : "Create"}
          </SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditProductModal;
