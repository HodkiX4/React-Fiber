import API from "./axiosInstance";
import {
  IDeleteMessageResponse,
  IGetMessagesResponse,
  ISendMessageRequest,
  ISendMessageResponse
} from "../models/message.model";

export const sendMessageApi = async (values: ISendMessageRequest): Promise<ISendMessageResponse> => {
  const response = await API.post<ISendMessageResponse>('/chat/message', values)
  return response.data
}

export const getMessagesApi = async (roomId: string): Promise<IGetMessagesResponse> => {
  const response = await API.get<IGetMessagesResponse>(`/chat/message/${roomId}`)
  return response.data
}

export const deleteMessagesApi = async (messageId: string): Promise<IDeleteMessageResponse> => {
  const response = await API.delete<IDeleteMessageResponse>(`/chat/message/${messageId}`)
  return response.data
}
