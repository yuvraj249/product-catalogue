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
import { CategorySelect } from "./Styles";
import ModalBox from "../ModalBox";
import xIcon from "../../Images/cross.svg";
import api from "../../Api/axios";
import { toast } from "react-toastify";
import { useEffect, useState, useMemo } from "react";
import FormField from "../FormField";
import { validateProductForm } from "../../utils/validation";

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

  const categoryOptions = useMemo(
    () =>
      (categories || []).map((c) => ({
        value: c.category_id,
        label: c.category_name,
      })),
    [categories]
  );

  const selectedCategory = useMemo(
    () =>
      categoryOptions.find(
        (o) => o.value === Number(product.product_category_id)
      ) || null,
    [categoryOptions, product.product_category_id]
  );


  const createProduct = async () => {
  if(!validateProductForm(product)) return
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
  if(!validateProductForm(product)) return
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

      <Form onSubmit={(e) => e.preventDefault()}>

  <FormField label="Name" required>
    <Input
      value={product.product_name}
      onChange={(e) =>
        setProduct((p) => ({ ...p, product_name: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Category" required>
    <CategorySelect
      classNamePrefix="react-select"
      options={categoryOptions}
      value={selectedCategory}
    />
  </FormField>

  <FormField label="Cost" required>
    <Input
      type="number"
      value={product.product_cost}
      onChange={(e) =>
        setProduct((p) => ({ ...p, product_cost: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Discount Type">
    <Input
      value={product.discount_type}
      onChange={(e) =>
        setProduct((p) => ({ ...p, discount_type: e.target.value }))
      }
    />
  </FormField>

  <FormField label="Discount Value">
    <Input
      type="number"
      value={product.discount_value}
      onChange={(e) =>
        setProduct((p) => ({ ...p, discount_value: e.target.value }))
      }
    />
  </FormField>

  <ModalActions>
    <Button
      variant="secondary"
      onClick={() => {
        onClose();
        setProduct(initialProduct);
      }}
    >
      Cancel
    </Button>

    <Button varaint="primary" onClick={updatingId ? updateProduct : createProduct}>
      {updatingId ? "Update" : "Create"}
    </Button>
  </ModalActions>

</Form>

    </ModalBox>
  );
};

export default AddOrEditProductModal;
