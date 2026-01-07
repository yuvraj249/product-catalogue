import styled from "styled-components";

export const Page = styled.div`
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;

  background:
    radial-gradient(circle at 30% 30%, rgba(42, 123, 155, 0.18), transparent 45%),
    radial-gradient(circle at 70% 70%, rgba(87, 199, 133, 0.18), transparent 45%),
    linear-gradient(135deg, #f3f8fb, #eef6f2);
`
export const Card = styled.div`
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  padding: 43px;
  border-radius: 16px;

  box-shadow:
    0 10px 30px rgba(0,0,0,0.12),
    0 0 0 1px rgba(42, 123, 155, 0.05);

  display: flex;
  flex-direction: column;
  gap: 20px;
`

export const Header = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
`

export const LogoBox = styled.div`
  width: 64px;
  height: 64px;
  margin: 0 auto;
  border-radius: 12px;
  background: linear-gradient(135deg ,#2A7B9B, #57C785);
  display: flex;
  align-items: center;
  justify-content: center;
`

export const Title = styled.h1`
  text-align: center;
  font-size: 24px;
  font-weight: 700;
  margin: 0;
`

export const Subtitle = styled.p`
  text-align: center;
  font-size: 16px;
  color: #6b6b6b;
  margin: 0;
`


export const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 16px;
`

export const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`

export const Label = styled.label`
  font-size: 14px;
  font-weight: 500;
  color: #2c2c2c;
`

export const Input = styled.input`
  width: 100%;
  box-sizing: border-box;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid #d6e4ec;
  font-size: 16px;

  &:focus {
    outline: none;
    border-color: #2A7B9B;
    box-shadow: 0 0 0 3px rgba(42, 123, 155, 0.18);
  }
`


export const PasswordWrapper = styled.div`
  position: relative;

  input {
    padding-right: 44px;
  }
`

export const IconButton = styled.button`
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
`

export const SubmitButton = styled.button`
  width: 100%;
  padding: 14px;
  border-radius: 8px;
  border: none;

  background: linear-gradient(135deg,#2A7B9B, #57C785);
  color: white;
  font-size: 16px;
  font-weight: 600;

  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(0,0,0,0.15);
  }
`
