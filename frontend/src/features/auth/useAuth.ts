import { useState, useEffect } from 'react';
import { authStore } from '../../store/authStore';
import { User, Tenant } from '../../types';

export function useAuth() {
  const [state, setState] = useState(() => authStore.getState());

  useEffect(() => {
    const unsubscribe = authStore.subscribe(() => {
      setState(authStore.getState());
    });
    return () => {
      unsubscribe();
    };
  }, []);

  return {
    user: state.user as User | null,
    tenant: state.tenant as Tenant | null,
    token: state.token as string | null,
    isAuthenticated: Boolean(state.token),
    setAuth: state.setAuth,
    logout: state.logout,
  };
}
