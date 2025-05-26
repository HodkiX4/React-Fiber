import { IFriendRequest } from "../models/friend.model";
import API from "./axiosInstance";

interface IMessageResponse {
  message: string;
}

interface IPendingFriendRequestsResponse {
  incoming: IFriendRequest[];
  outgoing: IFriendRequest[];
}

export const sendFriendRequestApi = async (receiverId: string): Promise<IMessageResponse> => {
  const response = await API.post<IMessageResponse>("/friend-request", { receiverId });
  return response.data
}

export const acceptFriendRequestApi = async (friendReqId: string): Promise<IMessageResponse> => {
  const response = await API.patch<IMessageResponse>(`/friend-request/accept/${friendReqId}`);
  return response.data
}

export const rejectFriendRequestApi = async (friendReqId: string): Promise<IMessageResponse> => {
  const response = await API.patch<IMessageResponse>(`/friend-request/reject/${friendReqId}`);
  return response.data
}

export const getPendingFriendRequestsApi = async (): Promise<IPendingFriendRequestsResponse> => {
  const response = await API.get<IPendingFriendRequestsResponse>("/friend-request/pending");
  return response.data
}