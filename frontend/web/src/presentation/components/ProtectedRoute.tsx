import { Navigate, useLocation } from 'react-router-dom';
import { useAuthContext } from '../contexts/AuthContext';
import type { ReactNode } from 'react';
import { Box, CircularProgress } from '@mui/material';

interface ProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated, isLoading } = useAuthContext();
  const location = useLocation();

  console.log('[ProtectedRoute] Checking authentication:', {
    isAuthenticated,
    isLoading,
    currentPath: location.pathname,
  });

  if (isLoading) {
    console.log('[ProtectedRoute] Loading authentication state...');
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!isAuthenticated) {
    console.log('[ProtectedRoute] User not authenticated, redirecting to login');
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  console.log('[ProtectedRoute] User authenticated, rendering protected content');
  return <>{children}</>;
};

export default ProtectedRoute; 