import { useState } from "react";
import { Outlet } from "react-router-dom";
import SidebarMenu from "../Sidebar";
import TopbarMenu from "../Topbar";
import { Main, ContentArea } from "../Sidebar/Styles";
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
