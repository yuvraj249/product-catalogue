import React, { useState, useEffect, useMemo } from 'react';
import api from '../../api/axios';
import { useCallback } from 'react';
import { getUserInfo } from '../../utils/auth';
import { useDebounce } from '../../DebounceSearch';
import {
  PageHeader,
  Title,
  Icon,
  SearchWrap,
  SearchIpt,
  ActionButtons,
  IconButton,
  HeaderAction,
  HeaderAddButton
} from '../CommonUIStyles/Styles'
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
  

  const fetchSuppliers = useCallback(async (q = "") => {
    try {
      setLoading(true)
         
        const res = await api.get('/suppliers', { params: { q } })
        setState(prev => ({ ...prev, data: res.data.suppliers || [] }))

    } catch {
      toast.error('Failed to load suppliers')
    } finally {
      setLoading(false)
    }
  }, [])

  const location = useLocation()
  useEffect(() => {
  const params = new URLSearchParams(location.search)
  const value = params.get("q")

  if (value !== null) {
    setGlobalFilter(value)
  } else {
    fetchSuppliers()
  }
}, [location.search, fetchSuppliers])


  const debouncedSearch = useDebounce(globalFilter, 350);
 useEffect(() => {
  fetchSuppliers(debouncedSearch.trim());
}, [debouncedSearch, fetchSuppliers]);
  const deleteSupplier = useCallback(async (id) => {
  try {
    await api.delete(`suppliers/${id}`)
    toast.success('Supplier deleted successfully');
    fetchSuppliers(globalFilter)
  } catch {
    toast.error('Failed to delete supplier');
  }
}, [fetchSuppliers,globalFilter])


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
    }))
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
              <IconButton data-cy="edit-supplier-btn" onClick={() => onClickEdit(row)}>
                <Icon src={editIcon} alt="edit" />
              </IconButton>

              <IconButton
                data-cy="delete-supplier-btn" 
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
                            data-cy="supplier-search"
                            value={globalFilter}
                            onChange={(e) => setGlobalFilter(e.target.value)}
                            placeholder="Search Suppliers"
                        />
                </SearchWrap>
                {user?.role === 'system_admin' && (
                <HeaderAction><HeaderAddButton data-cy="add-supplier-btn" onClick={onClickAdd}>
                        <Icon src={plusIcon} alt="add" />
                        Add Supplier
                </HeaderAddButton></HeaderAction>
                )}
          </PageHeader>

          {loading ? (
            <p>Loading suppliers...</p>
          ) : (
            <Datatable
          data={state.data}
          columns={columns}
          />
      )}

       <AddOrEditSupplierModal
        open={modalOpen}
        onClose={onModalClose}
        updatingId={state.updatingId}
        refetch={fetchSuppliers}
        setLoading={setLoading}
      />
      <ToastContainer position="top-right" autoClose={3000} />
  </>

  )
}

export default Suppliers