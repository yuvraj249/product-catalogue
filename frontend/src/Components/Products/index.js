import React, { useState, useEffect, useCallback, useMemo } from "react";
import api from "../../Api/axios";
import { getUserInfo } from "../../utils/auth";
import {
  PageHeader,
  Title,
  AddButton,
  Icon,
  SearchWrap,
  SearchIpt,
  ActionButtons,
  IconButton,
} from "./Styles";
import { toast, ToastContainer } from "react-toastify";
import Datatable from "../DataTable";
import AddOrEditProductModal from "./AddorEditModal";
import plusIcon from "../../Images/plus.svg";
import editIcon from "../../Images/edit.svg";
import trashIcon from "../../Images/trash.svg";
import searchIcon from "../../Images/search.svg";
import { useLocation } from "react-router-dom";

export const caseInsensitiveSort = (rowA, rowB, columnId) => {
  const valueA = String(rowA.getValue(columnId) ?? "").toLowerCase();
  const valueB = String(rowB.getValue(columnId) ?? "").toLowerCase();

  if (valueA > valueB) return 1;
  if (valueA < valueB) return -1;
  return 0;
}

const Products = () => {
  const user = getUserInfo()
  const [loading, setLoading] = useState(true)
  const [globalFilter, setGlobalFilter] = useState("")
  const [modalOpen, setModalOpen] = useState(false)
  const [state, setState] = useState({
    data: [],
    updatingId: null,
    form: {},
  })

  const [categories, setCategories] = useState([]);
  useEffect(() => {
  api.get("/categories")
    .then(res => setCategories(res.data.categories || []))
    .catch(() => toast.error("Failed to load categories"));
}, [])
  
    const categoryMap = useMemo(() => {
    const map = {};
    categories.forEach(c => {
      map[c.category_id] = c.category_name;
    });
    return map;
  }, [categories])

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get("/products")
      setState((p) => ({ ...p, data: res.data.products || [] }))
    } catch {
      toast.error("Failed to load products")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchProducts()
  }, [fetchProducts])

  const location = useLocation()
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
      const id = params.get("id")
      const product_name = params.get("product_name") 
      const category = params.get("category") 
      const value = id || product_name || category || ""
    if(value){
        setGlobalFilter(value)
    }
  }, [location.search])

  useEffect(() => {
    const handler = setTimeout(async () => {
        const search = globalFilter.trim()

        switch(true){
            case search === "":
                fetchProducts()
                return
            case /^\d+$/.test(search):
                try{
                    const res = await api.get(`/products/${search}`)
                    setState(prev => ({
                        ...prev,
                        data: [res.data.product]
                    }))
                } catch{
                    setState(prev => ({...prev, data: []}))
                }
                return 
            default:
                try {
                    const res = await api.get("/products",{
                        params: {product_name: search, category: search }
                    })

                    setState(prev => ({
                        ...prev,
                        data: res.data.products || []
                    }))
                } catch (err){
                    console.log(err.res.data || err.message)
                    toast.error("Failed to search suppliers")
                }

        }

    }, 350)
    return () => clearTimeout(handler)
  }, [globalFilter, fetchProducts])

  const deleteProduct = useCallback(
    async (id) => {
      try {
        await api.delete(`/products/${id}`)
        toast.success("Product deleted")
        fetchProducts()
      } catch {
        toast.error("Failed to delete product")
      }
    },
    [fetchProducts]
  );

  const onClickEdit = useCallback((row) => {
    setState((prev) => ({
    ...prev,      
    updatingId: row.original.product_id,
    form: row.original,
   }))
    setModalOpen(true);
  }, [])

  const onClickAdd = () => {
    setState({ data: [], updatingId: null, form: {} })
    setModalOpen(true)
  }

  const columns = useMemo(
    () => [
      { accessorKey: "product_id", 
        header: "ID", 
        cell: (info) => `#${info.getValue()}` 
      },
      { accessorKey: "product_name", 
        header: "Name",
        sortingFn: caseInsensitiveSort
      },
      { accessorKey: "product_cost", 
        header: "Cost" 
      },
      { accessorKey: "discount_type", 
        header: "Discount Type",
        sortingFn: caseInsensitiveSort
      },
      { accessorKey: "discount_value", 
        header: "Discount Value" 
      },
      {
        accessorKey: "product_category_id",
        header: "Category",
        cell: (info) => categoryMap[info.getValue()] || "—",
        sortingFn: caseInsensitiveSort
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) =>
          user.role === "supplier_admin" && (
            <ActionButtons>
              <IconButton onClick={() => onClickEdit(row)}>
                <Icon src={editIcon} />
              </IconButton>
              <IconButton onClick={() => deleteProduct(row.original.product_id)}>
                <Icon src={trashIcon} />
              </IconButton>
            </ActionButtons>
          ),
      },
    ],
    [user.role, onClickEdit, deleteProduct, categoryMap]
  )

  return (
    <>
      <PageHeader>
        <Title>Products</Title>
        <SearchWrap>
          <Icon src={searchIcon} />
          <SearchIpt
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Search Products"
          />
        </SearchWrap>

        {user.role === "supplier_admin" && (
          <AddButton onClick={onClickAdd}>
            <Icon src={plusIcon} />
            Add Product
          </AddButton>
        )}
      </PageHeader>

      {loading ? (
        <p>Loading products…</p>
      ) : (
        <Datatable
          data={state.data}
          columns={columns}
          globalFilter={globalFilter}
          setGlobalFilter={setGlobalFilter}
        />
      )}

      <AddOrEditProductModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        updatingId={state.updatingId}
        form={state.form}
        refetch={fetchProducts}
        categories={categories}
      />

      <ToastContainer />
    </>
  )
}

export default Products;
