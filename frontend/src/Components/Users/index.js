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
  IconButton,
  ActionButtons,
} from "./Styles";
import { toast, ToastContainer } from "react-toastify";
import plusIcon from "../../Images/plus.svg";
import searchIcon from "../../Images/search.svg";
import editIcon from "../../Images/edit.svg";
import trashIcon from "../../Images/trash.svg";
import Datatable from "../DataTable";
import AddOrEditUserModal from "./AddorEditModal";
import { useLocation } from "react-router-dom";

const Users = () => {
  const user = getUserInfo()
  const [loading, setLoading] = useState(true)
  const [globalFilter, setGlobalFilter] = useState("")
  const [modalOpen, setModalOpen] = useState(false)
  const [suppliers, setSuppliers] = useState([])

useEffect(() => {
  api.get("/suppliers")
    .then(res => setSuppliers(res.data.suppliers || []))
    .catch(() => toast.error("Failed to load suppliers"))
}, [])

const supplierMap = useMemo(() => {
  const map = {};
  suppliers.forEach(s => {
    map[s.supplier_id] = s.name || s.company || `Supplier #${s.supplier_id}`
  })
  return map
}, [suppliers])


  const [state, setState] = useState({
    data: [],
    updatingId: null,
    form: { name: "", email: "", password: "",role: "supplier_admin", supplier_id: "" },
  })

  
  const fetchUsers = useCallback(async (q = "") => {
    try {
      setLoading(true)
      const res = await api.get("/users/supplier-admin", {
        params: { q },
      });
      setState((prev) => ({ ...prev, data: res.data.users || [] }))
    } catch {
      toast.error("Failed to load users")
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
            fetchUsers()
    } else {
        fetchUsers(globalFilter.trim())
    }
    }, 350);

    return () => clearTimeout(handler)
  }, [globalFilter, fetchUsers])

 
  
  const deleteUser = useCallback(
    async (id) => {
      try {
        await api.delete(`users/supplier-admin/${id}`)
        toast.success("User deleted")
        fetchUsers(globalFilter)
      } catch {
        toast.error("Failed to delete user")
      }
    },
    [fetchUsers,globalFilter]
  )


  const onClickEdit = useCallback((row) => {
    setState((prev) => ({
      ...prev,
      updatingId: row.original.user_id,
      form: {
        name: row.original.name,
        email: row.original.email,
        role: row.original.role,
        supplier_id: row.original.supplier_id || "",
        password: "",
      },
    }))
    setModalOpen(true);
  }, [])

 
  const onClickAdd = () => {
    setState((prev) => ({
      ...prev,
      updatingId: null,
      form: { name: "", email: "", password: "", supplier_id: "" },
    }));
    setModalOpen(true)
  };

  const onModalClose = () => {
    setModalOpen(false)
    setState((prev) => ({ ...prev, updatingId: null }))
  };

 
  const columns = useMemo(
    () => [
      {
        accessorKey: "user_id",
        header: "ID",
        cell: (info) => `#${info.getValue()}`,
      },
      {
        accessorKey: "name",
        header: "Name",
      },
      {
        accessorKey: "email",
        header: "Email",
      },
      {
        accessorKey: "role",
        header: "Role",
      },
      {
        accessorKey: "supplier_id",
        header: "Supplier",
        cell: (info) => supplierMap[info.getValue()] || "—",
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) =>
          user.role === "system_admin" && (
            <ActionButtons>
              <IconButton data-cy="edit-user-btn" onClick={() => onClickEdit(row)}>
                <Icon src={editIcon} alt="edit" />
              </IconButton>

              <IconButton data-cy="delete-user-btn" onClick={() => deleteUser(row.original.user_id)}>
                <Icon src={trashIcon} alt="delete" />
              </IconButton>
            </ActionButtons>
          ),
      },
    ],
    [user.role, onClickEdit, deleteUser,supplierMap]
  );

console.log({state})
  return (
    <>
      <PageHeader>
        <Title>Users</Title>

        <SearchWrap>
          <Icon src={searchIcon} alt="search" />
          <SearchIpt
            data-cy="user-search"
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Search Users"
          />
        </SearchWrap>

        {user.role === "system_admin" && (
          <AddButton data-cy="add-user-btn" onClick={onClickAdd}>
            <Icon src={plusIcon} alt="add" />
            Add User
          </AddButton>
        )}
      </PageHeader>

      {loading ? (
        <p style={{ padding: 20 }}>Loading users…</p>
      ) : (
        <Datatable
          data={state.data}
          columns={columns}
        />
      )}

      <AddOrEditUserModal
        open={modalOpen}
        onClose={onModalClose}
        updatingId={state.updatingId}
        form={state.form}
        refetch={fetchUsers}
        suppliers={suppliers}
        setLoading={setLoading}
      />

      <ToastContainer position="top-right" autoClose={3000} />
    </>
  );
};

export default Users
