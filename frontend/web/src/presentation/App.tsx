import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { AuthProvider } from './contexts/AuthContext';
import { PlayerProvider } from './contexts/PlayerContext';
import ErrorBoundary from './components/ErrorBoundary';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import Login from './pages/Auth/Login';
import Register from './pages/Auth/Register';
import Home from './pages/Home';
import Profile from './pages/Profile';
import GenreList from './pages/Admin/GenreList';
import LicenseList from './pages/Admin/LicenseList';
import ArtistList from './pages/Admin/ArtistList';
import AlbumList from './pages/Admin/AlbumList';
import theme from './theme';

const App = () => {
  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <AuthProvider>
            <PlayerProvider>
              <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <Routes>
                  <Route path="/login" element={<Login />} />
                  <Route path="/register" element={<Register />} />
                  
                  <Route element={<Layout />}>
                    <Route element={<ProtectedRoute />}>
                      <Route path="/" element={<Home />} />
                      <Route path="/profile" element={<Profile />} />
                      
                      <Route path="/admin" element={<AdminRoute />}>
                        <Route index element={<Navigate to="genres" replace />} />
                        <Route path="genres" element={<GenreList />} />
                        <Route path="licenses" element={<LicenseList />} />
                        <Route path="artists" element={<ArtistList />} />
                        <Route path="albums" element={<AlbumList />} />
                      </Route>
                    </Route>
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