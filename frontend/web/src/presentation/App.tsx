import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { AuthProvider } from './contexts/AuthContext';
import { PlayerProvider } from './contexts/PlayerContext';
import ErrorBoundary from './components/ErrorBoundary';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import Home from './pages/Home';
import Profile from './pages/Profile';
import Me from './pages/Me';
import GenreList from './pages/Admin/GenreList';
import LicenseList from './pages/Admin/LicenseList';
import ArtistList from './pages/Admin/ArtistList';
import AlbumList from './pages/Admin/AlbumList';
import { theme } from './theme';

const ProtectedRoutes = () => {
  return (
    <ProtectedRoute>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/me" element={<Me />} />
        
        <Route path="/admin" element={<AdminRoute><Outlet /></AdminRoute>}>
          <Route index element={<Navigate to="genres" replace />} />
          <Route path="genres" element={<GenreList />} />
          <Route path="licenses" element={<LicenseList />} />
          <Route path="artists" element={<ArtistList />} />
          <Route path="albums" element={<AlbumList />} />
        </Route>
      </Routes>
    </ProtectedRoute>
  );
};

const App = () => {
  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <AuthProvider>
            <PlayerProvider>
              <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <Routes>
                  <Route path="/login" element={<Navigate to="/auth/login" replace />} />
                  <Route path="/register" element={<Navigate to="/auth/register" replace />} />
                  
                  <Route element={<Layout />}>
                    <Route path="/*" element={<ProtectedRoutes />} />
                  </Route>
                </Routes>
              </BrowserRouter>
            </PlayerProvider>
          </AuthProvider>
        </LocalizationProvider>
      </ThemeProvider>
    </ErrorBoundary>
  );
};

export default App; 