// axiosInstance.ts
import axios from "axios";

const API = axios.create({
  baseURL: '/api', // Directly point to the backend server
});

// Optionally add interceptors or other configurations
API.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token && config.url && !config.url.includes('/auth/')) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default API;
