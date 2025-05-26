import { useEffect, useState } from 'react';
import styles from './Header.module.css';
import { useUsersStore } from '../../store/usersStore';
import { useFriendRequestStore } from '../../store/friendReqStore';
import { FaBell, FaUsers } from 'react-icons/fa';
import { useFriends } from '../../hooks/useFriends';
import useAuthStore from '../../store/authStore';
import { useAuth } from '../../hooks/useAuth';

export const Header = () => {
  const { user } = useAuthStore();
  const { users } = useUsersStore();
  const { logout } = useAuth();
  const { incomingRequests, outgoingRequests } = useFriendRequestStore();

  const {
    sendFriendRequest,
    acceptFriendRequest,
    rejectFriendRequest,
  } = useFriends();

  console.log(user?.avatarUrl);
  
  const [showNotifications, setShowNotifications] = useState(false);
  const [showAddFriends, setShowAddFriends] = useState(false);
  const [showIncoming, setShowIncoming] = useState(false);
  
  
  
    const notFriends = users.filter(
    (u) =>
      !user?.friends?.includes(u.id) &&
      !outgoingRequests.some((req) => req.receiverId === u.id) &&
      !incomingRequests.some((req) => req.senderId === u.id)
  );
  
  const handleSendRequest = (receiverId: string) => {
    sendFriendRequest(receiverId);
  };

  const handleAccept = (id: string) => {
    acceptFriendRequest(id);
  };

  const handleReject = (id: string) => {
    rejectFriendRequest(id);
  };

  const handleBellClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowNotifications((prev) => !prev);
    setShowAddFriends(false);
  };

  const handleFriendsClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowAddFriends((prev) => !prev);
    setShowNotifications(false);
  };

  return (
    <header className={styles.header}>
      <div className={styles.icon} onClick={handleBellClick}>
        <FaBell />
        {showNotifications && (
          <div className={styles.dropdown} onClick={(e) => e.stopPropagation()}>
            {incomingRequests.length === 0 ? (
              <div>No new requests</div>
            ) : (
              incomingRequests.map((req) => {
                const fromUser = users.find((u) => u.id === req.senderId);
                return (
                  <div key={req.id} className={styles.request}>
                    <p>
                      <span>{fromUser?.username}</span> has sent you a friend request
                    </p>
                  </div>
                );
              })
            )}
          </div>
        )}
      </div>

      <div className={styles.icon} onClick={handleFriendsClick}>
        <FaUsers />
        {showAddFriends && (
          <div className={styles.dropdown} onClick={(e) => e.stopPropagation()}>
            <div className={styles.toggle}>
              <label>
                <input
                  type="checkbox"
                  checked={showIncoming}
                  onChange={() => setShowIncoming((prev) => !prev)}
                />
                Show Incoming Requests
              </label>
            </div>

            {showIncoming ? (
              incomingRequests.length === 0 ? (
                <div>No incoming requests</div>
              ) : (
                incomingRequests.map((req) => {
                  const fromUser = users.find((u) => u.id === req.senderId);
                  return (
                    <div key={req.id} className={styles.request}>
                      <span>{fromUser?.username}</span>
                      <button onClick={() => handleAccept(req.id)}>Accept</button>
                      <button onClick={() => handleReject(req.id)}>Reject</button>
                    </div>
                  );
                })
              )
            ) : (
              notFriends.length === 0 ? (
                <div>No users to add</div>
              ) : (
                notFriends.map((user) => (
                  <div key={user.id} className={styles.request}>
                    <span>{user.username}</span>
                    <button onClick={() => handleSendRequest(user.id)}>Add</button>
                  </div>
                ))
              )
            )}
          </div>
        )}
      </div>
      <div className={styles.userInfo}>
        <img
          src={user?.avatarUrl}
          alt="User Avatar"
          className={styles.avatar}
          onError={(e) =>
            ((e.target as HTMLImageElement).src =
              'https://ui-avatars.com/api/?name=User&background=random')
          }
        />
        <span className={styles.username}>{user?.username}</span>
      </div>
      <button onClick={logout}>Log out</button>
    </header>
  );
};
