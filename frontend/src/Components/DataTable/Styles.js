import styled from "styled-components";

export const TableWrapper = styled.div`
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  max-height: calc(100vh - 220px);
  overflow-y: auto;
  overflow-x: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.25);
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;

  tbody tr {
    transition: all 0.2s ease;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: rgba(56, 189, 248, 0.05);
      box-shadow: inset 3px 0 0 #38bdf8;
    }
  }
`;

export const Th = styled.th`
  padding: 14px 18px;
  background: rgba(30, 41, 59, 0.85);
  backdrop-filter: blur(8px);
  color: #94a3b8;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  text-align: ${({ align }) => align || 'left'};
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 12.5px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  position: sticky;
  top: 0;
  z-index: 10;
`;

export const Td = styled.td`
  padding: 15px 18px;
  color: #f1f5f9;
  text-align: ${({ align }) => align || 'left'};
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 14px;
  font-weight: 400;
`;

export const SortWrapper = styled.span`
  display: flex;
  flex-direction: column;
  margin-left: 6px;
  gap: 2px;
`;

export const Arrow = styled.span`
  color: ${({ $active }) => ($active ? '#38bdf8' : '#64748b')};
  opacity: ${({ $active }) => ($active ? 1 : 0.4)};
  font-size: 10px;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    color: #38bdf8;
    opacity: 1;
  }
`;

export const HeaderContent = styled.div`
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: ${({ $center }) => $center ? "center" : "flex-start"};
  gap: 8px;
`;
