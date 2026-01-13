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

const AddOrEditProductModal = ({ open, onClose, updatingId, form, refetch, categories }) => {
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

  const createProduct = async () => {
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
  }
}


const updateProduct = async () => {
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
  }
}


  return (
    <ModalBox open={open} onClose={onClose}>
      <ModalHeader>
        <ModalTitle>{updatingId ? "Edit Product" : "Add Product"}</ModalTitle>
        <CloseButton onClick={onClose}>
          <Icon src={xIcon} />
        </CloseButton>
      </ModalHeader>

      <Form onSubmit={(e) => e.preventDefault()}>
        <FormGroup>
          <Label>Name *</Label>
          <Input
            value={product.product_name}
            onChange={(e) =>
              setProduct((p) => ({ ...p, product_name: e.target.value }))
            }
            required
          />
        </FormGroup>

        <FormGroup>
          <Label>Category *</Label>
          <CategorySelect
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
            type="number"
            value={product.product_cost}
            onChange={(e) =>
              setProduct((p) => ({ ...p, product_cost: e.target.value }))
            }
            required
          />
        </FormGroup>

        <FormGroup>
          <Label>Discount Type</Label>
          <Input
            value={product.discount_type}
            onChange={(e) =>
              setProduct((p) => ({ ...p, discount_type: e.target.value }))
            }
          />
        </FormGroup>

        <FormGroup>
          <Label>Discount Value</Label>
          <Input
            type="number"
            value={product.discount_value}
            onChange={(e) =>
              setProduct((p) => ({ ...p, discount_value: e.target.value }))
            }
          />
        </FormGroup>

        <ModalActions>
          <CancelButton onClick={onClose}>Cancel</CancelButton>
          <SubmitButton onClick={updatingId ? updateProduct : createProduct}>
            {updatingId ? "Update" : "Create"}
          </SubmitButton>
        </ModalActions>
      </Form>
    </ModalBox>
  );
};

export default AddOrEditProductModal;
