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

// Функция для сохранения состояния в localStorage
const saveAuthState = (state: AuthState) => {
  localStorage.setItem('authState', JSON.stringify({
    isAuthenticated: state.isAuthenticated,
    isAdmin: state.isAdmin,
    user: state.user,
  }));
};

// Функция для загрузки состояния из localStorage
const loadAuthState = (): Partial<AuthState> => {
  const savedState = localStorage.getItem('authState');
  if (savedState) {
    try {
      return JSON.parse(savedState);
    } catch (e) {
      console.error('Error parsing saved auth state:', e);
    }
  }
  return {};
};

export const useAuth = (initialPath: string = '/') => {
  const [state, setState] = useState<AuthState>(() => {
    const savedState = loadAuthState();
    return {
      isAuthenticated: savedState.isAuthenticated ?? false,
      isAdmin: savedState.isAdmin ?? false,
      user: savedState.user ?? null,
      isLoading: true,
    };
  });

  const checkAuth = async () => {
    console.log('[useAuth] Starting auth check...');
    try {
      console.log('[useAuth] Calling getMe...');
      const userData = await api.getMe();
      console.log('[useAuth] Received user data:', userData);

      if (userData && userData.id) {
        console.log('[useAuth] Setting authenticated state with user data:', userData);
        const newState = {
          isAuthenticated: true,
          isAdmin: userData.access_lvl === 2,
          user: userData,
          isLoading: false,
        };
        setState(newState);
        saveAuthState(newState);
        return true;
      } else {
        console.log('[useAuth] No valid user data received');
        const newState = {
          isAuthenticated: false,
          isAdmin: false,
          user: null,
          isLoading: false,
        };
        setState(newState);
        saveAuthState(newState);
        return false;
      }
    } catch (error) {
      console.log('[useAuth] Setting unauthenticated state due to error');
      const newState = {
        isAuthenticated: false,
        isAdmin: false,
        user: null,
        isLoading: false,
      };
      setState(newState);
      saveAuthState(newState);
      return false;
    }
  };

  // Проверяем авторизацию при монтировании компонента
  useEffect(() => {
    const initAuth = async () => {
      // Skip auth check for public routes
      if (PUBLIC_ROUTES.includes(initialPath)) {
        setState(prev => ({ ...prev, isLoading: false }));
        return;
      }
      await checkAuth();
    };
    initAuth();
  }, []); // Пустой массив зависимостей означает, что эффект выполнится только при монтировании

  const login = async (login: string, password: string) => {
    console.log('[useAuth] Login attempt...');
    try {
      const response = await api.login(login, password);
      console.log('[useAuth] Login response:', response);
      const isAuthenticated = await checkAuth();
      if (!isAuthenticated) {
        throw new Error('Failed to authenticate after login');
      }
      return response;
    } catch (error) {
      console.error('[useAuth] Login error:', error);
      throw error;
    }
  };

  const register = async (login: string, password: string) => {
    try {
      const response = await api.register(login, password);
      const isAuthenticated = await checkAuth();
      if (!isAuthenticated) {
        throw new Error('Failed to authenticate after registration');
      }
      return response;
    } catch (error) {
      console.error('[useAuth] Register error:', error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await api.logout();
      const newState = {
        isAuthenticated: false,
        isAdmin: false,
        user: null,
        isLoading: false,
      };
      setState(newState);
      saveAuthState(newState);
      localStorage.removeItem('authState');
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
      const newState = {
        ...state,
        user: response,
      };
      setState(newState);
      saveAuthState(newState);
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