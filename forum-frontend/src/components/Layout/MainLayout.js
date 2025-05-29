import React from 'react';
import styled from 'styled-components';
import Navbar from './Navbar';

const Sidebar = styled.div`
  width: 240px;
  min-width: 200px;
  background: #a678f7;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  box-shadow: 2px 0 8px #e1d5ee44;
  font-family: 'Montserrat', sans-serif;
  padding-top: 0;
`;

const LayoutContainer = styled.div`
  display: flex;
  min-height: 100vh;
`;

const Content = styled.div`
  flex: 1;
  padding: 40px 40px 40px 0;
  background: #f8f4fc;
`;

const MainLayout = ({ children }) => {
    return (
        <LayoutContainer>
            <Sidebar>
                <Navbar />
            </Sidebar>
            <Content>
                {children}
            </Content>
        </LayoutContainer>
    );
};

export default MainLayout;