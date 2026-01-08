import React, { useState, useEffect, useMemo } from 'react';
import api from '../../Api/axios';
import { useCallback } from 'react';
import { getUserInfo } from '../../utils/auth';
import {
  PageHeader,
  Title,
  AddButton,
  Icon,
  SearchWrap,
  SearchIpt,
  ActionButtons,
  IconButton,
} from './Styles'
import { toast, ToastContainer } from 'react-toastify';
import Datatable from '../DataTable';
import AddOrEditSupplierModal from './AddorEditModal';




import plusIcon from '../../Images/plus.svg'
import editIcon from '../../Images/edit.svg'
import trashIcon from '../../Images/trash.svg'
import searchIcon from '../../Images/search.svg'
import { useLocation } from 'react-router-dom';


export const caseInsensitiveSort = (rowA, rowB, columnId) => {
  const valueA = String(rowA.getValue(columnId) ?? "").toLowerCase();
  const valueB = String(rowB.getValue(columnId) ?? "").toLowerCase();

  if (valueA > valueB) return 1;
  if (valueA < valueB) return -1;
  return 0;
};


const Suppliers = () => {
  const user = getUserInfo()
  const [modalOpen, setModalOpen] = useState(false)
  const [loading, setLoading] = useState(true)
  const [globalFilter, setGlobalFilter] = useState('')
  const [state, setState ] = useState({
   data: [],
   updatingId: null,
   form: {name: "", contact: "", email: "", company: ""}
 })
  

  const fetchSuppliers = useCallback(async () => {
  try {
    setLoading(true)
    const res = await api.get('/suppliers')
     setState((prev) => ({...prev, data: res.data.suppliers || []}))     
  } catch (err) {
    console.log('failed to fetch suppliers: ', err)
    toast.error("Failed to load suppliers")
  } finally {
    setLoading(false)
  }
}, [])

  useEffect(() => {
    fetchSuppliers()
  }, [fetchSuppliers])

  const location = useLocation()
  useEffect(() => {
    const getFromBrowser = new URLSearchParams(window.location.search)
    const id = getFromBrowser.get("id")
    const name = getFromBrowser.get("name")
    const contact = getFromBrowser.get("contact_info")
    const email = getFromBrowser.get("email")
    const company = getFromBrowser.get("company")

    const value = id || name || contact || email || company || ""
    if(value){
        setGlobalFilter(value)
    }
  }, [location.search])

  useEffect(() => {
    const handler = setTimeout(async () => {
        const search = globalFilter.trim()

        switch(true){
            case search === "":
                fetchSuppliers()
                return
            case /^\d+$/.test(search):
                try{
                    const res = await api.get(`/suppliers/${search}`)
                    setState(prev => ({
                        ...prev,
                        data: [res.data.supplier]
                    }))
                } catch{
                    setState(prev => ({...prev, data: []}))
                }
                return 
            default:
                try {
                    const res = await api.get("/suppliers",{
                        params: {name: search, company: search}
                    })

                    setState(prev => ({
                        ...prev,
                        data: res.data.suppliers || []
                    }))
                } catch (err){
                    console.log(err.res?.data || err.message)
                    toast.error("Failed to search suppliers")
                }

        }

    }, 350)
    return () => clearTimeout(handler)
  }, [globalFilter, fetchSuppliers])


  const deleteSupplier = useCallback(async (id) => {
  try {
    await api.delete(`suppliers/${id}`)
    toast.success('Supplier deleted successfully');
    fetchSuppliers()
  } catch {
    toast.error('Failed to delete supplier');
  }
}, [fetchSuppliers])


const onClickEdit = useCallback((row) => {
    setState((prev) => ({
      ...prev,
      updatingId: row.original.supplier_id,
      form: {
        name: row.original.name,
        contact: row.original.contact_info,
        email: row.original.email,
        company: row.original.company,
      },
    }));
    setModalOpen(true);
  }, [])

  const onClickAdd = () => {
    setState((prev) => ({
      ...prev,
      updatingId: null,
      form: { name: "", contact: "", email: "", company: "" },
    }));
    setModalOpen(true);
  }

  const onModalClose = () => {
    setModalOpen(false);
    setState((prev) => ({ ...prev, updatingId: null }));
  }

  const columns = useMemo(
    () => [
      {
        accessorKey: 'supplier_id',
        header: 'ID',
        cell: (info) => `#${info.getValue()}`,
      },
      {
        accessorKey: 'name',
        header: 'Name',
        sortingFn: caseInsensitiveSort,
      },
      {
        accessorKey: 'contact_info',
        header: 'Contact',
      },
      {
        accessorKey: 'email',
        header: 'Email',
        enableSorting: false,
      },
      {
        accessorKey: 'company',
        header: 'Company',
        sortingFn: caseInsensitiveSort,
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) =>
          user.role === "system_admin" && (
            <ActionButtons>
              <IconButton onClick={() => onClickEdit(row)}>
                <Icon src={editIcon} alt="edit" />
              </IconButton>

              <IconButton
                onClick={() => deleteSupplier(row.original.supplier_id)}
              >
                <Icon src={trashIcon} alt="delete" />
              </IconButton>
            </ActionButtons>
          ),
      },
    ],
    [user.role, onClickEdit, deleteSupplier]
  )
  return (
        <>
          <PageHeader>
            <Title>Suppliers</Title>
                <SearchWrap>
                    <Icon src={searchIcon} alt="search" />
                        <SearchIpt
                            value={globalFilter}
                            onChange={(e) => setGlobalFilter(e.target.value)}
                            placeholder="Search Suppliers"
                        />
                </SearchWrap>
                {user?.role === 'system_admin' && (
                    <AddButton onClick={onClickAdd}>
                        <Icon src={plusIcon} alt="add" />
                        Add Supplier
                    </AddButton>
                )}
          </PageHeader>

          {loading ? (
            <p>Loading suppliers...</p>
          ) : (
            <Datatable
          data={state.data}
          columns={columns}
          globalFilter={globalFilter}
          setGlobalFilter={setGlobalFilter}
          />
      )}

       <AddOrEditSupplierModal
        open={modalOpen}
        onClose={onModalClose}
        updatingId={state.updatingId}
        refetch={fetchSuppliers}
      />
      <ToastContainer position="top-right" autoClose={3000} />
  </>

  )
}

export default Suppliers