import { Navigate } from 'react-router-dom';
import type { ReactNode } from 'react';
import { isAdmin } from '../../utils/jwt';

interface AdminRouteProps {
  children: ReactNode;
}

const AdminRoute = ({ children }: AdminRouteProps) => {
  const token = localStorage.getItem('access_token');
  if (!isAdmin(token)) {
    return <Navigate to="/" replace />;
  }
  return <>{children}</>;
};

export default AdminRoute; 