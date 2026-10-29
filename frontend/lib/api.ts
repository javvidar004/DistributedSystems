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

const getStoredAccessToken = () => {
    if (typeof document === 'undefined') {
        return null;
    }

    const match = document.cookie.match('(?:^|; )accessToken=([^;]+)');
    return match && match[1] ? decodeURIComponent(match[1]) : null;
};

const extractUserIdFromToken = (token: string): string | null => {
    try {
        const payloadPart = token.split('.')[1];
        if (!payloadPart) {
            return null;
        }

        const normalized = payloadPart.replace(/-/g, '+').replace(/_/g, '/');
        const padded = normalized + '='.repeat((4 - (normalized.length % 4)) % 4);
        const decoded = atob(padded);
        const payload = JSON.parse(decoded) as { user_id?: string; sub?: string };

        return payload.user_id ?? payload.sub ?? null;
    } catch {
        return null;
    }
};

// Helper to clear client-side auth state
export const clearAuth = () => {
    try {
        if (apiClient?.defaults?.headers) {
            delete apiClient.defaults.headers.common['Authorization'];
        }
        if (typeof document !== 'undefined') {
            document.cookie = 'accessToken=; Path=/; Max-Age=0; SameSite=Lax';
            document.cookie = 'refreshToken=; Path=/; Max-Age=0; SameSite=Lax';
        }
        if (typeof localStorage !== 'undefined') {
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
        }
    } catch (e) {
        console.warn('clearAuth failed', e);
    }
};


export const userSignIn = async (
    credentials: LoginRequest
): Promise<AuthResponse> => {
    try {
        const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
        const user = response.data.user;
        const accessToken = response.data.access_token;
        const refreshToken = response.data.refresh_token;

        if (accessToken) {
            try {
                const maxAge = 60 * 60 * 24 * 7;
                const secure = window.location.protocol === 'https:' ? '; Secure' : '';
                document.cookie = `accessToken=${encodeURIComponent(accessToken)}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`;
                localStorage.setItem('accessToken', accessToken);
                apiClient.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
            } catch (e) {
                console.warn('Failed to store access token', e);
            }
        }

        if (refreshToken) {
            try {
                const maxAge = 60 * 60 * 24 * 30;
                const secure = window.location.protocol === 'https:' ? '; Secure' : '';
                document.cookie = `refreshToken=${encodeURIComponent(refreshToken)}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`;
                localStorage.setItem('refreshToken', refreshToken);
            } catch (e) {
                console.warn('Failed to store refresh token', e);
            }
        }

        return {
            message: response.data.message ?? 'login successful',
            user,
            access_token: accessToken,
            refresh_token: refreshToken,
        };
    } catch (error) {
        console.error('Error signing in:', error);
        throw error;
    }
};

export const validateJWT = async (token?: string): Promise<{ user_id: string; role: string; group: string }> => {
    const bearerToken = token ?? getStoredAccessToken();
    if (!bearerToken) {
        throw new Error('access token not found');
    }

    const response = await apiClient.post<{ user_id: string; role: string; group: string }>(
        '/auth/validate',
        { token: bearerToken },
        {
            headers: {
                Authorization: `Bearer ${bearerToken}`,
            },
        }
    );

    return response.data;
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
    const token = getStoredAccessToken();
    if (!token) {
        throw new Error('access token not found');
    }

    const userId = extractUserIdFromToken(token);

    if (!userId) {
        throw new Error('invalid access token payload');
    }

    const response = await apiClient.get<UserResponse>(`/users/user/${userId}`, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
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
    const token = getStoredAccessToken();
    if (!token) {
        throw new Error('access token not found');
    }

    const userId = extractUserIdFromToken(token);

    if (!userId) {
        throw new Error('invalid access token payload');
    }

    const response = await apiClient.put<UserResponse>(`/users/user/${userId}`, {
        email: payload.email,
        name: payload.name,
        last_name: payload.lastName,
        work_position: payload.workPosition,
        salary: payload.salary,
    }, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
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