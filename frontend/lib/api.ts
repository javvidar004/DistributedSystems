import axios from 'axios';
import {
        AuthResponse,
        CreateUserRequest,
        LogEntry,
        LoginRequest,
        RegisterRequest,
        UpdatePasswordRequest,
        UpdateUserRequest,
        User,
        UserResponse,
} from '../types';

const API_BASE_URL = (() => {
    return process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
})();

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

type BackendUserResponse = AuthResponse | { user: User };
type BackendUsersResponse = {users: User[]} | UserResponse;
type BackendLogsResponse = { logs: LogEntry[] };

const getStoredUserId = () => {
    if (typeof document === 'undefined') {
        return null;
    }

    const match = document.cookie.match('(?:^|; )userId=([^;]+)');
    return match && match[1] ? decodeURIComponent(match[1]) : null;
};

// Helper to clear client-side auth state
export const clearAuth = () => {
    try {
        if (apiClient?.defaults?.headers) {
            delete apiClient.defaults.headers.common['Authorization'];
        }
        if (typeof document !== 'undefined') {
            document.cookie = 'userId=; Path=/; Max-Age=0; SameSite=Lax';
        }
        if (typeof localStorage !== 'undefined') {
            localStorage.removeItem('userId');
        }
    } catch (e) {
        console.warn('clearAuth failed', e);
    }
};


export const userSignIn = async (
    credentials: LoginRequest
): Promise<AuthResponse> => {
    try {
        const response = await apiClient.post<BackendUserResponse>('/auth/login', credentials);
        const user = response.data.user;

        if (user?.id) {
            try {
                const maxAge = 60 * 60 * 24 * 7;
                const secure = window.location.protocol === 'https:' ? '; Secure' : '';
                document.cookie = `userId=${encodeURIComponent(String(user.id))}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`;
                localStorage.setItem('userId', String(user.id));
            } catch (e) {
                console.warn('Failed to set userId cookie', e);
            }
        }

        return {
            message: response.data.message ?? 'login successful',
            user,
        };
    } catch (error) {
        console.error('Error signing in:', error);
        throw error;
    }
};
export const registerUser = async (payload: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>('/auth/register', {
        name: payload.name,
        last_name: payload.lastName,
        email: payload.email,
        work_position: payload.workPosition,
        salary: payload.salary,
        password: payload.password,
    });
    return response.data;
};

export const getUserData = async (): Promise<User> => {
    const userId = getStoredUserId();

    if (!userId) {
        throw new Error('userId cookie not found');
    }

    const response = await apiClient.get<UserResponse>(`/users/user/${userId}`);
    const user = response.data.user;

    if (Array.isArray(user)) {
        throw new Error('expected a single user record');
    }

    return user;
};

export const getUsers = async (): Promise<User[]> => {
    const response = await apiClient.get<BackendUsersResponse>('/users/users');

    if (Array.isArray(response.data.users)) {
        return response.data.users;
    }

    if (response.data.user) {
        return [response.data.user];
    }

    return [];
};

export const createUser = async (payload: CreateUserRequest): Promise<User> => {
    const response = await apiClient.post<UserResponse>('/auth/register', {
        email: payload.email,
        name: payload.name,
        last_name: payload.lastName,
        work_position: payload.workPosition,
        salary: payload.salary,
        password: payload.password,
    });
    const user = response.data.user;

    if (Array.isArray(user)) {
        throw new Error('expected a single user record');
    }

    return user;
};

export const updateUserProfile = async (payload: UpdateUserRequest): Promise<User> => {
    const userId = getStoredUserId();

    if (!userId) {
        throw new Error('userId cookie not found');
    }

    const response = await apiClient.put<UserResponse>(`/users/user/${userId}`, {
        email: payload.email,
        name: payload.name,
        last_name: payload.lastName,
        work_position: payload.workPosition,
        salary: payload.salary,
    });
    const user = response.data.user;

    if (Array.isArray(user)) {
        throw new Error('expected a single user record');
    }

    return user;
};

export const updatePassword = async (formData: UpdatePasswordRequest): Promise<{ message: string }> => {
    const response = await apiClient.put<{ message: string }>('/auth/update', {
        email: formData.email,
        old_password: formData.currentPassword,
        new_password: formData.newPassword,
    });

    return response.data;
};

export const getLogs = async (): Promise<LogEntry[]> => {
    const response = await apiClient.get<BackendLogsResponse>('/logs/log');
    return response.data.logs ?? [];
};