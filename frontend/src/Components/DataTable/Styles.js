import styled from "styled-components";

export const TableWrapper = styled.div`
  background: #ffffff;
  border: 1px solid hsla(197, 57%, 39%, 0.25);
  border-radius: 12px;
  max-height: calc(100vh - 220px);
  overflow-y: auto;
  overflow-x: hidden;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;

  tbody tr {
    transition: background 0.15s ease;

    &:hover {
      background: rgba(42,123,155,0.05);
    }
  }
`;

export const Th = styled.th`
  padding: 12px 16px;
  background: #f3f8fb;
  color: #4f6b72;
  border-bottom: 2px solid rgba(42,123,155,0.25);
  text-align: ${({ align }) => align || 'left'};
  font-size: 14px;
  font-weight: 600;
  position: sticky;
  top: 0;
  align-items: center;
  gap: 6px;

`;

export const Td = styled.td`
  padding: 16px;
  color: #1f2d36;
  border-bottom: 1px solid rgba(42,123,155,0.18);
  text-align: ${({ align }) => align || 'left'};
  font-size: 14px;
`;

export const SortWrapper = styled.span`
  display: flex;
  flex-direction: column;
  margin-left: 4px;
  font-size: 9px;
  line-height: 9px;
`;

export const Arrow = styled.span`
  opacity: ${({ $active }) => ($active ? 1 : 0.5)};
  font-weight: ${({ $active }) => ($active ? 900 : 400)};
  font-size: 11px;
  line-height: 12px;
  cursor: pointer;

`;

export const HeaderContent = styled.div`
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: ${({ $center }) => $center ? "center" : "flex-start"};
  gap: 10px;
`;

