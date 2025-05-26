import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { IUser } from '../models/user.model';

interface AuthState {
  token: string | null;
  user: IUser | null;
  setToken: (token: string | null) => void;
  setUser: (user: IUser | null) => void;
  clearAuth: () => void;
}

const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      setToken: (token) => set({ token }),
      setUser: (user) => set({ user }),
      clearAuth: () => set({ token: null, user: null }),
    }),
    {
      name: 'auth-storage', // The key used for localStorage
      storage: createJSONStorage(() => localStorage), // Use localStorage for persistence
    }
  )
);

export default useAuthStore;
