import API from "./axiosInstance"
import {
  ICreateRoomRequest,
  ICreateRoomResponse,
  IDeleteRoomResponse,
  IGetRoomsResponse,
} from "../models/room.model";

export const getRoomsApi = async (): Promise<IGetRoomsResponse> => {
  const response = await API.get<IGetRoomsResponse>("/chat/rooms")  
  return response.data
}

export const createRoomApi = async (data: ICreateRoomRequest): Promise<ICreateRoomResponse> => {
  const response = await API.post<ICreateRoomResponse>("/chat/rooms", data)
  return response.data
}

export const deleteRoomApi = async (roomId: string): Promise<IDeleteRoomResponse> => {
  const response = await API.delete<IDeleteRoomResponse>(`/chat/rooms/${roomId}`)
  return response.data
}