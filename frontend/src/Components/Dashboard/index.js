import React, { useEffect, useState , useMemo} from "react";
import api from "../../Api/axios";
import { useNavigate } from "react-router-dom";
import Datatable from "../DataTable";
import { getUserInfo } from "../../utils/auth";
import { toast , ToastContainer} from "react-toastify";
import {
  PageHeader,
  Title,
  StatsGrid,
  StatCard,
  StatIcon,
  StatInfo,
  StatValue,
  StatLabel,
  Section,
  SectionHeader,
  SectionTitle,
  ThresholdBadge,
  Icon,
} from "./Styles";
import packageIcon from "../../Images/package.svg";
import folderIcon from "../../Images/folder-tree.svg";
import usersIcon from "../../Images/users.svg";
import truckIcon from "../../Images/truck.svg";
import alertIcon from "../../Images/alert-triangle.svg";

export const caseInsensitiveSort = (rowA, rowB, columnId) => {
  const valueA = String(rowA.getValue(columnId) ?? "").toLowerCase();
  const valueB = String(rowB.getValue(columnId) ?? "").toLowerCase();

  if (valueA > valueB) return 1;
  if (valueA < valueB) return -1;
  return 0;
}

const Dashboard = () => {
  const [stats, setStats] = useState({})
  const [lowStock, setLowStock] = useState([])
  const [threshold, setThreshold] = useState(10)
  const user = getUserInfo()

  const navigate = useNavigate()

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const res = await api.get("/dashboard")

        setStats(res.data);
        setLowStock(res.data.low_stock_products || [])
        setThreshold(res.data.low_stock_threshold || 10)
      } catch (err) {
        console.error("Failed to load dashboard", err)
      }
    }

    fetchDashboard();
  }, [])

  const cards = [
    {
      icon: packageIcon,
      label: "Total Products",
      value: stats.total_products || 0,
      path: "/admin/products",
    },
    {
      icon: folderIcon,
      label: "Categories",
      value: stats.total_categories || 0,
      path: "/admin/categories",
    },
    {
      icon: truckIcon,
      label: "Suppliers",
      value: stats.total_suppliers || stats.company_suppliers || 0,
      path: "/admin/suppliers",
    },
    {
      icon: usersIcon,
      label: "Supplier Admins",
      value: stats.total_supplier_admins || 0,
      path: "/admin/users",
      restricted: user.role !== "system_admin"
    },
  ] 

  const lowStockColumns = useMemo(() => [
    {
      accessorKey: "product_id",
      header: "Product ID",
      cell: (info) => `#${info.getValue()}`,
    },
    {
      accessorKey: "product_name",
      header: "Product Name",
      sortingFn: caseInsensitiveSort
    },
    {
      accessorKey: "current_stock",
      header: "Stock",
      cell: (info) => (
        <span style={{ color: "#c75858", fontWeight: 600 }}>
          {info.getValue()}
        </span>
      ),
    },
  ], [])

  return (
    <>
      <PageHeader>
        <Title data-cy="dashboard-title">Dashboard</Title>
      </PageHeader>

      <StatsGrid>
        {cards.map((card, i) => (
          <StatCard key={i} data-cy={`dashboard-card-${card.label.toLowerCase().replace(" ", "-")}`} onClick={() => {if (card.restricted){toast.error("access denied"); return} navigate(card.path)}}>
            <StatIcon>
              <Icon src={card.icon} alt={card.label} />
            </StatIcon>
            <StatInfo>
              <StatValue>{card.value}</StatValue>
              <StatLabel>{card.label}</StatLabel>
            </StatInfo>
          </StatCard>
        ))}
      </StatsGrid>

      <Section>
        <SectionHeader>
          <SectionTitle data-cy="low-stock-title">
            <Icon src={alertIcon} alt="alert" />
            Low Stock Alert
          </SectionTitle>
          <ThresholdBadge data-cy="low-stock-threshold">Threshold: {threshold} units</ThresholdBadge>
        </SectionHeader>

        <Datatable
          data-cy="low-stock-table"
          data={lowStock}
          columns={lowStockColumns}
          globalFilter=""
          setGlobalFilter={() => {}}
        />
        <ToastContainer position="top-right" autoClose={3000} />
      </Section>
    </>
  );
};

export default Dashboard
