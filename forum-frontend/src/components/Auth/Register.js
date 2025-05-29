import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import Chat from '../Chat/chat';

const RegisterPageWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
`;

const RegisterContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 8px #e1d5ee22;
  padding: 40px 48px 32px 48px;
  margin-bottom: 32px;
`;

const RegisterTitle = styled.h2`
  color: #a678f7;
  font-size: 2em;
  font-weight: 700;
  margin-bottom: 24px;
  text-align: center;
  font-family: 'Montserrat', sans-serif;
`;

const Input = styled.input`
  padding: 12px;
  margin: 10px 0;
  border: 2px solid #a678f7;
  border-radius: 8px;
  width: 320px;
  box-sizing: border-box;
  font-size: 1.1em;
  font-family: 'Montserrat', sans-serif;
  color: #3d2466;
  background: #fff;
  transition: border 0.2s;
  &:focus {
    border: 2px solid #8e44ad;
    outline: none;
  }
`;

const RegisterButton = styled.button`
  background: #a678f7;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 12px 0;
  width: 100%;
  font-family: 'Montserrat', sans-serif;
  font-weight: 600;
  font-size: 1.1em;
  margin-top: 10px;
  margin-bottom: 10px;
  transition: background 0.2s;
  box-shadow: 0 2px 8px #e1d5ee44;
  &:hover {
    background: #8e44ad;
  }
`;

const ErrorMessage = styled.div`
  color: #dc3545;
  margin: 10px 0;
  text-align: center;
`;

const PasswordRequirements = styled.div`
  margin: 10px 0;
  padding: 10px;
  background-color: #f8f9fa;
  border-radius: 5px;
  width: 300px;
  box-sizing: border-box;
`;

const ChatWrapper = styled.div`
  display: flex;
  justify-content: center;
  width: 100%;
`;

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleRegister = async () => {
        try {
            setError('');
            const response = await axios.post('http://localhost:8080/auth/register', 
              { 
                username: username,
                password: password,
                role: "user",
              });
            localStorage.setItem('token', response.data.token);
            navigate('/Login');
        } catch (error) {
            console.error('Не удалось зарегистрироваться:', error);
            setError(error.response?.data?.error || 'Не удалось зарегистрироваться. Пожалуйста, попробуйте снова.');
        }
    };

    return (
        <RegisterPageWrapper>
            <RegisterContainer>
                <RegisterTitle>Регистрация</RegisterTitle>
                <Input
                    type="text"
                    placeholder="Логин"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <Input
                    type="password"
                    placeholder="Пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                {error && <ErrorMessage>{error}</ErrorMessage>}
                <RegisterButton onClick={handleRegister}>Зарегистрироваться</RegisterButton>
            </RegisterContainer>
            <ChatWrapper>
                <Chat />
            </ChatWrapper>
        </RegisterPageWrapper>
    );
};

export default Register;