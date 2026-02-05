import { getAccessToken } from "./auth/storage";

// Base API URL - configure via environment variable
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Default request headers
const defaultHeaders: Record<string, string> = {
  "Content-Type": "application/json",
};

export interface ApiError {
  error: string;
}

export class ApiRequestError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: ApiError
  ) {
    super(message);
    this.name = "ApiRequestError";
  }
}

/**
 * Make an API request with automatic JSON parsing and error handling
 */
export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // Get auth token if available
  const token = typeof window !== "undefined" ? getAccessToken() : null;
  
  const headers: Record<string, string> = {
    ...defaultHeaders,
    ...((options.headers as Record<string, string>) || {}),
  };
  
  // Add auth header if token exists
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  
  const config: RequestInit = {
    ...options,
    headers,
  };
  
  const response = await fetch(url, config);
  
  // Parse JSON response
  let data: unknown;
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.includes("application/json")) {
    data = await response.json();
  } else {
    data = await response.text();
  }
  
  // Handle error responses
  if (!response.ok) {
    const errorMessage = 
      typeof data === "object" && data !== null && "error" in data
        ? String((data as ApiError).error)
        : response.statusText;
    
    throw new ApiRequestError(
      errorMessage,
      response.status,
      typeof data === "object" ? (data as ApiError) : undefined
    );
  }
  
  return data as T;
}

/**
 * HTTP methods shortcuts
 */
export const api = {
  get: <T>(endpoint: string, options?: RequestInit) =>
    apiRequest<T>(endpoint, { ...options, method: "GET" }),
    
  post: <T>(endpoint: string, body: unknown, options?: RequestInit) =>
    apiRequest<T>(endpoint, {
      ...options,
      method: "POST",
      body: JSON.stringify(body),
    }),
    
  put: <T>(endpoint: string, body: unknown, options?: RequestInit) =>
    apiRequest<T>(endpoint, {
      ...options,
      method: "PUT",
      body: JSON.stringify(body),
    }),
    
  patch: <T>(endpoint: string, body: unknown, options?: RequestInit) =>
    apiRequest<T>(endpoint, {
      ...options,
      method: "PATCH",
      body: JSON.stringify(body),
    }),
    
  delete: <T>(endpoint: string, options?: RequestInit) =>
    apiRequest<T>(endpoint, { ...options, method: "DELETE" }),
};

export { API_BASE_URL };
