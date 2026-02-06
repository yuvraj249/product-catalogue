import styled from 'styled-components';
import ReactSelect from 'react-select'

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

export const CategoryDescription = styled.div`
  color: #4f6b72;
  line-height: 1.5;
  max-width: 400px;
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

  ${({ $danger }) =>
    $danger &&
    `
    background: rgba(199,88,88,0.12);
    color: #c75858;
  `}

  &:hover {
    opacity: 0.9;
  }
`;


export const ModalHeader = styled.div`
  padding: 22px 28px;
  border-bottom: 1px solid rgba(42,123,155,0.25);
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
  max-height: 70vh;
  overflow-y: auto;
`;

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
`;

export const Label = styled.label`
  font-size: 16px;
  font-weight: 600;
  color: #1f2d36;
`;

export const Input = styled.input`
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid rgba(42,123,155,0.30);
  background: #f3f8fb;
  font-size: 16px;
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

export const Textarea = styled.textarea`
  padding: 14px 17px;
  border-radius: 12px;
  border: 1px solid rgba(42,123,155,0.30);
  background: #f3f8fb;
  font-size: 15px;
  outline: none;
  color: #1f2d36;
  font-family: inherit;
  resize: vertical;
  transition: all 0.15s ease;
  overflow-y: auto;

  &::-webkit-scrollbar {
    width: 0;
  }

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
  padding: 12px 26px;
  background: transparent;
  color: #4f6b72;
  border: 1px solid rgba(42,123,155,0.30);
  border-radius: 10px;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: rgba(42,123,155,0.10);
  }
`;

export const SubmitButton = styled.button`
  padding: 12px 26px;
  background: linear-gradient(135deg, #2A7B9B, #57C785);
  color: white;
  border-radius: 10px;
  border: none;
  font-weight: 600;
  font-size: 15px;
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


export const CategorySelect = styled(ReactSelect)`
  .react-select__control {
    border-radius: 8px;
    background: #f3f8fb;
    border: 1px solid rgba(42,123,155,0.30);
    min-height: 38px;
  }

  .react-select__control--is-focused {
    border-color: #2A7B9B;
    box-shadow: 0 0 0 3px rgba(42,123,155,0.15);
    background: white;
  }

  .react-select__menu-list {
    max-height: 180px;
    overflow-y: auto;
  }

  .react-select__option {
    padding: 6px 10px;
    font-size: 14px;
  }

  .react-select__option--is-selected {
    background: rgba(42,123,155,0.12);
    color: #2A7B9B;
  }

  .react-select__option--is-focused {
    background: rgba(42,123,155,0.08);
  }

  .react-select__value-container {
    padding: 4px 8px;
  }

  .react-select__indicator-separator {
    display: none;
  }
`;



