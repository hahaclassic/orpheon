import { useState, useEffect } from 'react';
import { api } from '../services/api';

interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface AuthState {
  isAuthenticated: boolean;
  isAdmin: boolean;
  user: User | null;
  isLoading: boolean;
}

// List of public routes that don't require authentication
const PUBLIC_ROUTES = ['/', '/search', '/login', '/register'];

export const useAuth = (initialPath: string = '/') => {
  const [state, setState] = useState<AuthState>({
    isAuthenticated: false,
    isAdmin: false,
    user: null,
    isLoading: true,
  });

  const checkAuth = async () => {
    console.log('[useAuth] Starting auth check...');
    try {
      console.log('[useAuth] Calling getMe...');
      const userData = await api.getMe();
      console.log('[useAuth] Received user data:', userData);

      if (userData && userData.id) {
        console.log('[useAuth] Setting authenticated state with user data:', userData);
        setState({
          isAuthenticated: true,
          isAdmin: userData.access_lvl === 2,
          user: userData,
          isLoading: false,
        });
      } else {
        console.log('[useAuth] No valid user data received');
        setState({
          isAuthenticated: false,
          isAdmin: false,
          user: null,
          isLoading: false,
        });
      }
    } catch (error) {
      console.log('[useAuth] Setting unauthenticated state due to error');
      setState({
        isAuthenticated: false,
        isAdmin: false,
        user: null,
        isLoading: false,
      });
    }
  };

  // Проверяем авторизацию при монтировании компонента
  useEffect(() => {
    checkAuth();
  }, []); // Пустой массив зависимостей означает, что эффект выполнится только при монтировании

  // Дополнительная проверка при изменении пути
  useEffect(() => {
    console.log('[useAuth] Auth effect triggered');
    // Skip auth check for public routes
    if (PUBLIC_ROUTES.includes(initialPath)) {
      setState(prev => ({ ...prev, isLoading: false }));
      return;
    }
    checkAuth();
  }, [initialPath]); // Re-run when path changes

  const login = async (login: string, password: string) => {
    console.log('[useAuth] Login attempt...');
    try {
      const response = await api.login(login, password);
      console.log('[useAuth] Login response:', response);
      await checkAuth();
      return response;
    } catch (error) {
      console.error('[useAuth] Login error:', error);
      throw error;
    }
  };

  const register = async (login: string, password: string) => {
    try {
      const response = await api.register(login, password);
      await checkAuth();
      return response;
    } catch (error) {
      console.error('[useAuth] Register error:', error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await api.logout();
      setState({
        isAuthenticated: false,
        isAdmin: false,
        user: null,
        isLoading: false,
      });
    } catch (error) {
      console.error('[useAuth] Logout error:', error);
      throw error;
    }
  };

  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      const response = await api.changePassword(oldPassword, newPassword);
      return response;
    } catch (error) {
      console.error('[useAuth] Change password error:', error);
      throw error;
    }
  };

  const updateUser = async (user: User) => {
    try {
      const response = await api.updateUser(user);
      setState(prev => ({
        ...prev,
        user: response,
      }));
      return response;
    } catch (error) {
      console.error('[useAuth] Update user error:', error);
      throw error;
    }
  };

  return {
    ...state,
    login,
    register,
    logout,
    changePassword,
    updateUser,
  };
}; 