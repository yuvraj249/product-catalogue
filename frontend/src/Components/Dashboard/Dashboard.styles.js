import styled from 'styled-components'

export const Layout = styled.div`
  display: flex;
  min-height: 100vh;

  background:
    radial-gradient(circle at 20% 20%, rgba(42, 123, 155, 0.10), transparent 40%),
    radial-gradient(circle at 80% 80%, rgba(87, 199, 133, 0.10), transparent 40%),
    #f3f8fb;

  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
`;

export const Sidebar = styled.aside`
  width: ${({ $isOpen }) => $isOpen ? '260px' : '0'};
  background: #ffffff;
  border-right: 1px solid rgba(42, 123, 155, 0.20);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 1000;
  transition: width 0.25s ease;
  overflow: hidden;
`;

export const SidebarHeader = styled.div`
  padding: 24px;
  border-bottom: 1px solid rgba(42, 123, 155, 0.20);
`;

export const SidebarTitle = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: #2A7B9B;
`;

export const Nav = styled.nav`
  flex: 1;
  padding: 16px;
  overflow-y: auto;
`;

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

export const Main = styled.main`
  flex: 1;
  margin-left: ${({ $sidebarOpen }) => $sidebarOpen ? '260px' : '0'};
  transition: margin-left 0.25s ease;

  @media (max-width: 768px) {
    margin-left: 0;
  }
`;

export const TopBar = styled.div`
  height: 64px;
  background: #ffffff;
  border-bottom: 1px solid rgba(42, 123, 155, 0.2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
`;

export const MenuButton = styled.button`
  background: none;
  border: none;
  color: #2A7B9B;
  padding: 8px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.15s ease;

  &:hover {
    background: rgba(42, 123, 155, 0.10);
  }
`;

export const UserRole = styled.div`
  font-size: 14px;
  font-weight: 500;
  background: rgba(42, 123, 155, 0.15);
  padding: 6px 12px;
  border-radius: 8px;
  color: #2A7B9B;
`;

export const ContentArea = styled.div`
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;

  @media (max-width: 768px) {
    padding: 24px;
  }
`;

export const PageHeader = styled.div`
  margin-bottom: 32px;
`;

export const Title = styled.h1`
  font-size: 32px;
  font-weight: 700;
  color: #1f2d36;
  margin: 0;
`;

export const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
`;

export const StatCard = styled.div`
  background: linear-gradient(135deg, #ffffff, #f3f8fb);
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  padding: 24px;

  display: flex;
  align-items: center;
  gap: 16px;

  transition: all 0.25s ease;
  cursor: pointer;

  &:hover {
    border-color: rgba(87, 199, 133, 0.6);
    box-shadow: 0 10px 25px rgba(42, 123, 155, 0.18);
    transform: translateY(-2px);
  }
`;

export const StatIcon = styled.div`
  width: 56px;
  height: 56px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: linear-gradient(
    135deg,
    rgba(42, 123, 155, 0.18),
    rgba(87, 199, 133, 0.18)
  );

  color: #2A7B9B;
  border-radius: 8px;
`;

export const StatInfo = styled.div`
  flex: 1;
`;

export const StatValue = styled.div`
  font-size: 24px;
  font-weight: 700;
  color: #1f2d36;
  margin-bottom: 4px;
`;

export const StatLabel = styled.div`
  font-size: 14px;
  color: #4f6b72;
`;

export const Section = styled.div`
  background: linear-gradient(135deg, #ffffff, #f3f8fb);
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  padding: 32px;

  @media (max-width: 768px) {
    padding: 20px;
  }
`;

export const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
`;

export const SectionTitle = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #1f2d36;
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
`;

export const ThresholdBadge = styled.div`
  padding: 6px 12px;
  background: rgba(87, 199, 133, 0.15);
  border: 1px solid rgba(87, 199, 133, 0.3);
  border-radius: 4px;
  color: #2A7B9B;
  font-size: 12px;
  font-weight: 500;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

export const Th = styled.th`
  text-align: ${({ $align }) => $align || 'left'};
  padding: 12px 16px;
  background: #f3f8fb;
  color: #4f6b72;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 2px solid rgba(42, 123, 155, 0.25);
`;

export const Tr = styled.tr`
  transition: background 0.15s ease;

  &:hover {
    background: rgba(42, 123, 155, 0.05);
  }
`;

export const Td = styled.td`
  text-align: ${({ $align }) => $align || 'left'};
  padding: 16px;
  border-bottom: 1px solid rgba(42, 123, 155, 0.18);
  color: #1f2d36;
  font-size: 14px;
  font-weight: ${({ $align }) => $align === 'right' ? 400 : 500};
`;

export const StockBadge = styled.span`
  display: inline-block;
  padding: 4px 12px;
  background: rgba(199, 88, 88, 0.12);
  color: #c75858;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;

export const TableWrapper = styled.div`
  background: #ffffff;
  border: 1px solid rgba(42,123,155,0.25);
  border-radius: 12px;
  overflow: hidden;
  margin-top: 24px;
`;

