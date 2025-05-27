import React, { createContext, useState, useContext, useEffect } from 'react';
import PropTypes from 'prop-types';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    // Загрузка данных при монтировании
    useEffect(() => {
        const loadAuthData = () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    setIsLoading(false);
                    return;
                }

                const id = localStorage.getItem('userID');
                const username = localStorage.getItem('username');
                const role = localStorage.getItem('role');

                if (!id || !username || !role) {
                    console.warn('Неполные данные пользователя в localStorage');
                    localStorage.removeItem('token');
                    setIsLoading(false);
                    return;
                }

                setUser({
                    id: Number(id) || 0,
                    username: username,
                    role: role,
                    token: token
                });
                setIsAuthenticated(true);
            } catch (error) {
                console.error('Ошибка при загрузке данных авторизации:', error);
                localStorage.clear();
            } finally {
                setIsLoading(false);
            }
        };

        loadAuthData();
    }, []);

    const login = (token, userData) => {
        try {
            if (!token || !userData?.id || !userData.username || !userData.role) {
                throw new Error('Неверные данные для входа');
            }

            localStorage.setItem('token', token);
            localStorage.setItem('userID', String(userData.id));
            localStorage.setItem('username', userData.username);
            localStorage.setItem('role', userData.role);

            setUser({
                id: Number(userData.id),
                username: userData.username,
                role: userData.role,
                token: userData.token
            });
            setIsAuthenticated(true);
        } catch (error) {
            console.error('Ошибка при входе:', error);
            throw error;
        }
    };

    const logout = () => {
        try {
            localStorage.clear();
            setUser(null);
            setIsAuthenticated(false);
        } catch (error) {
            console.error('Ошибка при выходе:', error);
        }
    };

    return (
        <AuthContext.Provider value={{
            user,
            isAuthenticated,
            isLoading,
            login,
            logout
        }}>
            {children}
        </AuthContext.Provider>
    );
};

AuthProvider.propTypes = {
    children: PropTypes.node.isRequired
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth должен использоваться внутри AuthProvider');
    }
    return context;
};