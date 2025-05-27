import React from 'react';
import styled from 'styled-components';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../Chat/AuthContext';

const Sidebar = styled.nav`
  background: #a259ff;
  color: white;
  width: 220px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 30px 20px 0 20px;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 100;
`;

const SidebarTitle = styled.h1`
  font-size: 1.7em;
  margin-bottom: 40px;
  color: #fff;
  align-self: center;
`;

const SidebarLinks = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  a {
    color: #fff;
    margin: 10px 0;
    text-decoration: none;
    font-size: 1.1em;
    padding: 8px 12px;
    border-radius: 6px;
    transition: background 0.2s, color 0.2s;
    &:hover {
      background: #6c2eb7;
      color: #e0c3fc;
    }
  }
`;

const Navbar = () => {
    const navigate = useNavigate();
    const { isAuthenticated, logout } = useAuth();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <Sidebar>
            <SidebarTitle>Форум</SidebarTitle>
            <SidebarLinks>
                {isAuthenticated ? (
                    <>
                        <Link to="/posts">Посты</Link>
                        <Link to="#" onClick={handleLogout}>Выйти</Link>
                    </>
                ) : (
                    <>
                        <Link to="/login">Войти</Link>
                        <Link to="/register">Зарегистрироваться</Link>
                    </>
                )}
            </SidebarLinks>
        </Sidebar>
    );
};

export default Navbar;