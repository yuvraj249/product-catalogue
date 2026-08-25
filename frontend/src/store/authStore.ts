import { User, Tenant } from '../types';

interface AuthState {
  user: User | null;
  tenant: Tenant | null;
  token: string | null;
  setAuth: (user: User | null, tenant: Tenant | null, token: string | null) => void;
  logout: () => void;
}

class SimpleAuthStore {
  private listeners: Set<() => void> = new Set();
  private state: AuthState;

  constructor() {
    const savedToken = localStorage.getItem('erp_token');
    const savedUser = localStorage.getItem('erp_user');
    const savedTenant = localStorage.getItem('erp_tenant');

    this.state = {
      user: savedUser ? JSON.parse(savedUser) : null,
      tenant: savedTenant ? JSON.parse(savedTenant) : null,
      token: savedToken,
      setAuth: this.setAuth.bind(this),
      logout: this.logout.bind(this),
    };
  }

  public getState = (): AuthState => this.state;

  public subscribe = (listener: () => void) => {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  };

  public setAuth = (user: User | null, tenant: Tenant | null, token: string | null) => {
    this.state.user = user;
    this.state.tenant = tenant;
    this.state.token = token;

    if (token) {
      localStorage.setItem('erp_token', token);
    } else {
      localStorage.removeItem('erp_token');
    }

    if (user) {
      localStorage.setItem('erp_user', JSON.stringify(user));
    } else {
      localStorage.removeItem('erp_user');
    }

    if (tenant) {
      localStorage.setItem('erp_tenant', JSON.stringify(tenant));
    } else {
      localStorage.removeItem('erp_tenant');
    }

    this.notify();
  };

  public logout = () => {
    this.setAuth(null, null, null);
  };

  private notify = () => {
    this.listeners.forEach((l) => l());
  };
}

export const authStore = new SimpleAuthStore();
