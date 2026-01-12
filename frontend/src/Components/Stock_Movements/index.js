import React, { useState, useEffect, useCallback, useMemo } from "react";
import api from "../../Api/axios";
import { getUserInfo } from "../../utils/auth";
import {
  PageHeader,
  Title,
  SearchWrap,
  SearchIpt,
  Icon,
  AddButton,
  ActionButtons,
  IconButton,

} from "./Styles";
import { toast, ToastContainer } from "react-toastify";
import Datatable from "../DataTable";
import searchIcon from "../../Images/search.svg";
import plusIcon from '../../Images/plus.svg';
import AddOrEditStockModal from "./AddorEditModal";
import editIcon from '../../Images/edit.svg';
import trashIcon from '../../Images/trash.svg';
import { useLocation } from "react-router-dom";



const StockMovements = () => {
  const user = getUserInfo()
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [globalFilter, setGlobalFilter] = useState("")
  const [data, setData] = useState([])
  const [products, setProducts] = useState([])
  const [editingItem, setEditingItem] = useState(null)

  useEffect(() => {
    api.get("/products")
    .then(res => setProducts(res.data.products || []))
    .catch(() => toast.error("failed to load products"))
  },[])


  const fetchMovements = useCallback(async (q = "") => {
    try {
      setLoading(true);
      const res = await api.get("/stock_movements", { params: { q } })
      setData(res.data.movements || [])
    } catch {
      toast.error("failed to load stock movements")
    } finally {
      setLoading(false)
    }
  }, [])

  const location = useLocation()

  useEffect(() => {
    const qs = new URLSearchParams(window.location.search)
    const value =  qs.get("q") || ""

    if (value) setGlobalFilter(value)
  }, [location.search]);


   useEffect(() => {
    const handler = setTimeout(() => {
      if (globalFilter.trim() === "") {
            fetchMovements()
    } else {
        fetchMovements(globalFilter.trim())
    }
    }, 350)
    return () => clearTimeout(handler)
  }, [globalFilter, fetchMovements])



  const deleteMovement = useCallback(async (id) => {
    try {
      await api.delete(`/stock_movements/${id}`)
      toast.success("stock movement deleted")
      fetchMovements();
    } catch (err) {
      toast.error(err.response?.data?.error || "delete failed")
    }
  },[fetchMovements])

  const columns = useMemo(
    () => [
      { accessorKey: "stock_id",
        header: "ID", 
        cell: (info) => `#${info.getValue()}`
      },
      {
        accessorKey: "product_name",
        header: "Product",
      },
      { 
        accessorKey: "quantity",
        header: "Qty" 
      },
      { 
        accessorKey: "movement_type",
        header: "Type" 
      },
      { 
        accessorKey: "reason",
        header: "Reason" 
      },
      ...( user.role === "system_admin" ? [
      { 
        accessorKey: "username", 
        header: "User" ,
      },
    ] : []
    ),
    ...( user.role === "supplier_admin" ? [
      {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        user.role === "supplier_admin" && (
          <ActionButtons>
            <IconButton
              onClick={() => {
                setEditingItem(row.original)
                setModalOpen(true)
              }}
            >
              <Icon src={editIcon} />
            </IconButton>

            <IconButton onClick={() => deleteMovement(row.original.stock_id)}>
              <Icon src={trashIcon} />
            </IconButton>
          </ActionButtons>
        )
      )
    } ] : []
)
    ],
    [user.role, deleteMovement]
  ) 



  return (
    <>
      <PageHeader>
        <Title>Stock Table</Title>
        <SearchWrap>
          <Icon src={searchIcon} />
          <SearchIpt
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Search stock movements"
          />
        </SearchWrap>

        {user.role === "supplier_admin" && (
            <AddButton onClick={() => {setEditingItem(null); setModalOpen(true)}}>
            <Icon src={plusIcon} />
            Add Stock
            </AddButton>
  )}
      </PageHeader>

      {loading ? (
        <p>Loading stock movements…</p>
      ) : (
        <Datatable
          data={data}
          columns={columns}
          globalFilter={globalFilter}
          setGlobalFilter={setGlobalFilter}
        />
      )}
      <AddOrEditStockModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        refetch={fetchMovements}
        products={products}
        editingItem={editingItem}
    />

      <ToastContainer position="top-right" autoClose={3000}/>
    </>
  )
}

export default StockMovements
