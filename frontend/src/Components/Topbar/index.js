import { TopBar, MenuButton, UserRole, Icon } from "../Sidebar/Styles";
import menuIcon from '../../Images/menu.svg';
import { getUserInfo } from "../../utils/auth";

const TopbarMenu = ({ sidebarOpen, setSidebarOpen }) => {
  const user = getUserInfo();
  const displayName = user.name
  console.log(getUserInfo())

  return (
    <TopBar $sidebarOpen={sidebarOpen}>
      <MenuButton onClick={() => setSidebarOpen(!sidebarOpen)}>
        <Icon src={menuIcon} alt="menu" />
      </MenuButton>

      <UserRole>
         Hey! {displayName}
      </UserRole>
    </TopBar>
  );
};

export default TopbarMenu;
