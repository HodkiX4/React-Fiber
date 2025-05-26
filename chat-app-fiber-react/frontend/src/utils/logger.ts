import API from "../api/axiosInstance"

export const logError = (message: string) => {
  // Send the error message to a logging service or backend
  API.post('/log-error', {
    message,
    timeStamp: new Date().toISOString(),
  })
  .catch((err) => {
    console.error("Error logging failed:", err)
  });

  console.error("Logged Error:", message)
}