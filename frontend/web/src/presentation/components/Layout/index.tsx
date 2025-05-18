import { Box, styled } from '@mui/material';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import PlayerBar from './PlayerBar';

const LayoutContainer = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  height: '100vh',
  backgroundColor: 'background.default',
  overflow: 'hidden',
});

const MainContent = styled(Box)({
  flex: 1,
  display: 'flex',
  overflow: 'hidden',
  position: 'relative',
});

const ContentArea = styled(Box)({
  flex: 1,
  overflowY: 'auto',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
});

const Layout = () => {
  return (
    <LayoutContainer>
      <MainContent>
        <Sidebar />
        <ContentArea>
          <Outlet />
        </ContentArea>
      </MainContent>
      <PlayerBar />
    </LayoutContainer>
  );
};

export default Layout; 