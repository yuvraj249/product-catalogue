import { TopBar, MenuButton, UserRole, Icon } from "../Sidebar/Sidebar.styles";
import menuIcon from '../Dashboard/assets1/menu.svg';
import { getUserInfo } from "../../utils/auth";

const TopbarMenu = ({ sidebarOpen, setSidebarOpen }) => {
  const user = getUserInfo();

  return (
    <TopBar $sidebarOpen={sidebarOpen}>
      <MenuButton onClick={() => setSidebarOpen(!sidebarOpen)}>
        <Icon src={menuIcon} alt="menu" />
      </MenuButton>

      <UserRole>
        {user?.role || "User"}
      </UserRole>
    </TopBar>
  );
};

export default TopbarMenu;
