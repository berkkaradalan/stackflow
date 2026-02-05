import { api } from "../api-client";
import type { LoginCredentials, LoginResponse, User } from "./types";
import { setAuthData, clearAuthData } from "./storage";

const AUTH_ENDPOINTS = {
  LOGIN: "/api/auth/login",
  LOGOUT: "/api/auth/logout",
  REFRESH: "/api/auth/refresh",
  ME: "/api/auth/me",
  PROFILE: "/api/auth/profile",
} as const;

/**
 * Login with email and password
 */
export async function login(credentials: LoginCredentials): Promise<User> {
  const response = await api.post<LoginResponse>(
    AUTH_ENDPOINTS.LOGIN,
    credentials
  );
  
  // Store auth data
  setAuthData(
    { access_token: response.access_token, refresh_token: response.refresh_token },
    response.user
  );
  
  return response.user;
}

/**
 * Logout the current user
 */
export async function logout(): Promise<void> {
  try {
    // Call logout endpoint if needed
    await api.post(AUTH_ENDPOINTS.LOGOUT, {});
  } catch {
    // Ignore errors, clear local data anyway
  } finally {
    clearAuthData();
  }
}

interface ProfileResponse {
  user: User;
}

export interface UpdateProfileData {
  username?: string;
  email?: string;
  avatar_url?: string;
  old_password?: string;
  new_password?: string;
}

/**
 * Get current user profile
 */
export async function getCurrentUser(): Promise<User | null> {
  try {
    const user = await api.get<User>(AUTH_ENDPOINTS.ME);
    return user;
  } catch {
    return null;
  }
}

/**
 * Get user profile from /api/auth/profile endpoint
 * Response format: { user: { id, username, email, avatar_url, role, is_active, created_at, updated_at } }
 */
export async function getUserProfile(): Promise<User | null> {
  try {
    const response = await api.get<ProfileResponse>(AUTH_ENDPOINTS.PROFILE);
    return response.user;
  } catch {
    return null;
  }
}

/**
 * Update user profile
 * PUT /api/auth/profile
 */
export async function updateProfile(data: UpdateProfileData): Promise<User> {
  const response = await api.put<ProfileResponse>(AUTH_ENDPOINTS.PROFILE, data);
  return response.user;
}

/**
 * Refresh access token
 */
export async function refreshToken(): Promise<boolean> {
  try {
    const refreshToken = typeof window !== "undefined" 
      ? localStorage.getItem("refresh_token") 
      : null;
      
    if (!refreshToken) return false;
    
    const response = await api.post<LoginResponse>(AUTH_ENDPOINTS.REFRESH, {
      refresh_token: refreshToken,
    });
    
    setAuthData(
      { access_token: response.access_token, refresh_token: response.refresh_token },
      response.user
    );
    
    return true;
  } catch {
    clearAuthData();
    return false;
  }
}

export { AUTH_ENDPOINTS };
