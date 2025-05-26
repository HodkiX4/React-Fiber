import { create } from 'zustand';

interface WebSocketStore {
  socket: WebSocket | null;
  typing: boolean;
  setSocket: (socket: WebSocket) => void;
  setTyping: (typing: boolean) => void;
}

export const useWebSocketStore = create<WebSocketStore>((set) => ({
  socket: null,
  typing: false,
  setSocket: (socket) => set({ socket }),
  setTyping: (typing) => set({ typing }),
}));
