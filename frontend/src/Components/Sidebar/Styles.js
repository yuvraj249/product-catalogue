import styled from "styled-components";

export const Main = styled.main`
  min-height: 100vh;
  background:
    radial-gradient(circle at 10% 10%, rgba(56, 189, 248, 0.05), transparent 35%),
    radial-gradient(circle at 90% 90%, rgba(99, 102, 241, 0.05), transparent 35%),
    #0a0d12;

  margin-left: ${({ $sidebarOpen }) => ($sidebarOpen ? '260px' : '0')};
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  @media (max-width: 768px) {
    margin-left: 0;
  }
`;

export const ContentArea = styled.div`
  padding-top: 88px;
  padding-left: 32px;
  padding-right: 32px;
  padding-bottom: 40px;

  max-width: 1440px;
  margin: 0 auto;

  @media (max-width: 768px) {
    padding-left: 16px;
    padding-right: 16px;
    padding-top: 80px;
  }
`;

export const Sidebar = styled.aside`
  width: ${({ $isOpen }) => ($isOpen ? '260px' : '0')};
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  overflow: hidden;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 800;
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.3);
`;

export const SidebarHeader = styled.div`
  padding: 22px 24px;
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-weight: 700;
  font-size: 17px;
  color: #f8fafc;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(15, 23, 42, 0.4);
`;

export const Nav = styled.nav`
  flex: 1;
  padding: 20px 14px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const NavItem = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 11px 16px;

  background: ${({ $active }) =>
    $active
      ? 'linear-gradient(90deg, rgba(56, 189, 248, 0.15) 0%, rgba(56, 189, 248, 0.03) 100%)'
      : 'transparent'};

  border: none;
  border-left: ${({ $active }) =>
    $active
      ? '3px solid #38bdf8'
      : '3px solid transparent'};

  border-radius: ${({ $active }) =>
    $active
      ? '0 10px 10px 0'
      : '10px'};

  color: ${({ $active }) =>
    $active
      ? '#38bdf8'
      : '#94a3b8'};

  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 14.5px;
  font-weight: ${({ $active }) => ($active ? 600 : 500)};
  text-align: left;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    background: rgba(56, 189, 248, 0.08);
    color: #f1f5f9;
    transform: translateX(3px);
  }
`;

export const LogoutWrapper = styled.div`
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(15, 23, 42, 0.4);
`;

export const LogoutButton = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 11px 16px;
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.15);
  border-radius: 10px;

  color: #f87171;
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(239, 68, 68, 0.12);
    border-color: rgba(239, 68, 68, 0.3);
    color: #ef4444;
  }
`;

export const TopBar = styled.div`
  height: 68px;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 28px;
  position: fixed;
  top: 0;
  left: ${({ $sidebarOpen }) => ($sidebarOpen ? '260px' : '0')};
  right: 0;
  z-index: 700;
  transition: left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
`;

export const MenuButton = styled.button`
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #38bdf8;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(56, 189, 248, 0.15);
    border-color: rgba(56, 189, 248, 0.3);
  }
`;

export const UserRole = styled.div`
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.2);
  color: #38bdf8;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px 14px;
  height: 32px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 20px;
  letter-spacing: 0.03em;
  text-transform: uppercase;
`;

export const Icon = styled.img`
  width: 18px;
  height: 18px;
  filter: invert(70%) sepia(30%) saturate(1000%) hue-rotate(180deg);
`;
