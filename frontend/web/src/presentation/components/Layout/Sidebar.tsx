import { Box, List, ListItem, ListItemIcon, ListItemText, styled, Divider, IconButton } from '@mui/material';
import { Home, LibraryMusic, Search, Logout, Person, AdminPanelSettings } from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthContext } from '../../contexts/AuthContext';
import { isAdmin } from '../../../utils/jwt';

const SidebarContainer = styled(Box)({
  width: 240,
  backgroundColor: 'background.paper',
  borderRight: '1px solid rgba(255, 255, 255, 0.1)',
  padding: '24px 0',
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
});

const Logo = styled(Box)({
  padding: '0 24px',
  marginBottom: '24px',
  fontSize: '24px',
  fontWeight: 700,
  color: 'primary.main',
});

const StyledListItem = styled(ListItem)(({ theme }) => ({
  padding: '12px 24px',
  '&:hover': {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
  },
  '&.Mui-selected': {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.15)',
    },
  },
}));

const navigationItems = [
  { text: 'Home', icon: <Home />, path: '/' },
  { text: 'Search', icon: <Search />, path: '/search' },
  { text: 'Library', icon: <LibraryMusic />, path: '/library' },
];

const ProfileSection = styled(Box)({
  marginTop: 'auto',
  padding: '16px 24px',
  borderTop: '1px solid rgba(255, 255, 255, 0.1)',
});

const Sidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { logout } = useAuthContext();

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const token = localStorage.getItem('access_token');
  const admin = isAdmin(token);

  return (
    <SidebarContainer>
      <Logo>Orpheon</Logo>
      <List>
        {navigationItems.map((item) => (
          <StyledListItem
            key={item.text}
            onClick={() => navigate(item.path)}
            sx={{
              cursor: 'pointer',
              backgroundColor: location.pathname === item.path ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
            }}
          >
            <ListItemIcon sx={{ color: 'inherit' }}>{item.icon}</ListItemIcon>
            <ListItemText primary={item.text} />
          </StyledListItem>
        ))}
        {admin && (
          <StyledListItem
            onClick={() => navigate('/admin')}
            sx={{
              cursor: 'pointer',
              backgroundColor: location.pathname === '/admin' ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
            }}
          >
            <ListItemIcon sx={{ color: 'inherit' }}><AdminPanelSettings /></ListItemIcon>
            <ListItemText primary="Админ-панель" />
          </StyledListItem>
        )}
      </List>

      <ProfileSection>
        <StyledListItem
          onClick={() => navigate('/profile')}
          sx={{
            cursor: 'pointer',
            backgroundColor: location.pathname === '/profile' ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
          }}
        >
          <ListItemIcon sx={{ color: 'inherit' }}>
            <Person />
          </ListItemIcon>
          <ListItemText primary="Profile" />
        </StyledListItem>
        <StyledListItem onClick={handleLogout} sx={{ cursor: 'pointer' }}>
          <ListItemIcon sx={{ color: 'inherit' }}>
            <Logout />
          </ListItemIcon>
          <ListItemText primary="Logout" />
        </StyledListItem>
      </ProfileSection>
    </SidebarContainer>
  );
};

export default Sidebar; 