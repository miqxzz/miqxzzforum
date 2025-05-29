import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../Chat/AuthContext'; // Импортируем useAuth

const LoginContainer = styled.div`
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

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();
    const { login } = useAuth(); // Используем useAuth для получения метода login

    const handleLogin = async () => {
        try {
            const response = await axios.post('http://localhost:8080/login', { username, password });
            login(response.data.token, {
              id: response.data.userID, // Проверьте что здесь число
              username: response.data.username,
              role: response.data.role,
              token: response.data.token
          }); // Вызываем метод login
            navigate('/posts');
        } catch (error) {
            console.error('Ошибка при входе:', error);
            alert('Ошибка при входе. Пожалуйста, проверьте свои учетные данные.');
        }
    };

    return (
        <LoginContainer>
            <h2 style={{ marginBottom: '20px', color: '#a259ff' }}>Войти</h2>
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
            <Button onClick={handleLogin}>Войти</Button>
        </LoginContainer>
    );
};

export default Login;