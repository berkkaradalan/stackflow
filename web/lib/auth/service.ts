import { api } from "../api-client";
import type { LoginCredentials, LoginResponse, User } from "./types";
import { setAuthData, clearAuthData } from "./storage";

const AUTH_ENDPOINTS = {
  LOGIN: "/api/auth/login",
  LOGOUT: "/api/auth/logout",
  REFRESH: "/api/auth/refresh",
  ME: "/api/auth/me",
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
