import axios, { type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '../store/useAuthStore.ts';

const baseURL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080';

export const axiosClient = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor to attach the JWT token to every request
axiosClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().token;
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor to unwrap responses from Envelope { success: boolean, data: T }
axiosClient.interceptors.response.use(
  (response) => {
    // If response data is in the Envelope format, unwrap it
    if (response.data && response.data.success !== undefined && response.data.data !== undefined) {
      return { ...response, data: response.data.data };
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);