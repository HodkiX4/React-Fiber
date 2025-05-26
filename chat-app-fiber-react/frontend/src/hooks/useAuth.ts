import { useState } from "react";
import { loginApi, signupApi } from "../api/authApi";
import { ILoginRequest, ISignupRequest } from "../models/auth.model";
import useAuthStore from "../store/authStore";
import { getMeApi } from "../api/userApi";
import { validateLoginInput, validateSignupInput } from "../utils/validators";
import { useToastContext } from "../context/ToastContext";
import handleErrors from "../utils/errorUtil";

export const useAuth = () => {
  const { setToken, setUser } = useAuthStore();
  
  const { showError, showSuccess } = useToastContext();
  
  const [signupLoading, setSignupLoading] = useState<boolean>(false);
  const [loginLoading, setLoginLoading] = useState<boolean>(false);
  const [logoutLoading, setLogoutLoading] = useState<boolean>(false);

  const signup = async (values: ISignupRequest) => {
    if (signupLoading) return;
    
    const validationError = validateSignupInput(values)
    if(validationError != null) {
      showError(validationError);
      return;
    }

    setSignupLoading(true);
    try {
      const data = await signupApi(values);
      showSuccess(data.message)
    } catch (error) {
      showError(handleErrors(error))
    } finally {
      setSignupLoading(false);
    }
  };
  
  const login = async (values: ILoginRequest): Promise<boolean> => {
    if (loginLoading) return false; // Avoid duplicate login requests
    
    const validationError = validateLoginInput(values);
    if (validationError) {
      showError(validationError);
      return false; // Invalid input
    }
    
    setLoginLoading(true);
    try {
      const loginData = await loginApi(values);
      const userData = await getMeApi();
      setToken(loginData.token);
      setUser(userData);
      showSuccess(loginData.message);
      return true;
    } catch (error) {
      showError(handleErrors(error));
      return false;
    } finally {
      setLoginLoading(false);
    }
  };
  
  
  const logout = async () => {
    if (logoutLoading) return;
    setLogoutLoading(true)
    try {
      setToken(null);
      setUser(null)
    } catch (error) {
      showError(handleErrors(error));
    } finally {
      setLogoutLoading(false)
    }
  }

  return {
    signup, signupLoading,
    login, loginLoading,
    logout, logoutLoading
  };
};
