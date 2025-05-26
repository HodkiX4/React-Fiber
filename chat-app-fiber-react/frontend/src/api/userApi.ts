import { IGetUsersResponse, IUser } from "../models/user.model";
import API from "./axiosInstance";

export const getMeApi = async (): Promise<IUser> => {
  const response = await API.get<IUser>("/user/me")
  return response.data
}

export const getUsersApi = async (): Promise<IGetUsersResponse> => {
  const response = await API.get<IGetUsersResponse>("/user/all")
  return response.data
}