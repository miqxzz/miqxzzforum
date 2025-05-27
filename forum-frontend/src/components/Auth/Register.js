import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';

const RegisterContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 0;
  padding-top: 30px;
  width: 100%;
`;

const Input = styled.input`
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #a259ff;
  border-radius: 5px;
  width: 100%;
  max-width: 300px;
  box-sizing: border-box;
  font-family: 'Montserrat', Arial, sans-serif;
  &::placeholder {
    font-family: 'Montserrat', Arial, sans-serif;
    color: #bfa6e6;
    opacity: 1;
  }
`;

const Button = styled.button`
  padding: 10px 15px;
  background-color: #a259ff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  font-family: 'Montserrat', Arial, sans-serif;
  &:hover {
    background-color: #6c2eb7;
  }
`;

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleRegister = async () => {
        try {
            const response = await axios.post('http://localhost:8080/register', 
              { 
                username: username,
                 password: password,
                 role: "user",
                 });
            localStorage.setItem('token', response.data.token);
            navigate('/Login');
        } catch (error) {
            console.error('Ошибка при регистрации:', error);
            alert('Ошибка при регистрации. Пожалуйста, попробуйте снова.');
        }
    };

    return (
        <RegisterContainer>
            <h2 style={{ marginBottom: '20px', color: '#a259ff', fontFamily: 'Montserrat, Arial, sans-serif' }}>Регистрация</h2>
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
            <Button onClick={handleRegister}>Зарегистрироваться</Button>
        </RegisterContainer>
    );
};

export default Register;