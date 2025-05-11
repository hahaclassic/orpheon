import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import { theme } from './presentation/styles/theme';
import Layout from './presentation/components/Layout';
import { PlayerProvider } from './presentation/contexts/PlayerContext';
import { AuthProvider } from './presentation/contexts/AuthContext';
import ErrorBoundary from './presentation/components/ErrorBoundary';
import Home from './presentation/pages/Home';
import Library from './presentation/pages/Library';
import Search from './presentation/pages/Search';
import Login from './presentation/pages/Login';
import Register from './presentation/pages/Register';
import Profile from './presentation/pages/Profile';
import ProtectedRoute from './presentation/components/ProtectedRoute';
import { AdminPanel } from './presentation/pages/Admin/AdminPanel';
import AdminRoute from './presentation/components/AdminRoute';
import AlbumCreate from './presentation/pages/Admin/AlbumCreate';
import GenreList from './presentation/pages/Admin/GenreList';
import LicenseList from './presentation/pages/Admin/LicenseList';
import ArtistCreate from './presentation/pages/Admin/ArtistCreate';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';

function App() {
  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <AuthProvider>
          <PlayerProvider>
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <Router>
                <Routes>
                  <Route path="/login" element={<Login />} />
                  <Route path="/register" element={<Register />} />
                  <Route element={<ProtectedRoute><Layout><Outlet /></Layout></ProtectedRoute>}>
                    <Route path="/" element={<Home />} />
                    <Route path="/library" element={<Library />} />
                    <Route path="/search" element={<Search />} />
                    <Route path="/profile" element={<Profile />} />
                    <Route path="/admin" element={<AdminPanel />} />
                    <Route path="/admin" element={<AdminRoute><Outlet /></AdminRoute>}>
                      <Route path="albums" element={<AlbumCreate />} />
                      <Route path="genres" element={<GenreList />} />
                      <Route path="licenses" element={<LicenseList />} />
                      <Route path="artists" element={<ArtistCreate />} />
                    </Route>
                  </Route>
                </Routes>
              </Router>
            </LocalizationProvider>
          </PlayerProvider>
        </AuthProvider>
      </ThemeProvider>
    </ErrorBoundary>
  );
}

export default App;
