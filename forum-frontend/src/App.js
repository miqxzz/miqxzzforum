import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import PostList from './components/Posts/PostList';
import CreatePost from './components/Posts/CreatePost';
import MainLayout from './components/Layout/MainLayout';
import Chat from './components/Chat/chat'; // Импортируем компонент чата
import { AuthProvider, useAuth } from './components/Chat/AuthContext'; // Импортируем AuthProvider
import styled from 'styled-components';

const AppContainer = styled.div`
  text-align: center;
`;

const PrivateRoute = ({ children }) => {
    const { isAuthenticated } = useAuth();
    return children;
};

const App = () => {
    const onPostCreated = () => {
        console.log('Post was created');
    };

    return (
        <AuthProvider>
            <AppContainer>
                <Router>
                    <MainLayout>
                        <Routes>
                            <Route path="/login" element={<Login />} />
                            <Route path="/register" element={<Register />} />
                            <Route path="/posts" element={
                                <PrivateRoute>
                                    <CreatePost onPostCreated={onPostCreated} />
                                    <PostList />
                                </PrivateRoute>
                            } />
                            <Route path="/" element={<Navigate to="/login" />} />
                        </Routes>
                        <Chat /> {/* Включаем компонент чата */}
                    </MainLayout>
                </Router>
            </AppContainer>
        </AuthProvider>
    );
};

export default App;