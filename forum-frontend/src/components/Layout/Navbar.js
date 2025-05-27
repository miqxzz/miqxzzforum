import React from 'react';
import styled from 'styled-components';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../Chat/AuthContext'; // Импортируем useAuth

const Nav = styled.nav`
  background-color:rgb(24, 255, 16);
  color: white;
  padding: 10px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const NavTitle = styled.h1`
  margin-left: 45%;
  font-size: 1.5em;
`;

const NavLinks = styled.div`
  a {
    color: white;
    margin-left: 20px;
    text-decoration: none;
    transition: color 0.3s ease;

    &:hover {
      color: #f4f4f4;
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
            <NavTitle>MyForumGo</NavTitle>
            <NavLinks>
                {isAuthenticated ? (
                    <>
                        <Link to="/posts">Posts</Link>
                        <Link to="#" onClick={handleLogout}>Logout</Link>
                    </>
                ) : (
                    <>
                        <Link to="/login">Login</Link>
                        <Link to="/register">Register</Link>
                    </>
                )}
            </NavLinks>
        </Nav>
    );
};

export default Navbar;