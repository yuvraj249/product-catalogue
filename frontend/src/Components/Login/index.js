
import React, { useState } from 'react';
import {useNavigate} from 'react-router-dom';
import { Page, Card, LogoBox, Title, Subtitle,Input, PasswordWrapper, IconButton, SubmitButton , Form, FormGroup, Header} from './Styles';
 
import eye from '../../Images/eye-open.svg'
import eyeOff from '../../Images/eye-closed.svg'
import box from '../../Images/logo.svg'
import arrow from '../../Images/arrow.svg'
import api from '../../Api/axios'
import { setToken } from '../../utils/auth'
import { toast, ToastContainer } from 'react-toastify';

export const LoginPage = () => {
  const navigate = useNavigate()
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data } = await api.post("auth/login", {
        email,
        password,
      });
      console.log("LOGIN RESPONSE:", data);
      setToken(data.token);
      toast.success("Login successful")  
      navigate("/admin/categories");
    } catch (err) {
      console.error("LOGIN ERROR:", err);
      const msg = err.response.data.error || "";
      if (msg === "Invalid credentials") {
        toast.error("Incorrect email or password");
      } else {
        toast.error("Login failed");
      }
    } finally {
      setLoading(false);
    }
  };
  return (
    <Page>
      <Card>
        <LogoBox>
          <img src={box} alt='logo' width={32} />
        </LogoBox>       
        <Header>
        <Title>Product Catalogue</Title>
        <Subtitle>Sign in to your account</Subtitle>
        </Header>
        <Form onSubmit={handleSubmit}>
         <FormGroup>
          <label>Email Address</label>
        <Input placeholder='abc@gmail.com' value={email} onChange={e => setEmail(e.target.value)}/>
         </FormGroup>

         <FormGroup>
         <label>Password</label>
        <PasswordWrapper>
          <Input placeholder='password' type={!showPassword ? "text" : "password"} value={password} onChange={e => setPassword(e.target.value)} />
          <IconButton type='button' onClick={() => setShowPassword(!showPassword)}>
            <img src={showPassword ? eyeOff : eye} alt='toggle password' width={18} />
          </IconButton>
        </PasswordWrapper>
        </FormGroup> 
       
       <SubmitButton type="submit" disabled={loading}>
            {loading ? "Signing In..." : "Sign In"}
            <img src={arrow} alt="" width={18} style={{ marginLeft: 8 }} />
        </SubmitButton>
      
        </Form>
      </Card>
      <ToastContainer position="top-right" autoClose={3000} />
    </Page>
  );
}


