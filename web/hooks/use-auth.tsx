"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import type { User, LoginCredentials } from "@/lib/auth/types";
import { login as loginService, logout as logoutService, getUserProfile, updateProfile as updateProfileService } from "@/lib/auth/service";
import type { UpdateProfileData } from "@/lib/auth/service";
import { getStoredUser, hasAuthTokens, clearAuthData, setAuthData } from "@/lib/auth/storage";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  updateProfile: (data: UpdateProfileData) => Promise<void>;
  error: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  // Initialize auth state from storage and fetch fresh profile
  useEffect(() => {
    const initializeAuth = async () => {
      const storedUser = getStoredUser();
      
      if (hasAuthTokens()) {
        // Fetch fresh profile data from API
        const freshUser = await getUserProfile();
        if (freshUser) {
          setUser(freshUser);
          // Update stored user data
          const tokens = {
            access_token: localStorage.getItem("access_token") || "",
            refresh_token: localStorage.getItem("refresh_token") || "",
          };
          setAuthData(tokens, freshUser);
        } else if (storedUser) {
          // Fallback to stored user if API fails
          setUser(storedUser);
        }
      }
      
      setIsLoading(false);
    };
    
    initializeAuth();
  }, []);

  const login = useCallback(async (credentials: LoginCredentials) => {
    setError(null);
    setIsLoading(true);
    
    try {
      const user = await loginService(credentials);
      setUser(user);
      router.push("/"); // Redirect to dashboard after login
    } catch (err) {
      const message = err instanceof Error ? err.message : "Login failed";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [router]);

  const logout = useCallback(async () => {
    setIsLoading(true);
    
    try {
      await logoutService();
    } finally {
      setUser(null);
      clearAuthData();
      setIsLoading(false);
      router.push("/login"); // Redirect to login after logout
    }
  }, [router]);

  const updateProfile = useCallback(async (data: UpdateProfileData) => {
    setError(null);
    setIsLoading(true);
    
    try {
      const updatedUser = await updateProfileService(data);
      setUser(updatedUser);
      // Update stored user data
      const tokens = {
        access_token: localStorage.getItem("access_token") || "",
        refresh_token: localStorage.getItem("refresh_token") || "",
      };
      setAuthData(tokens, updatedUser);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to update profile";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    logout,
    updateProfile,
    error,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

export function useRequireAuth() {
  const { user, isLoading, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  return { user, isLoading, isAuthenticated };
}

export function useRedirectIfAuthenticated() {
  const { isLoading, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push("/");
    }
  }, [isLoading, isAuthenticated, router]);
}
