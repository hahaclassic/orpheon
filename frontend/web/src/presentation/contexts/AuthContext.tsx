import { createContext, useContext, type ReactNode } from 'react';
import { useAuth } from '../hooks/useAuth';

// Separate interfaces for better interface segregation
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
  isLoading: boolean;
  user: User | null;
}

interface AuthActions {
  login: (login: string, password: string) => Promise<User>;
  register: (login: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
  changePassword: (oldPassword: string, newPassword: string) => Promise<void>;
}

interface AuthContextType extends AuthState, AuthActions {}

// Create separate contexts for state and actions
const AuthStateContext = createContext<AuthState | undefined>(undefined);
const AuthActionsContext = createContext<AuthActions | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const auth = useAuth();

  // Split auth into state and actions
  const state: AuthState = {
    isAuthenticated: auth.isAuthenticated,
    isAdmin: auth.isAdmin,
    isLoading: auth.isLoading,
    user: auth.user,
  };

  const actions: AuthActions = {
    login: auth.login,
    register: auth.register,
    logout: auth.logout,
    changePassword: auth.changePassword,
  };

  return (
    <AuthStateContext.Provider value={state}>
      <AuthActionsContext.Provider value={actions}>
        {children}
      </AuthActionsContext.Provider>
    </AuthStateContext.Provider>
  );
};

// Custom hooks for accessing specific parts of auth
export const useAuthState = () => {
  const context = useContext(AuthStateContext);
  if (context === undefined) {
    throw new Error('useAuthState must be used within an AuthProvider');
  }
  return context;
};

export const useAuthActions = () => {
  const context = useContext(AuthActionsContext);
  if (context === undefined) {
    throw new Error('useAuthActions must be used within an AuthProvider');
  }
  return context;
};

// Main hook that combines both state and actions
export const useAuthContext = () => {
  const state = useAuthState();
  const actions = useAuthActions();

  return {
    ...state,
    ...actions,
  };
}; 