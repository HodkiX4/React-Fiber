import { FormEvent, useState } from 'react';
import styles from './RoomSidebar.module.css'
import { useRoomsStore } from '../../store/roomStore';
import { useRooms } from '../../hooks/useRooms';
import { IUser } from '../../models/user.model';
import { useUsersStore } from '../../store/usersStore';
import useAuthStore from '../../store/authStore';
import { FaPlus } from 'react-icons/fa';
import { FaX } from 'react-icons/fa6';


const RoomsSidebar = () => {
  
  const [isCreating, setIsCreating] = useState(false);
  const [name, setName] = useState('');
  const [isPrivate, setIsPrivate] = useState(false);
  const [selectedMembers, setSelectedMembers] = useState<IUser[]>([]);


  const { createRoom, deleteRoom, getLoading } = useRooms();
  const { rooms, setActiveRoomId } = useRoomsStore();
  const { user } = useAuthStore();
  const { users } = useUsersStore();

  const friends = users.filter(u => u.friends?.includes(user!.id)) 
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
  
    const members = selectedMembers.map(m => m.id);
  
    await createRoom({ name, isPrivate, members });
  
    setName('');
    setIsPrivate(false);
    setSelectedMembers([]);
    setIsCreating(false);
  };
  

  if(getLoading) { 
    return <h2>Loading</h2>
  }
  return (
    <div className={styles.sidebar}>
      <h2 className={styles.header}>My Server</h2>
      <div className={styles.channelList}>
      <h4>TEXT CHANNELS</h4>
      {rooms.map((room) => (
        <div onClick={() => setActiveRoomId(room.id)} key={room.id} className={styles.channel}>
          <span>#{room.name}</span>
          <button onClick={() => deleteRoom(room.id)}><FaX/></button>
        </div>
      ))}
      </div>
      <button
        className={styles.openButton}
        onClick={() => setIsCreating(true)}
        type="submit"
      ><FaPlus/></button>
      {isCreating && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Create New Channel</h3>
            <form onSubmit={handleSubmit} className={styles.modalForm}>
              <label>
                Channel Name:
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Write something"
                  className={styles.input}
                />
              </label>

              <label>
                <input
                  type="checkbox"
                  checked={isPrivate}
                  onChange={(e) => setIsPrivate(e.target.checked)}
                />
                Private Channel
              </label>

              {/* Only show member selection if the channel is private */}
              {isPrivate && (
                <>
                  <label>
                    Add Members:
                    <select
                      onChange={(e) => {
                        const friendId = e.target.value;
                        const friend = friends.find((f) => f.id === friendId);
                        if (friend && !selectedMembers.find((m) => m.id === friend.id)) {
                          setSelectedMembers([...selectedMembers, friend]);
                        }
                        e.target.value = ''; // Reset selection
                      }}
                      className={styles.input}
                    >
                      <option value="">Select a friend...</option>
                      {friends
                        .filter((f) => !selectedMembers.find((m) => m.id === f.id))
                        .map((friend) => (
                          <option key={friend.id} value={friend.id}>
                            {friend.username}
                          </option>
                        ))}
                    </select>
                  </label>

                  <ul className={styles.friends}>
                    {selectedMembers.map((member) => (
                      <li key={member.id} className={styles.selectedMember}>
                        <span className={styles.username}>{member.username}</span>
                        <button
                          type="button"
                          className={styles.removeButton}
                          onClick={() =>
                            setSelectedMembers((prev) =>
                              prev.filter((m) => m.id !== member.id)
                            )
                          }
                        >
                          <FaX />
                        </button>
                      </li>
                    ))}
                  </ul>
                </>
              )}

              <div className={styles.modalButtons}>
                <button type="submit" className={styles.createButton}>
                  Create
                </button>
                <button
                  type="button"
                  onClick={() => setIsCreating(false)}
                  className={styles.cancelButton}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
      
      </div>
    )
  }
  
  /*
  <form className={styles.createForm} onSubmit={handleSubmit}>
      <input
      className={styles.input}
      placeholder="New channel..."
      value={newRoom}
      onChange={(e) => setNewRoom(e.target.value)}
      />
        <button className={styles.createButton} type="submit">+</button>
      </form>
  */
  export default RoomsSidebar