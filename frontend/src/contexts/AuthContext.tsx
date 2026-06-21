'use client';

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import type { UserResponse } from '@/types';

interface AuthContextType {
  user: UserResponse | null;
  token: string | null;
  setAuth: (token: string, user: UserResponse) => void;
  logout: () => void;
  isLoggedIn: boolean;
  isReady: boolean;
}

const AuthContext = createContext<AuthContextType>({
  user: null, token: null, setAuth: () => {}, logout: () => {}, isLoggedIn: false, isReady: false,
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    if (savedToken && savedUser) {
      setToken(savedToken);
      setUser(JSON.parse(savedUser));
    }
    setIsReady(true);
  }, []);

  const setAuth = useCallback((t: string, u: UserResponse) => {
    localStorage.setItem('token', t);
    localStorage.setItem('user', JSON.stringify(u));
    setToken(t);
    setUser(u);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken(null);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, token, setAuth, logout, isLoggedIn: !!token, isReady }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
