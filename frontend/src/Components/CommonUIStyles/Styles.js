import styled, { css } from "styled-components";

export const PageHeader = styled.div`
  display: grid;
  grid-template-columns: 0.9fr minmax(300px, 500px) 1.1fr;
  align-items: center;
  margin-bottom: 28px;
  gap: 16px;

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
    row-gap: 16px;
  }
`;

export const Title = styled.h1`
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 28px;
  color: #f3f4f6;
  margin: 0;
  font-weight: 700;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #9ca3af 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
`;

export const HeaderAction = styled.div`
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;

  @media (max-width: 900px) {
    margin-left: 0;
  }
`;

export const Button = styled.button`
  font-family: var(--font-body, 'Inter', sans-serif);
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(56, 189, 248, 0.25);
  }

  &:active {
    transform: translateY(0);
  }

  ${({ size = "md" }) => {
    switch (size) {
      case "header":
        return css`
          padding: 10px 20px;
          font-size: 14px;
          border-radius: 10px;
        `;
      case "sm":
        return css`
          padding: 7px 14px;
          font-size: 13px;
          border-radius: 8px;
        `;
      default:
        return css`
          padding: 11px 22px;
          font-size: 14px;
          border-radius: 10px;
        `;
    }
  }}

  ${({ variant = "primary" }) => {
    switch (variant) {
      case "secondary":
        return css`
          background: rgba(30, 41, 59, 0.7);
          color: #94a3b8;
          border: 1px solid rgba(255, 255, 255, 0.1);
          backdrop-filter: blur(8px);

          &:hover {
            background: rgba(51, 65, 85, 0.9);
            color: #f8fafc;
            border-color: rgba(255, 255, 255, 0.2);
            box-shadow: 0 4px 14px rgba(0, 0, 0, 0.3);
          }
        `;
      case "danger":
        return css`
          background: linear-gradient(135deg, #e11d48, #be123c);
          color: white;

          &:hover {
            box-shadow: 0 6px 20px rgba(225, 29, 72, 0.35);
          }
        `;
      default:
        return css`
          background: linear-gradient(135deg, #0284c7 0%, #2563eb 100%);
          color: white;

          &:hover {
            box-shadow: 0 6px 20px rgba(37, 99, 235, 0.35);
          }
        `;
    }
  }}
`;

export const HeaderAddButton = styled(Button)`
  padding: 10px 20px;
  font-size: 14px;
  border-radius: 10px;
  white-space: nowrap;
  background: linear-gradient(135deg, #38bdf8 0%, #3b82f6 100%);
`;

export const ActionButtons = styled.div`
  display: flex;
  gap: 8px;
  justify-content: center;
`;

export const IconButton = styled.button`
  padding: 7px;
  border-radius: 8px;
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(56, 189, 248, 0.2);
    border-color: rgba(56, 189, 248, 0.4);
    transform: scale(1.05);
  }
`;

export const SearchWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  max-width: 500px;
  backdrop-filter: blur(12px);
  transition: all 0.25s ease;

  &:focus-within {
    border-color: #38bdf8;
    box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
    background: rgba(15, 23, 42, 0.85);
  }
`;

export const SearchIpt = styled.input`
  border: none;
  flex: 1;
  background: transparent;
  outline: none;
  font-size: 14px;
  color: #f1f5f9;

  &::placeholder {
    color: #64748b;
  }
`;

export const ModalHeader = styled.div`
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(15, 23, 42, 0.5);
`;

export const ModalTitle = styled.h2`
  margin: 0;
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 19px;
  font-weight: 700;
  color: #f8fafc;
`;

export const CloseButton = styled.button`
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #94a3b8;
  cursor: pointer;
  padding: 6px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(239, 68, 68, 0.15);
    color: #ef4444;
    border-color: rgba(239, 68, 68, 0.3);
  }
`;

export const Form = styled.form`
  padding: 24px;
  max-height: 75vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 12px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
`;

export const Input = styled.input`
  padding: 12px 16px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(15, 23, 42, 0.7);
  color: #f1f5f9;
  font-size: 14px;
  outline: none;
  transition: all 0.2s ease;

  &:focus {
    border-color: #38bdf8;
    box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
    background: rgba(15, 23, 42, 0.9);
  }

  &::placeholder {
    color: #64748b;
  }
`;

export const Icon = styled.img`
  width: 18px;
  height: 18px;
  opacity: 0.7;
`;
