import axios from "axios";

const handleErrors = (error: unknown): string => {
  console.error(error);
  
  if (axios.isAxiosError(error)) {
    return (
      error.response?.data?.message ||
      error.message ||
      "An unexpected Axios error occurred"
    );
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "Something went wrong";
};

export default handleErrors;