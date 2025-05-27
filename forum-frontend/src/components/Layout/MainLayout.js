import React from 'react';
import styled from 'styled-components';
import Navbar from './Navbar';

const LayoutContainer = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 100vh;
`;

const Content = styled.div`
  display: flex;
  flex: 1;
  padding: 20px;
`;

const PostListContainer = styled.div`
  flex: 1;
  margin-right: 320px; /* Ширина чата + отступ */
`;

const ChatContainer = styled.div`
  width: 300px;
  position: fixed;
  top: 60px; /* Высота шапки */
  right: 20px;
  height: calc(100vh - 80px); /* Высота окна минус шапка */
`;

const MainLayout = ({ children }) => {
    return (
        <LayoutContainer>
            <Navbar />
            <Content>
                <PostListContainer>
                    {React.Children.map(children, child => {
                        if (child.type && child.type.name !== 'Chat') {
                            return child;
                        }
                        return null;
                    })}
                </PostListContainer>
                <ChatContainer>
                    {React.Children.map(children, child => {
                        if (child.type && child.type.name === 'Chat') {
                            return child;
                        }
                        return null;
                    })}
                </ChatContainer>
            </Content>
        </LayoutContainer>
    );
};

export default MainLayout;