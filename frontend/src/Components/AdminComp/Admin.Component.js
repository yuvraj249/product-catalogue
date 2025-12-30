import { useState } from "react";
import { Outlet } from "react-router-dom";
import SidebarMenu from "../Sidebar/Sidebar.Component";
import TopbarMenu from "../Topbar/Topbar.Component";
import { Main, ContentArea } from "../Suppliers/Suppliers.styles";
const AdminComponent = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <>
      <SidebarMenu
        sidebarOpen={sidebarOpen}
        setSidebarOpen={setSidebarOpen}
      />

      <Main $sidebarOpen={sidebarOpen}>
        <TopbarMenu
          sidebarOpen={sidebarOpen}
          setSidebarOpen={setSidebarOpen}
        />

        <ContentArea>
          <Outlet />
        </ContentArea>
      </Main>
    </>
  );
};

export default AdminComponent
