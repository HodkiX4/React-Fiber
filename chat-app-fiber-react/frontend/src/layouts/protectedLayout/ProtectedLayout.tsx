import { Navigate, Outlet } from 'react-router-dom'
import RoomsSidebar from '../../components/rooms_sidebar/RoomsSidebar'
import useAuthStore from '../../store/authStore';
import styles from './ProtectedLayout.module.css'
import UsersSidebar from '../../components/users_sidebar/UsersSidebar';
import { Header } from '../../components/header/Header';
import { useEffect } from 'react';
import { useUsers } from '../../hooks/useUsers';
import { useRooms } from '../../hooks/useRooms';
import { useFriends } from '../../hooks/useFriends';
const ProtectedLayout = () => {
  const { token } = useAuthStore();
  const { getMe, getUsers } = useUsers();
  const { getRooms } = useRooms();
  const { getPendingFriendRequests } = useFriends();
  useEffect(() => {
    if(token) {
      getMe();
      getUsers();
      getRooms();
      getPendingFriendRequests();
    }
  }, []);
  if(token === null) {
    return <Navigate to='/login'/>
  }
  return (
    <>
    <Header/>
    <main className={styles.main}>
      <RoomsSidebar/>
      <Outlet/>
      <UsersSidebar/>
    </main>
    </>
  )
}

export default ProtectedLayout