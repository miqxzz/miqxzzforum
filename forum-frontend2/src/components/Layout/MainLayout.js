import React from 'react';
import styled from 'styled-components';
import Navbar from './Navbar';

const LayoutContainer = styled.div`
  display: flex;
  min-height: 100vh;
`;

const SidebarSpace = styled.div`
  width: 220px;
  min-height: 100vh;
`;

const Content = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 40px 0 0 0;
  min-height: 100vh;
  margin-left: 220px;
`;

const MainLayout = ({ children }) => {
    // Разделяем Chat и остальной контент
    const chatChild = React.Children.toArray(children).find(child => child.type && child.type.name === 'Chat');
    const otherChildren = React.Children.toArray(children).filter(child => !(child.type && child.type.name === 'Chat'));

    return (
        <LayoutContainer>
            <Navbar />
            <SidebarSpace />
            <Content>
                <div style={{ width: '100%', maxWidth: 400, marginTop: '40px' }}>
                    {otherChildren}
                </div>
                <div style={{ width: '100%', maxWidth: 400, marginTop: '30px' }}>
                    {chatChild}
                </div>
            </Content>
        </LayoutContainer>
    );
};

export default MainLayout;