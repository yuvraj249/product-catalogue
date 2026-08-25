import styled, { keyframes } from "styled-components";

const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
`;

export const Modal = styled.div`
  position: fixed;
  inset: 0;
  background: rgba(10, 13, 18, 0.75);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  z-index: 2000;
`;

export const ModalContent = styled.div`
  background: rgba(15, 23, 42, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.12);
  width: 100%;
  max-width: 540px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.5), 0 0 40px rgba(56, 189, 248, 0.1);
  animation: ${fadeIn} 0.25s cubic-bezier(0.4, 0, 0.2, 1) forwards;
  color: #f1f5f9;
`;
