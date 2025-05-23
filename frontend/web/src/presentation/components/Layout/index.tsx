import { Box, styled } from '@mui/material';
import Sidebar from './Sidebar';
import PlayerBar from './PlayerBar';
import { useState } from 'react';

const LayoutContainer = styled(Box)({
  display: 'flex',
  height: '100vh',
  overflow: 'hidden',
});

const MainContent = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isSidebarCollapsed'
})<{ isSidebarCollapsed: boolean }>(({ theme, isSidebarCollapsed }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  transition: theme.transitions.create('margin', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
}));

const ContentArea = styled(Box)({
  flex: 1,
  overflow: 'auto',
  padding: '24px',
  backgroundColor: 'background.default',
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
  borderRadius: '12px',
  margin: '12px 12px 12px 0',
  boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
  border: '1px solid rgba(255, 255, 255, 0.1)',
});

const PlayerBarContainer = styled(Box)({
  height: 80,
  backgroundColor: 'background.paper',
  borderTop: '1px solid',
  borderColor: 'divider',
  margin: '0 12px 12px 0',
  borderRadius: '12px',
  boxShadow: '0 -4px 6px -1px rgba(0, 0, 0, 0.1)',
});

const Layout = ({ children }: { children: React.ReactNode }) => {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  return (
    <LayoutContainer>
      <Sidebar 
        isCollapsed={isSidebarCollapsed} 
        onToggleCollapse={() => setIsSidebarCollapsed(!isSidebarCollapsed)} 
      />
      <MainContent isSidebarCollapsed={isSidebarCollapsed}>
        <ContentArea>
          {children}
        </ContentArea>
        <PlayerBarContainer>
          <PlayerBar />
        </PlayerBarContainer>
      </MainContent>
    </LayoutContainer>
  );
};

export default Layout; 