import styled from "styled-components";

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
