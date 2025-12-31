import styled from "styled-components";

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
  margin-left: auto;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #2a7b9b, #57c785);
  border: none;
  color: white;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;

  &:hover {
    box-shadow: 0 6px 16px rgba(42, 123, 155, 0.25);
    transform: translateY(-1px);
  }
`;

export const SearchWrap = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  justify-self: center;
  margin: 0 auto;
  padding: 10px 16px;
  background: #ffffff;
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 8px;
  max-width: 500px;
  width: 100%;

  &:focus-within {
    border-color: #2a7b9b;
    box-shadow: 0 0 0 3px rgba(42, 123, 155, 0.15);
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
  background: white;
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  overflow: hidden;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

export const Th = styled.th`
  padding: 16px;
  background: #f3f8fb;
  color: #4f6b72;
  border-bottom: 2px solid rgba(42, 123, 155, 0.25);
  text-align: ${({ align }) => align || "left"};
  font-size: 14px;
  font-weight: 600;
`;

export const Td = styled.td`
  padding: 16px;
  color: #1f2d36;
  border-bottom: 1px solid rgba(42, 123, 155, 0.18);
  text-align: ${({ align }) => align || "left"};
  font-size: 14px;
`;

export const ActionButtons = styled.div`
  display: flex;
  gap: 8px;
  justify-content: center;
`;

export const IconButton = styled.button`
  padding: 8px;
  border-radius: 6px;
  background: rgba(42, 123, 155, 0.12);
  border: none;
  cursor: pointer;

  ${({ $danger }) =>
    $danger &&
    `
    background: rgba(199,88,88,0.12);
    color: #c75858;
  `}
`;

export const Modal = styled.div`
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  z-index: 2000;
`;

export const ModalContent = styled.div`
  background: white;
  width: 100%;
  max-width: 600px;
  border-radius: 14px;
  overflow: hidden;
`;

export const ModalHeader = styled.div`
  padding: 20px 24px;
  border-bottom: 1px solid rgba(42, 123, 155, 0.25);
  display: flex;
  justify-content: space-between;
`;

export const ModalTitle = styled.h2`
  margin: 0;
`;

export const CloseButton = styled.button`
  border: none;
  background: transparent;
`;

export const Form = styled.form`
  padding: 24px;
`;

export const FormGroup = styled.div`
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

export const Label = styled.label`
  font-size: 14px;
  font-weight: 600;
`;

export const Input = styled.input`
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid rgba(42, 123, 155, 0.3);
`;

export const Select = styled.select`
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid rgba(42, 123, 155, 0.3);
`;

export const Textarea = styled.textarea`
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid rgba(42, 123, 155, 0.3);
`;

export const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
`;

export const CancelButton = styled.button`
  padding: 10px 24px;
`;

export const SubmitButton = styled.button`
  padding: 10px 24px;
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;

export const ToastContainer = styled.div`
  position: fixed;
  top: 20px;
  right: 20px;
`;

export const Toast = styled.div`
  padding: 10px 16px;
  background: white;
  border-left: 6px solid
    ${({ $type }) => ($type === "success" ? "#2ecc71" : "#e74c3c")};
`;
