import axios from "axios";
import toast from "react-hot-toast";

export const showErrorToast = (
  error: unknown,
  fallback: string,
  setError: (message: string | null) => void
) => {
  toast.dismiss();

  let message: string;

  if (axios.isAxiosError(error)) {
    message = typeof error.response?.data?.message === 'string' ? error.response?.data?.message : fallback;
  } else if (error instanceof Error) {
    message = error.message || fallback;
  } else {
    message = fallback;
  }

  setError(message);
  toast.error(message);
};
