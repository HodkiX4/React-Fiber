import API from "./axiosInstance";
import {
  ILoginRequest,
  ILoginResponse,
  ISignupRequest,
  ISignupResponse
} from "../models/auth.model";

export const signupApi = async (values: ISignupRequest): Promise<ISignupResponse> => {  
  const response = await API.post<ISignupResponse>("/auth/register", values);
  return response.data
};

export const loginApi = async (values: ILoginRequest): Promise<ILoginResponse> => {
  const response = await API.post<ILoginResponse>("/auth/login", values);
  return response.data
};