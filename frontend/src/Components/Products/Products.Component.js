import React, { useState, useEffect, useMemo, useCallback } from "react";
import api from "../../Api/axios";
import { getUserInfo } from "../../utils/auth";

import {
  PageHeader,
  Title,
  AddButton,
  SearchWrap,
  SearchIpt,
  TableWrapper,
  Table,
  Th,
  Td,
  ActionButtons,
  IconButton,
  Modal,
  ModalContent,
  ModalHeader,
  ModalTitle,
  CloseButton,
  Form,
  FormGroup,
  Label,
  Input,
  Select,
  Textarea,
  ModalActions,
  CancelButton,
  SubmitButton,
  Icon,
  ToastContainer,
  Toast,
} from "./Products.styles";

import searchIcon from "../Products/assets2/search.svg";
import plusIcon from "../Products/assets2/plus.svg";
import editIcon from "../Products/assets2/edit.svg";
import trashIcon from "../Products/assets2/trash.svg";
import xIcon from "../Products/assets2/cross.svg";

const Products = () => {
  const user = getUserInfo();

  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  const [globalFilter, setGlobalFilter] = useState("");

  const [showModal, setShowModal] = useState(false);
  const [updatingId, setUpdatingId] = useState(null);

  const isUpdating = updatingId !== null;

  const [name, setName] = useState("");
  const [desc, setDesc] = useState("");
  const [cost, setCost] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [discountType, setDiscountType] = useState("");
  const [discountValue, setDiscountValue] = useState("");
  const [categories, setCategories] = useState([]);
  const [toast, setToast] = useState({ message: "", type: "" });

  const showToast = (message, type = "success") => {
    setToast({ message, type });
    setTimeout(() => setToast({ message: "", type: "" }), 3000);
  };

  const fetchCategoriesList = useCallback(async () => {
  try {
    const res = await api.get("/categories");
    setCategories(res?.data?.categories || []);
  } catch (err) {
    console.error("Failed to fetch categories", err);
    showToast("Failed to load categories", "error");
  }
}, []);


  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get("/products");
      setProducts(res.data.products || []);
    } catch (err) {
      console.error("Failed to fetch products", err);
      showToast("Failed to load products", "error");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);


  const createProduct = async (e) => {
    e.preventDefault();

    try {
      await api.post("/products", {
        product_name: name,
        product_description: desc,
        product_cost: parseFloat(cost),
        product_category_id: parseInt(categoryId),
        discount_type: discountType,
        discount_value: parseFloat(discountValue || 0),
      });

      showToast("Product created successfully", "success");
      setShowModal(false);

      setName("");
      setDesc("");
      setCost("");
      setCategoryId("");
      setDiscountType("");
      setDiscountValue("");
      fetchCategoriesList();

      fetchProducts();
    } catch (err) {
      console.error("Error while creating Product", err);
      showToast(
        err?.response?.data?.error || "Failed to create product",
        "error"
      );
    }
  };

  const startUpdating = (p) => {
    setUpdatingId(p.product_id);
    setName(p.product_name || "");
    setDesc(p.product_description || "");
    setCost(p.product_cost || "");
    setCategoryId(p.product_category_id || "");
    setDiscountType(p.discount_type || "");
    setDiscountValue(p.discount_value || "");
    fetchCategoriesList();
    setShowModal(true);
  };

  const updateProduct = async (e) => {
    e.preventDefault();

    try {
      await api.put(`/products/${updatingId}`, {
        product_name: name,
        product_description: desc,
        product_cost: parseFloat(cost),
        product_category_id: parseInt(categoryId),
        discount_type: discountType,
        discount_value: parseFloat(discountValue || 0),
      });

      showToast("Product updated successfully", "success");
      setShowModal(false);
      setUpdatingId(null);
      fetchProducts();
    } catch (err) {
      console.error("Error while Updating Product:", err);
      showToast(
        err?.response?.data?.error || "Failed to update product",
        "error"
      );
    }
  };

  const deleteProduct = async (id) => {
    try {
      await api.delete(`/products/${id}`);
      showToast("Product deleted successfully", "success");
      fetchProducts();
    } catch (err) {
      showToast("Failed to delete product", "error");
    }
  };

  const filteredProducts = useMemo(() => {
    if (!globalFilter) return products;
    return products.filter((p) =>
      JSON.stringify(p).toLowerCase().includes(globalFilter.toLowerCase())
    );
  }, [products, globalFilter]);

  return (
    <>
      <PageHeader>
        <Title>Products</Title>

        <SearchWrap>
          <Icon src={searchIcon} alt="search" />
          <SearchIpt
            type="text"
            placeholder="Search products..."
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
          />
        </SearchWrap>

        {user?.role === "supplier_admin" && (
          <AddButton
            onClick={() => {
              setUpdatingId(null);
              setName("");
              setDesc("");
              setCost("");
              setCategoryId("");
              setDiscountType("");
              setDiscountValue("");
              fetchCategoriesList();
              setShowModal(true);
            }}
          >
            <Icon src={plusIcon} alt="add" />
            Add Product
          </AddButton>
        )}
      </PageHeader>

      {loading ? (
        <p>Loading products...</p>
      ) : (
        <TableWrapper>
          <Table>
            <thead>
              <tr>
                <Th>ID</Th>
                <Th>Name</Th>
                <Th>Description</Th>
                <Th>Cost</Th>
                <Th align="center">Actions</Th>
              </tr>
            </thead>

            <tbody>
              {filteredProducts.length > 0 ? (
                filteredProducts.map((p) => (
                  <tr key={p.product_id}>
                    <Td>#{p.product_id}</Td>
                    <Td>{p.product_name}</Td>
                    <Td>{p.product_description || "-"}</Td>
                    <Td >{p.product_cost}</Td>
                    <Td align="center">
                      <ActionButtons>
                        {user?.role === "supplier_admin" && (
                          <>
                            <IconButton onClick={() => startUpdating(p)}>
                              <Icon src={editIcon} alt="edit" />
                            </IconButton>

                            <IconButton
                              $danger
                              onClick={() => deleteProduct(p.product_id)}
                            >
                              <Icon src={trashIcon} alt="delete" />
                            </IconButton>
                          </>
                        )}
                      </ActionButtons>
                    </Td>
                  </tr>
                ))
              ) : (
                <tr>
                  <Td colSpan="5" align="center">
                    No products found
                  </Td>
                </tr>
              )}
            </tbody>
          </Table>
        </TableWrapper>
      )}

      {user?.role === "supplier_admin" && showModal && (
        <Modal onClick={() => setShowModal(false)}>
          <ModalContent onClick={(e) => e.stopPropagation()}>
            <ModalHeader>
              <ModalTitle>
                {isUpdating ? "Edit Product" : "Add New Product"}
              </ModalTitle>

              <CloseButton onClick={() => setShowModal(false)}>
                <Icon src={xIcon} alt="close" />
              </CloseButton>
            </ModalHeader>

            <Form onSubmit={isUpdating ? updateProduct : createProduct}>
              <FormGroup>
                <Label>Product Name *</Label>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Cost *</Label>
                <Input
                  type="number"
                  step="0.01"
                  value={cost}
                  onChange={(e) => setCost(e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Category *</Label>
                <Select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} required>
                   <option value="">Select Category</option>
                   {categories.map((c) => (
                    <option key={c.category_id} value={c.category_id}>
                      {c.category_name}
                    </option>
                   ))}
                </Select>
              </FormGroup>

              <FormGroup>
                <Label>Description</Label>
                <Textarea
                  rows="3"
                  value={desc}
                  onChange={(e) => setDesc(e.target.value)}
                />
              </FormGroup>

              <FormGroup>
                <Label>Discount Type</Label>
                <Select
                  value={discountType}
                  onChange={(e) => setDiscountType(e.target.value)}
                >
                  <option value="">No Discount</option>
                  <option value="flat">Flat</option>
                  <option value="percent">Percent</option>
                </Select>
              </FormGroup>

              <FormGroup>
                <Label>Discount Value</Label>
                <Input
                  type="number"
                  step="0.01"
                  value={discountValue}
                  onChange={(e) => setDiscountValue(e.target.value)}
                />
              </FormGroup>

              <ModalActions>
                <CancelButton type="button" onClick={() => setShowModal(false)}>
                  Cancel
                </CancelButton>

                <SubmitButton type="submit">
                  {isUpdating ? "Update" : "Create"}
                </SubmitButton>
              </ModalActions>
            </Form>
          </ModalContent>
        </Modal>
      )}
      <ToastContainer>
        {toast.message && <Toast $type={toast.type}>{toast.message}</Toast>}
      </ToastContainer>
    </>
  );
};

export default Products;
