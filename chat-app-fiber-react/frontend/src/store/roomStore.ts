import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';
import { IRoom } from "../models/room.model";
import { IMessage } from '../models/message.model';

interface RoomsState {
  rooms: IRoom[];
  activeRoomId: string | null;
  messages: Record<string, IMessage[]>;
  setRooms: (rooms: IRoom[]) => void;
  addRoom: (room: IRoom) => void;
  removeRoom: (roomId: string) => void;
  setActiveRoomId: (roomId: string | null) => void;
  addMessage: (roomId: string, message: IMessage) => void;
  removeMessage: (roomId: string, messageId: string) => void;
  setMessages: (roomId: string, messages: IMessage[]) => void;
  typingUsers: Record<string, string[]>; // roomId -> list of usernames
  setTypingUsers: (roomId: string, users: string[]) => void;
  addTypingUser: (roomId: string, username: string) => void;
  removeTypingUser: (roomId: string, username: string) => void;
}


export const useRoomsStore = create<RoomsState>()(
  persist(
    (set) => ({
      rooms: [],
      activeRoomId: null,
      messages: {},
      typingUsers: {},
      setRooms: (rooms) => set({ rooms }),
      addRoom: (room) => set((state) => ({ rooms: [...state.rooms, room] })),
      setActiveRoomId: (roomId) => set({ activeRoomId: roomId }),
      addMessage: (roomId, message) =>
        set((state) => ({
          messages: {
            ...state.messages,
            [roomId]: [...(state.messages[roomId] || []), message],
          },
        })),
      removeMessage: (roomId, messageId) =>
        set((state) => ({
          messages: {
            ...state.messages,
            [roomId]: state.messages[roomId]?.filter(m => m.id !== messageId),
          },
        })),
      setMessages: (roomId, messages) =>
        set((state) => ({
          messages: {
            ...state.messages,
            [roomId]: messages,
          },
        })),
      removeRoom: (roomId) =>
        set((state) => {
          const newRooms = state.rooms.filter(room => room.id !== roomId); // Remove the room from the rooms array
          const { [roomId]: _, ...remainingMessages } = state.messages;
          return {
            rooms: newRooms,
            messages: remainingMessages,
          };
        }),
      setTypingUsers: (roomId, users) =>
        set(state => ({
          typingUsers: { ...state.typingUsers, [roomId]: users }
        })),
      addTypingUser: (roomId, username) =>
        set(state => {
          const existing = state.typingUsers[roomId] || [];
          if (existing.includes(username)) return state;
          return {
            typingUsers: {
              ...state.typingUsers,
              [roomId]: [...existing, username]
            }
          };
      }),
      removeTypingUser: (roomId, username) =>
        set(state => {
          const existing = state.typingUsers[roomId] || [];
          return {
            typingUsers: {
              ...state.typingUsers,
              [roomId]: existing.filter(u => u !== username)
            }
          };
      }),
    }),
    {
      name: 'rooms-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);

