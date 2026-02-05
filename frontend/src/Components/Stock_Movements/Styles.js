import styled, { css } from "styled-components";
import Select from "react-select";

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

export const Button = styled.button`
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  ${({ size = "md" }) => {
    switch (size) {
      case "header":
        return css`
          padding: 12px 24px;
          font-size: 16px;
          border-radius: 8px;
        `;
      case "sm":
        return css`
          padding: 8px 18px;
          font-size: 14px;
          border-radius: 8px;
        `;
      default:
        return css`
          padding: 12px 26px;
          font-size: 15px;
          border-radius: 10px;
        `;
    }
  }}

  ${({ variant = "primary" }) => {
    switch (variant) {
      case "secondary":
        return css`
          background: transparent;
          color: #4f6b72;
          border: 1px solid rgba(42, 123, 155, 0.3);
        `;
      case "danger":
        return css`
          background: linear-gradient(135deg, #c0392b, #e74c3c);
          color: white;
        `;
      default:
        return css`
          background: linear-gradient(135deg, #2a7b9b, #57c785);
          color: white;
        `;
    }
  }}
`;


export const SearchWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 8px;
  max-width: 500px;
  transition: all 0.15s ease;

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
  font-size: 14px;
  color: #1f2d36;

  &::placeholder {
    color: #4f6b72;
  }
`;

export const ModalHeader = styled.div`
  padding: 22px 28px;
  border-bottom: 1px solid rgba(42, 123, 155, 0.25);
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const ModalTitle = styled.h2`
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: #1f2d36;
`;

export const CloseButton = styled.button`
  background: none;
  border: none;
  color: #2a7b9b;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  display: flex;

  &:hover {
    background: rgba(42, 123, 155, 0.12);
  }
`;

export const Form = styled.form`
  padding: 24px;
  max-height: 70vh;
  overflow-y: auto;
`;

export const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 8px;
`;

export const Input = styled.input`
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid rgba(42, 123, 155, 0.3);
  background: #f3f8fb;
  font-size: 16px;
  outline: none;
  transition: all 0.15s ease;

  &:focus {
    border-color: #2a7b9b;
    box-shadow: 0 0 0 3px rgba(42, 123, 155, 0.15);
    background: white;
  }
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;

export const SelectBox = styled(Select)`
  .react-select__control {
    min-height: 52px;
    border-radius: 12px;
    border: 1px solid rgba(42, 123, 155, 0.3);
    background: #f3f8fb;
    padding: 6px 8px;
    box-shadow: none;
  }

  .react-select__control--is-focused {
    border-color: #2a7b9b;
    box-shadow: 0 0 0 3px rgba(42, 123, 155, 0.15);
    background: white;
  }

  .react-select__option--is-selected {
    background: #2a7b9b;
    color: white;
  }

  .react-select__option--is-focused {
    background: rgba(42, 123, 155, 0.12);
  }
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

 export const HeaderAction = styled.div`
  margin-left: auto;
  display: flex;
  align-items: center;
`;

export const HeaderAddButton = styled(Button)`
  padding: 10px 22px;
  font-size: 14px;
  border-radius: 8px;
  white-space: nowrap;
`;
