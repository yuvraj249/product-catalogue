import styled from "styled-components";

export const Sidebar = styled.aside`
  width: ${({ $isOpen }) => ($isOpen ? '260px' : '0')};
  background: #ffffff;
  border-right: 1px solid rgba(42,123,155,0.20);
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  overflow: hidden;
  transition: width 0.25s ease;
  z-index: 800;
`;

export const SidebarHeader = styled.div`
  padding: 24px;
  font-weight: 700;
  font-size: 18px;
  color: #2A7B9B;
  border-bottom: 1px solid rgba(42,123,155,0.20);
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const Nav = styled.nav`
  flex: 1;
  padding: 16px;
  overflow-y: auto;
`

export const NavItem = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;

  background: ${({ $active }) =>
    $active
      ? 'rgba(42, 123, 155, 0.12)'
      : 'transparent'};

  border: none;
  border-radius: 8px;

  color: ${({ $active }) =>
    $active
      ? '#2A7B9B'
      : '#4a6b75'};

  font-size: 16px;
  font-weight: ${({ $active }) => $active ? 600 : 400};
  text-align: left;
  margin-bottom: 4px;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: rgba(42, 123, 155, 0.12);
    color: #2A7B9B;
  }
`;

export const LogoutWrapper = styled.div`
  padding: 16px;
  border-top: 1px solid rgba(42, 123, 155, 0.20);
`;

export const LogoutButton = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  border-radius: 8px;

  color: #c75858;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: rgba(199, 88, 88, 0.10);
  }
`;

export const TopBar = styled.div`
  height: 64px;
  background: #ffffff;
  border-bottom: 1px solid rgba(42,123,155,0.20);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  position: fixed;
  top: 0;
  left: ${({ $sidebarOpen }) => ($sidebarOpen ? '260px' : '0')};
  right: 0;
  z-index: 1200;
  transition: left .25s ease;
`;

export const MenuButton = styled.button`
  background: none;
  border: none;
  color: #2A7B9B;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;

  &:hover {
    background: rgba(42,123,155,0.12);
  }
`;

export const UserRole = styled.div`
  background: rgba(42,123,155,0.15);
  color: #2A7B9B;
  padding: 6px 12px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;
