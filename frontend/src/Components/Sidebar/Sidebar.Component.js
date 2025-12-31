import { Sidebar, SidebarHeader, Nav, NavItem, LogoutWrapper,Icon, LogoutButton } from "./Sidebar.styles"
import packageIcon from '../Dashboard/assets1/package.svg'
import dashboardIcon from '../Dashboard/assets1/layout-dashboard.svg'
import folderIcon from '../Dashboard/assets1/folder-tree.svg'
import usersIcon from '../Dashboard/assets1/users.svg'
import truckIcon from '../Dashboard/assets1/truck.svg'
import trendingIcon from '../Dashboard/assets1/trending-up.svg'
import logoutIcon from '../Dashboard/assets1/log-out.svg'
import { useNavigate, useLocation } from "react-router-dom"
import { getUserInfo } from "../../utils/auth"


const SidebarMenu = ({ sidebarOpen}) => {
  const navigate = useNavigate()
  const user = getUserInfo()
  const location = useLocation()

  const menuItems = user?.role === 'system_admin'
    ? [
        { icon: dashboardIcon, label: 'Dashboard', path: '/admin/dashboard' },
        { icon: packageIcon, label: 'Products', path: '/admin/products' },
        { icon: folderIcon, label: 'Categories', path: '/admin/categories'},
        { icon: truckIcon, label: 'Suppliers', path: '/admin/suppliers' },
        { icon: trendingIcon, label: 'Stock Movements', path: '/admin/stock-movements' },
        { icon: usersIcon, label: 'Users', path: '/admin/users' },
      ]
    : [
        { icon: dashboardIcon, label: 'Dashboard', path: '/admin/dashboard' },
        { icon: packageIcon, label: 'Products', path: '/admin/products' },
        { icon: folderIcon, label: 'Categories', path: '/admin/categories' },
        { icon: truckIcon, label: 'Suppliers', path: '/admin/suppliers' },
        { icon: trendingIcon, label: 'Stock Movements', path: '/admin/stock-movements' },
      ]

  const handleMenuClick = (item) => {
    navigate(item.path)
  }

  const logout = () => navigate('/')
  
  return (
    <Sidebar $isOpen={sidebarOpen}>
      <SidebarHeader>
        <Icon src={packageIcon} alt="logo" />
        Product Catalogue
      </SidebarHeader>

      <Nav>
        {menuItems.map((item) => (
          <NavItem
            key={item.label}
            $active={location.pathname === item.path}
            onClick={() => handleMenuClick(item)}
          >
            <Icon src={item.icon} alt={item.label} />
            <span>{item.label}</span>
          </NavItem>
        ))}
      </Nav>

      <LogoutWrapper>
        <LogoutButton onClick={logout}>
          <Icon src={logoutIcon} alt="logout" />
          <span>Logout</span>
        </LogoutButton>
      </LogoutWrapper>
    </Sidebar>
  )
}

export default SidebarMenu


