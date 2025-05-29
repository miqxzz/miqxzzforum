import React from 'react';
import styled from 'styled-components';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../Chat/AuthContext'; // Импортируем useAuth

const Nav = styled.nav`
  height: 100vh;
  background: #a678f7;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding-top: 60px;
  font-family: 'Montserrat', sans-serif;
`;

const NavTitle = styled.h1`
  font-size: 2.2em;
  color: #fff;
  font-weight: 700;
  margin-bottom: 60px;
  text-align: center;
  font-family: 'Montserrat', sans-serif;
`;

const NavLinks = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  a {
    color: #fff;
    margin: 18px 0;
    font-size: 1.2em;
    text-decoration: none;
    font-family: 'Montserrat', sans-serif;
    font-weight: 400;
    transition: color 0.2s;
    &:hover {
      color: #e0c8ff;
    }
  }
`;

const Navbar = () => {
    const navigate = useNavigate();
    const { isAuthenticated, logout } = useAuth(); // Используем useAuth для получения состояния аутентификации и метода logout

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <Nav>
            <NavTitle>Форум</NavTitle>
            <NavLinks>
                {isAuthenticated ? (
                    <>
                        <Link to="/posts">Посты</Link>
                        <Link to="#" onClick={handleLogout}>Выход</Link>
                    </>
                ) : (
                    <>
                        <Link to="/login">Вход</Link>
                        <Link to="/register">Регистрация</Link>
                    </>
                )}
            </NavLinks>
        </Nav>
    );
};

export default Navbar;