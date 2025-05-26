import { useState } from "react"
import { useToastContext } from "../context/ToastContext";
import { createRoomApi, deleteRoomApi, getRoomsApi } from "../api/roomApi";
import { useRoomsStore } from "../store/roomStore";
import handleErrors from "../utils/errorUtil";
import { ICreateRoomRequest } from "../models/room.model";

export const useRooms = () => {
  const [getLoading, setGetLoading] = useState<boolean>(false);
  const [createLoading, setCreateLoading] = useState<boolean>(false);
  const [deleteLoading, setDeleteLoading] = useState<boolean>(false);
  const { setRooms, addRoom, removeRoom } = useRoomsStore();
  const { showSuccess, showError } = useToastContext();
  
  const getRooms = async () => {
    if(getLoading) return;
    setGetLoading(true);
    try {
      const data = await getRoomsApi();
      setRooms(data.rooms);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setGetLoading(false);
    }
  }

  const createRoom = async (values: ICreateRoomRequest) => {
    if(createLoading) return;
    setCreateLoading(true)
    try {
      const data = await createRoomApi(values);
      addRoom(data.room);
      showSuccess(data.message);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setCreateLoading(false);
    }
  }

  const deleteRoom = async (roomId: string) => {
    if(deleteLoading) return;
    setDeleteLoading(true);
    try {
      const data = await deleteRoomApi(roomId);
      removeRoom(roomId);
      showSuccess(data.message);
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setDeleteLoading(false);
    }
  }

  
  return {
    getRooms, getLoading ,
    createRoom, createLoading,
    deleteRoom, deleteLoading
  }
}