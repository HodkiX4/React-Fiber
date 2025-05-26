// hooks/useWebSocketConnection.ts

import { useState } from 'react';
import { useWebSocketStore } from '../store/useWebSocketStore';

export const useWebSocketConnection = () => {
  const { setSocket, setTyping } = useWebSocketStore();
  const [socket, setSocketState] = useState<WebSocket | null>(null);

  // Function to establish WebSocket connection
  const connectWebSocket = () => {
    const ws = new WebSocket('ws://localhost:5001/ws');
    setSocketState(ws);
    setSocket(ws);

    // Listen for messages from the WebSocket server
    ws.onmessage = (event) => {
      const data = event.data;

      if (data.startsWith('typing:')) {
        const user = data.split(':')[1];
        console.log(`${user} is typing...`);
      } else if (data.startsWith('stop-typing:')) {
        const user = data.split(':')[1];
        console.log(`${user} stopped typing.`);
      }
    };

    // Cleanup on unmount
    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };
  };

  // Function to notify the server that the user is typing
  const sendTyping = () => {
    if (socket) {
      socket.send('typing');
      setTyping(true);
    }
  };

  // Function to stop sending typing notification after a delay
  const stopTyping = () => {
    setTyping(false);
  };

  return { connectWebSocket, sendTyping, stopTyping };
};
