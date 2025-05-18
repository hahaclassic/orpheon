import { Box, List, ListItem, ListItemIcon, ListItemText, styled, Divider, IconButton } from '@mui/material';
import { Home, LibraryMusic, Search, Logout, Person, AdminPanelSettings } from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthContext } from '../../contexts/AuthContext';

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
  display: 'flex',
  alignItems: 'center',
  gap: '12px',
  '& svg': {
    width: '32px',
    height: '32px',
    display: 'block',
  },
  '& span': {
    fontSize: '24px',
    fontWeight: 700,
    color: 'primary.main',
    lineHeight: 1,
  },
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
  const { logout, isAdmin } = useAuthContext();

  console.log('Sidebar isAdmin value:', isAdmin); // Debug log

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <SidebarContainer>
      <Logo
        sx={{ cursor: 'pointer' }}
        onClick={() => navigate('/')}
      >
        <svg viewBox="0 0 310 350" xmlns="http://www.w3.org/2000/svg">
          <g transform="translate(0,350) scale(0.1,-0.1)" stroke="none">
            <path fill="currentColor" d="M0 1750 l0 -1750 1550 0 1550 0 0 1750 0 1750 -1550 0 -1550 0 0
-1750z m996 1531 c144 -64 217 -213 217 -443 l0 -118 74 0 73 0 0 133 c0 157
8 177 70 177 62 0 70 -20 70 -177 l0 -133 85 0 85 0 0 219 c0 244 2 251 64
251 65 0 66 -1 68 -245 l3 -220 93 -3 92 -3 0 133 c0 158 8 178 70 178 62 0
70 -20 70 -177 l0 -133 54 0 c45 0 58 -4 70 -22 40 -56 4 -108 -74 -108 l-50
0 -2 -746 -3 -746 -24 -19 c-28 -23 -59 -24 -89 -3 l-22 15 0 750 0 749 -95 0
-95 0 0 -873 c0 -666 -3 -876 -12 -885 -7 -7 -31 -12 -55 -12 -32 0 -45 5 -53
19 -7 13 -10 308 -10 885 l0 866 -85 0 -85 0 0 -797 c0 -439 -4 -803 -8 -809
-27 -41 -106 -40 -123 3 -5 14 -9 345 -9 814 l0 789 -93 0 -93 0 -18 -55 c-9
-31 -73 -167 -141 -303 -198 -394 -242 -526 -252 -742 -7 -157 13 -264 74
-390 39 -82 58 -107 137 -185 107 -106 188 -151 341 -191 120 -32 359 -44 472
-25 308 54 582 264 702 539 70 160 99 369 81 580 -14 157 -34 251 -86 405 -29
83 -41 134 -35 140 14 14 23 5 78 -85 112 -184 196 -381 230 -539 26 -121 26
-449 -1 -569 -77 -354 -317 -659 -647 -825 -343 -173 -807 -179 -1163 -15
-499 229 -744 776 -586 1310 32 109 128 314 203 435 24 39 108 162 187 275
155 223 193 292 211 393 10 56 10 76 -4 133 -10 37 -19 69 -21 71 -2 3 -23 -8
-46 -23 -57 -38 -150 -45 -197 -15 -46 29 -83 107 -83 176 0 100 55 173 156
206 74 24 161 20 230 -10z m1884 -96 c53 -27 53 -99 1 -208 -79 -165 -282
-316 -394 -292 -106 23 -127 81 -81 227 37 115 145 228 255 268 67 24 178 26
219 5z"/>
          </g>
        </svg>
        <span>Orpheon</span>
      </Logo>

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
        {isAdmin && (
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