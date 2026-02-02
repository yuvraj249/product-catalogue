import styled from 'styled-components';

export const PageHeader = styled.div`
  display: grid;
  grid-template-columns: 0.9fr minmax(350px, 500px) 1.1fr;
  align-items: center;
  margin-bottom: 32px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    row-gap: 16px;
  }
`;

export const Title = styled.h1`
  font-size: 32px;
  color: #1f2d36;
  margin: 0;
  font-weight: 700;
`;

export const AddButton = styled.button`
  display: flex;
  align-items: center;
  margin-left: auto ;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #2A7B9B, #57C785);
  border: none;
  color: white;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    box-shadow: 0 6px 16px rgba(42,123,155,0.25);
    transform: translateY(-1px);
  }
`;

export const SearchWrap = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  justify-self: center;
  margin: 0 auto ;
  gap: 8px;
  padding: 10px 16px;
  background: #ffffff;
  border: 1px solid rgba(42,123,155,0.25);
  border-radius: 8px;
  max-width: 500px;
  width: 100%;
  transition: all 0.15s ease;

  &:focus-within {
    border-color: #2A7B9B;
    box-shadow: 0 0 0 3px rgba(42,123,155,0.15);
  }
`;

export const SearchIpt = styled.input`
  border: none;
  flex: 1;
  background: transparent;
  outline: none;
  color: #1f2d36;
  font-size: 14px;

  &::placeholder {
    color: #4f6b72;
  }
`;

export const TableWrapper = styled.div`
  background: #ffffff;
  border: 1px solid rgba(42,123,155,0.25);
  border-radius: 12px;
  overflow: hidden;
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
`;

export const Td = styled.td`
  padding: 16px;
  color: #1f2d36;
  border-bottom: 1px solid rgba(42,123,155,0.18);
  text-align: ${({ align }) => align || 'left'};
  font-size: 14px;
`;

export const ActionButtons = styled.div`
  display: flex;
  gap: 8px;
  justify-content: center;
`;

export const IconButton = styled.button`
  padding: 8px;
  border-radius: 4px;
  background: rgba(42,123,155,0.12);
  border: none;
  color: #2A7B9B;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: all 0.15s ease;

  &:hover {
    opacity: 0.9;
  }
`;

export const Modal = styled.div`
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  justify-content: center;
  padding: 24px;
  z-index: 3000;
  overscroll-behavior: contain;
`;

export const ModalContent = styled.div`
  background: #ffffff;
  width: 100%;
  max-width: 500px;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  max-height: 90vh;
  overflow-y: visible;
`;

export const ModalHeader = styled.div`
  padding: 20px 24px;
  border-bottom: 1px solid rgba(42,123,155,0.25);
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const ModalTitle = styled.h2`
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #1f2d36;
`;

export const CloseButton = styled.button`
  background: none;
  border: none;
  color: #2A7B9B;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;

  &:hover {
    background: rgba(42,123,155,0.12);
  }
`;

export const Form = styled.form`
  padding: 24px;
`;

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
`;

export const Label = styled.label`
  font-size: 14px;
  font-weight: 600;
  color: #1f2d36;
`;

export const Input = styled.input`
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid rgba(42,123,155,0.30);
  background: #f3f8fb;
  font-size: 14px;
  outline: none;
  color: #1f2d36;
  transition: all 0.15s ease;

  &:focus {
    border-color: #2A7B9B;
    box-shadow: 0 0 0 3px rgba(42,123,155,0.15);
    background: white;
  }

  &::placeholder {
    color: #4f6b72;
  }
`;

export const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 8px;
`;

export const CancelButton = styled.button`
  padding: 10px 24px;
  background: transparent;
  color: #4f6b72;
  border: 1px solid rgba(42,123,155,0.30);
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: rgba(42,123,155,0.10);
  }
`;

export const SubmitButton = styled.button`
  padding: 10px 24px;
  background: linear-gradient(135deg, #2A7B9B, #57C785);
  color: white;
  border-radius: 8px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    box-shadow: 0 6px 16px rgba(42,123,155,0.25);
  }
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;


export const ToastContainer = styled.div`
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 12px;
`;


export const Toast = styled.div`
  background: #ffffff;
  padding: 12px 18px;
  border-radius: 2px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.15);
  font-weight: 600;
  min-width: 260px;
  min-height: 28px;
  color: #1f2d36;
  border-left: 6px solid
    ${({ $type }) => ($type === 'success' ? '#2ecc71' : '#e74c3c')};

  animation: slideIn 0.3s ease forwards;

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateX(60px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
`;