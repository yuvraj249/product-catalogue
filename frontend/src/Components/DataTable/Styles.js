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

  tbody tr {
    transition: background 0.15s ease;

    &:hover {
      background: rgba(42,123,155,0.05);
    }
  }
`;

export const Th = styled.th`
  padding: 16px;
  background: #f3f8fb;
  color: #4f6b72;
  border-bottom: 2px solid rgba(42,123,155,0.25);
  text-align: ${({ align }) => align || 'left'};
  font-size: 14px;
  font-weight: 600;
  position: sticky;
  top: 0;
`;

export const Td = styled.td`
  padding: 16px;
  color: #1f2d36;
  border-bottom: 1px solid rgba(42,123,155,0.18);
  text-align: ${({ align }) => align || 'left'};
  font-size: 14px;
`;