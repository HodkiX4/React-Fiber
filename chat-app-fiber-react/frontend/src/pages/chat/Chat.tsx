import { ChangeEvent, useEffect, useRef, useState } from 'react';
import { useRoomsStore } from '../../store/roomStore';
import { useUsersStore } from '../../store/usersStore';
import styles from './Chat.module.css';
import { useMessages } from '../../hooks/useMessages';
import useAuthStore from '../../store/authStore';
import { FaX } from 'react-icons/fa6';
import { useWebSocketStore } from '../../store/useWebSocketStore';

const Chat = () => {
  const { user: authUser } = useAuthStore();
  const { users } = useUsersStore();
  const { messages, rooms, activeRoomId, typingUsers } = useRoomsStore();
  const { socket } = useWebSocketStore(); 
  const [content, setContent] = useState('');
  const { sendMessage, getMessages, deleteMessage } = useMessages();
  
  useEffect(() => {
    getMessages();
  }, []);
  
  useEffect(() => {
    const socket = useWebSocketStore.getState().socket;
    if (!socket || !activeRoomId || !authUser) return;
  
    const handleTyping = (data: { roomId: string; username: string }) => {
      if (data.roomId !== activeRoomId || data.username === authUser.username) return;
      useRoomsStore.getState().addTypingUser(data.roomId, data.username);
      setTimeout(() => {
        useRoomsStore.getState().removeTypingUser(data.roomId, data.username);
      }, 3000); // Remove after 3s of inactivity
    };
  
    socket.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'typing') {
          handleTyping(data.payload);
        }
      } catch (e) {
        console.log(e);
      }
    });
  
    return () => {
      socket?.close(); // Optionally remove this line if your socket is global
    };
  }, [activeRoomId, authUser]);
  
  
  const currentMessages = messages[activeRoomId ?? ''] || [];
  const activeRoom = rooms.find(room => room.id === activeRoomId);
  const roomName = activeRoom?.name;
  const roomAdminsOrModerators = activeRoom?.admins.concat(activeRoom.moderators);

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setContent(value);
  
    if (socket && socket.readyState === WebSocket.OPEN && authUser && activeRoomId) {
      socket.send(
        JSON.stringify({
          type: 'typing',
          payload: {
            roomId: activeRoomId,
            username: authUser.username,
          },
        })
      );
    }
  };

  const handleSendMessage = () => {
    if (content.trim() === '') return;
    sendMessage(content)
    setContent('');
  };

  if (activeRoomId === null) {
    return <h1>Select a room</h1>;
  }

  return (
    <div className={styles.chat}>
      <div className={styles.header}>
        <h2># {roomName}</h2>
      </div>

      <div className={styles.messages}>
        {currentMessages.map(msg => {
          const msgUser = msg.senderId === authUser?.id ? authUser : users.find(u => u.id === msg.senderId);
          const authorizedToDelete = authUser === msgUser || roomAdminsOrModerators?.includes(authUser!.id);
          const { id, createdAt, content } = msg;
          const formatedTime = new Date(createdAt).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
          });

          return (
            <div key={id} className={styles.message}>
              <img src={msgUser?.avatarUrl} className={styles.avatar} />
              <div>
                <span className={styles.sender}>{msgUser?.username}</span>
                <span className={styles.time}>{formatedTime}</span>
                <p className={styles.text}>{content}</p>
              </div>
              {authorizedToDelete && (
                <button onClick={() => deleteMessage(id)}>
                  <FaX />
                </button>
              )}
            </div>
          );
        })}
        {typingUsers[activeRoomId]?.length > 0 && (
          <div className={styles.typingIndicator}>
            {typingUsers[activeRoomId]
              .map(username => username === authUser?.username ? 'You' : username)
              .join(', ')} {typingUsers[activeRoomId].length === 1 ? 'is' : 'are'} typing...
          </div>
        )}
      </div>

      <div className={styles.inputBox}>
        <button className={styles.addBtn}>+</button>
        <input
          className={styles.input}
          type="text"
          onChange={handleInputChange}
          value={content}
          placeholder={`Message #${roomName}`}
        />
        <button onClick={handleSendMessage} className={styles.sendBtn}>
          ➤
        </button>
      </div>
    </div>
  );
};

export default Chat;
