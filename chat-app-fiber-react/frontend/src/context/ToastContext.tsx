import { createContext, useContext, useRef, ReactNode } from "react";
import toast from "react-hot-toast";
import { logError } from "../utils/logger"; // optional if you're logging errors

type ToastContextType = {
  showError: (message: string, log?: boolean) => void;
  showSuccess: (message: string) => void;
  showInfo: (message: string) => void;
  clearToast: () => void;
};

const ToastContext = createContext<ToastContextType>({
  showError: () => {},
  showSuccess: () => {},
  showInfo: () => {},
  clearToast: () => {},
});

export const ToastProvider = ({ children }: { children: ReactNode }) => {
  const toastId = useRef<string | undefined>(undefined);

  const clearToast = () => {
    if (toastId.current) {
      toast.remove(toastId.current);
      toastId.current = undefined;
    }
  };

  const showError = (message: string, log = true) => {
    clearToast();
    toastId.current = toast.error(message);
    if (log) logError(message); // optional error logging
  };

  const showSuccess = (message: string) => {
    clearToast();
    toastId.current = toast.success(message);
  };

  const showInfo = (message: string) => {
    clearToast();
    toastId.current = toast(message); // Default info toast
  };

  return (
    <ToastContext.Provider value={{ showError, showSuccess, showInfo, clearToast }}>
      {children}
    </ToastContext.Provider>
  );
};

export const useToastContext = () => useContext(ToastContext);
