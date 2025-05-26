import { useUsersStore } from '../../store/usersStore';
import styles from './UsersSidebar.module.css'
import useAuthStore from '../../store/authStore';


const UsersSidebar = () => {
  const { users } = useUsersStore();
  const { user } = useAuthStore();
  
  const friends = users.filter(u => user?.friends?.includes(u.id)) 
  
  return (
    <div className={styles.sidebar}>
      <div className={styles.header}>ONLINE — {friends.length}</div>
      <ul className={styles.userList}>
        {
          friends.map((u, index) => {
            return <li key={index} className={styles.user}>
              <div className={styles.avatarContainer}>
                <img src={u.avatarUrl} className={styles.avatar}/>
                <span 
                  className={
                    `${styles.statusDot} ${styles[u.status || 'offline']}`
                  }
                />
              </div>

              <span className={styles.name}>{u.username}</span>
            </li>
          })
      }
      </ul>
    </div>
  );
}

export default UsersSidebar