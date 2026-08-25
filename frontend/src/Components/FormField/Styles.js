import styled from "styled-components";

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 18px;
`;

export const Label = styled.label`
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 13.5px;
  font-weight: 600;
  color: #e2e8f0;

  span {
    color: #ef4444;
    margin-left: 2px;
  }
`;