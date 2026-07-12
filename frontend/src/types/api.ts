import type { User } from './models';

export interface AuthResponse {
  token: string;
  user: User;
  message?: string;
}

export interface ApiError {
  message: string;
  statusCode: number;
}