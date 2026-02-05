import type { AuthTokens, User } from "./types";

const STORAGE_KEYS = {
  ACCESS_TOKEN: "access_token",
  REFRESH_TOKEN: "refresh_token",
  USER: "user",
} as const;

// Safe localStorage access (handles SSR)
const storage = {
  get: (key: string): string | null => {
    if (typeof window === "undefined") return null;
    try {
      return localStorage.getItem(key);
    } catch {
      return null;
    }
  },
  set: (key: string, value: string): void => {
    if (typeof window === "undefined") return;
    try {
      localStorage.setItem(key, value);
    } catch {
      // Ignore storage errors (e.g., private mode)
    }
  },
  remove: (key: string): void => {
    if (typeof window === "undefined") return;
    try {
      localStorage.removeItem(key);
    } catch {
      // Ignore storage errors
    }
  },
};

export function getAccessToken(): string | null {
  return storage.get(STORAGE_KEYS.ACCESS_TOKEN);
}

export function getRefreshToken(): string | null {
  return storage.get(STORAGE_KEYS.REFRESH_TOKEN);
}

export function getStoredUser(): User | null {
  const userJson = storage.get(STORAGE_KEYS.USER);
  if (!userJson) return null;
  try {
    return JSON.parse(userJson) as User;
  } catch {
    return null;
  }
}

export function setAuthData(tokens: AuthTokens, user: User): void {
  storage.set(STORAGE_KEYS.ACCESS_TOKEN, tokens.access_token);
  storage.set(STORAGE_KEYS.REFRESH_TOKEN, tokens.refresh_token);
  storage.set(STORAGE_KEYS.USER, JSON.stringify(user));
}

export function clearAuthData(): void {
  storage.remove(STORAGE_KEYS.ACCESS_TOKEN);
  storage.remove(STORAGE_KEYS.REFRESH_TOKEN);
  storage.remove(STORAGE_KEYS.USER);
}

export function hasAuthTokens(): boolean {
  return !!getAccessToken();
}
