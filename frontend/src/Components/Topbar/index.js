import { TopBar, MenuButton, UserRole, Icon } from "../Sidebar/Styles";
import menuIcon from '../../Images/menu.svg';
import { getUserInfo } from "../../utils/auth";
import { useLocation } from "react-router-dom";
import styled from "styled-components";

const LeftSection = styled.div`
  display: flex;
  align-items: center;
  gap: 16px;
`;

const Breadcrumb = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 15px;
  font-weight: 600;
  color: #f1f5f9;

  span.root {
    color: #64748b;
    font-weight: 500;
  }

  span.slash {
    color: #475569;
  }
`;

const UserProfile = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const Avatar = styled.div`
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: linear-gradient(135deg, #38bdf8 0%, #3b82f6 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-weight: 700;
  font-size: 14px;
  box-shadow: 0 4px 12px rgba(56, 189, 248, 0.25);
`;

const UserDetails = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-end;

  @media (max-width: 640px) {
    display: none;
  }
`;

const UserName = styled.span`
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 13.5px;
  font-weight: 600;
  color: #f8fafc;
`;

const getPathTitle = (pathname) => {
  if (pathname.includes("products")) return "Products Management";
  if (pathname.includes("categories")) return "Categories";
  if (pathname.includes("suppliers")) return "Suppliers Hub";
  if (pathname.includes("stock-movements")) return "Stock Movements";
  if (pathname.includes("users")) return "User Management";
  return "Overview";
};

const TopbarMenu = ({ sidebarOpen, setSidebarOpen }) => {
  const user = getUserInfo();
  const location = useLocation();
  const displayName = user?.name || "User";
  const initials = displayName.split(" ").map(n => n[0]).join("").toUpperCase().slice(0, 2);
  const title = getPathTitle(location.pathname);

  return (
    <TopBar $sidebarOpen={sidebarOpen}>
      <LeftSection>
        <MenuButton onClick={() => setSidebarOpen(!sidebarOpen)} aria-label="Toggle Navigation Sidebar">
          <Icon src={menuIcon} alt="menu" />
        </MenuButton>

        <Breadcrumb>
          <span className="root">Nexus</span>
          <span className="slash">/</span>
          <span>{title}</span>
        </Breadcrumb>
      </LeftSection>

      <UserProfile>
        <UserDetails>
          <UserName>{displayName}</UserName>
        </UserDetails>
        <Avatar>{initials}</Avatar>
        <UserRole>
          {user?.role === "system_admin" ? "Admin" : "Supplier"}
        </UserRole>
      </UserProfile>
    </TopBar>
  );
};

export default TopbarMenu;
