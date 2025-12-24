
import React, { useState } from 'react';
import { Page, Card, LogoBox, Title, Subtitle,Input, PasswordWrapper, IconButton, SubmitButton , Form, FormGroup, Header} from './LoginPage.Styles';
 
import eye from '../Login/assets/eye-open.svg'
import eyeOff from '../Login/assets/eye-closed.svg'
import box from '../Login/assets/logo.svg'
import arrow from '../Login/assets/arrow.svg'

export const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
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
        <Form>
         <FormGroup>
          <label>Email Address</label>
        <Input placeholder='abc@gmail.com' value={email} onChange={e => setEmail(e.target.value)}/>
         </FormGroup>

         <FormGroup>
         <label>Password</label>
        <PasswordWrapper>
          <Input placeholder='password' type={showPassword ? "text" : "password"} value={password} onChange={e => setPassword(e.target.value)} />
          <IconButton onClick={() => setShowPassword(!showPassword)}>
            <img src={showPassword ? eyeOff : eye} alt='toggle password' width={18} />
          </IconButton>
        </PasswordWrapper>
        </FormGroup> 
       
        <SubmitButton>
          Sign In
          <img src={arrow} alt="" width={18} style={{ marginLeft: 8 }} />
        </SubmitButton>
      
        </Form>
      </Card>
    </Page>
  );
}


