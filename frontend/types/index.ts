// types/index.ts

export interface User {
  id: string;
  email: string;
  name: string;
  lastName: string;
  workPosition: string;
  salary: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface AuthResponse {
  message: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  lastName: string;
  email: string;
  workPosition: string;
  salary: number;
  password: string;
}

export interface CreateUserRequest {
  email: string;
  name: string;
  lastName: string;
  workPosition: string;
  salary: number;
  password: string;
}

export interface UpdateUserRequest {
  email: string;
  name: string;
  lastName: string;
  workPosition: string;
  salary: number;
}

export interface UpdatePasswordRequest {
  email: string;
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

export interface DeleteUserRequest {
  userId: string;
}

export interface LogEntry {
  id: string;
  timestamp: string;
  username: string;
  action: string;
  status: string;
}

export interface LogsResponse {
  logs: LogEntry[];
}

export interface UserResponse {
  message: string;
  user: User | User[];
}
