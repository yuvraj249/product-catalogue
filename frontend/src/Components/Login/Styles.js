import styled, { keyframes } from "styled-components";

const floatAnim = keyframes`
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-8px) rotate(2deg); }
`;

export const Page = styled.div`
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  overflow: hidden;

  background:
    radial-gradient(circle at 15% 15%, rgba(56, 189, 248, 0.12), transparent 40%),
    radial-gradient(circle at 85% 85%, rgba(129, 140, 248, 0.12), transparent 40%),
    #0a0d12;
`;

export const Card = styled.div`
  width: 100%;
  max-width: 420px;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  padding: 40px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);

  box-shadow:
    0 24px 60px rgba(0, 0, 0, 0.5),
    0 0 40px rgba(56, 189, 248, 0.1);

  display: flex;
  flex-direction: column;
  gap: 24px;
  animation: ${floatAnim} 6s ease-in-out infinite;
`;

export const Header = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
`;

export const LogoBox = styled.div`
  width: 64px;
  height: 64px;
  margin: 0 auto 12px;
  border-radius: 18px;
  background: linear-gradient(135deg, #38bdf8 0%, #3b82f6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(56, 189, 248, 0.35);

  img {
    filter: brightness(0) invert(1);
  }
`;

export const Title = styled.h1`
  text-align: center;
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 26px;
  font-weight: 800;
  color: #f8fafc;
  margin: 0;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #9ca3af 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
`;

export const Subtitle = styled.p`
  text-align: center;
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 14.5px;
  color: #94a3b8;
  margin: 0;
`;

export const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-family: var(--font-body, 'Inter', sans-serif);
    font-size: 13.5px;
    font-weight: 600;
    color: #e2e8f0;
  }
`;

export const Label = styled.label`
  font-size: 13.5px;
  font-weight: 600;
  color: #e2e8f0;
`;

export const Input = styled.input`
  width: 100%;
  box-sizing: border-box;
  padding: 13px 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(15, 23, 42, 0.8);
  color: #f1f5f9;
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 14px;
  transition: all 0.2s ease;

  &::placeholder {
    color: #64748b;
  }

  &:focus {
    outline: none;
    border-color: #38bdf8;
    box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
    background: rgba(15, 23, 42, 0.95);
  }
`;

export const PasswordWrapper = styled.div`
  position: relative;

  input {
    padding-right: 44px;
  }
`;

export const IconButton = styled.button`
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  opacity: 0.6;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 1;
  }

  img {
    filter: invert(80%);
  }
`;

export const SubmitButton = styled.button`
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  border: none;
  margin-top: 6px;

  background: linear-gradient(135deg, #38bdf8 0%, #2563eb 100%);
  color: white;
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 15px;
  font-weight: 600;

  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 20px rgba(56, 189, 248, 0.3);

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(56, 189, 248, 0.45);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
`;
