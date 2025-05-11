import { useState, useCallback, useEffect } from 'react';
import { apiService } from '../services/api';

interface AuthState {
  isAuthenticated: boolean;
}

export const useAuth = () => {
  const [state, setState] = useState<AuthState>(() => ({
    isAuthenticated: !!localStorage.getItem('access_token'),
  }));

  const login = useCallback(async (login: string, password: string) => {
    await apiService.login(login, password);
    setState({ isAuthenticated: true });
  }, []);

  const register = useCallback(async (login: string, password: string) => {
    await apiService.register(login, password);
    setState({ isAuthenticated: true });
  }, []);

  const logout = useCallback(async () => {
    await apiService.logout();
    setState({ isAuthenticated: false });
  }, []);

  const changePassword = useCallback(async (oldPassword: string, newPassword: string) => {
    await apiService.changePassword(oldPassword, newPassword);
  }, []);

  return {
    ...state,
    login,
    register,
    logout,
    changePassword,
  };
}; 