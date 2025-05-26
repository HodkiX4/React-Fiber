import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';
import { IUser } from '../models/user.model';

interface UsersState {
  users: IUser[];
  setUsers: (users: IUser[]) => void;
  addUser: (user: IUser) => void;
}

export const useUsersStore = create<UsersState>()(
  persist(
    (set) => ({
      users: [],
      setUsers: (users) => set({ users }),
      addUser: (user) => set((state) => ({ users: [...state.users, user] })),
    }),
    {
      name: 'users-storage', // the key used to store the state in localStorage
      storage: createJSONStorage(() => localStorage), // persist in localStorage
    }
  )
);
