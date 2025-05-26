import { useState } from "react"
import { useUsersStore } from "../store/usersStore";
import { useToastContext } from "../context/ToastContext";
import handleErrors from "../utils/errorUtil";
import { getMeApi, getUsersApi } from "../api/userApi";
import useAuthStore from "../store/authStore";

export const useUsers = () => {
  const [usersLoading, setUsersLoading] = useState<boolean>(false);
  const [meLoading, setMeLoading] = useState<boolean>(false);
  const { setUser } = useAuthStore();
  const { setUsers } = useUsersStore();
  const { showSuccess, showError } = useToastContext();

  const getMe = async () => {
    setMeLoading(true);
    try {
      const data = await getMeApi();
      setUser(data);
      showSuccess("I got my self");
    } catch (error) {
      showError(handleErrors(error));
    }
  }
  
  const getUsers = async () => {
    setUsersLoading(true);
    try {
      const data = await getUsersApi();
      setUsers(data.users);
      showSuccess("Users fetched");
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setUsersLoading(false);
    }
  }

  return {
    getMe, meLoading, 
    getUsers, usersLoading
  }
}