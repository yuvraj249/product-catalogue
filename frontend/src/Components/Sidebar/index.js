import { Sidebar, SidebarHeader, Nav, NavItem, LogoutWrapper,Icon, LogoutButton } from "./Styles"
import packageIcon from '../../Images/package.svg'
import dashboardIcon from '../../Images/layout-dashboard.svg'
import folderIcon from '../../Images/folder-tree.svg'
import usersIcon from '../../Images/users.svg'
import truckIcon from '../../Images/truck.svg'
import trendingIcon from '../../Images/trending-up.svg'
import logoutIcon from '../../Images/log-out.svg'
import { useNavigate, useLocation } from "react-router-dom"
import { getUserInfo, handleLogout } from "../../utils/auth"


const user = getUserInfo()

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


const SidebarMenu = ({ sidebarOpen}) => {
  const navigate = useNavigate()
  const location = useLocation()
  const handleMenuClick = (item) => {
    navigate(item.path)
  }
 
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
        <LogoutButton onClick={handleLogout}>
          <Icon src={logoutIcon} alt="logout" />
          <span>Logout</span>
        </LogoutButton>
      </LogoutWrapper>
    </Sidebar>
  )
}

export default SidebarMenu


