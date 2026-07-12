import axios, { type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '../store/useAuthStore.ts';

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

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