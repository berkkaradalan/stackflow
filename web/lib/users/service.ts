import { api } from "@/lib/api-client";
import type {
  User,
  InviteUserData,
  InviteUserResponse,
  UpdateUserData,
  GetAllUsersResponse,
} from "./types";

/**
 * Get all users (admin only)
 */
export async function getAllUsers(): Promise<User[]> {
  try {
    const response = await api.get<GetAllUsersResponse>("/api/users");
    return response.users || [];
  } catch (error) {
    console.error("Failed to fetch users:", error);
    throw error;
  }
}

/**
 * Get user by ID (admin only)
 */
export async function getUserById(id: number): Promise<User> {
  try {
    const user = await api.get<User>(`/api/users/${id}`);
    return user;
  } catch (error) {
    console.error(`Failed to fetch user ${id}:`, error);
    throw error;
  }
}

/**
 * Update user (admin only)
 */
export async function updateUser(
  id: number,
  data: UpdateUserData
): Promise<User> {
  try {
    const user = await api.put<User>(`/api/users/${id}`, data);
    return user;
  } catch (error) {
    console.error(`Failed to update user ${id}:`, error);
    throw error;
  }
}

/**
 * Delete user (admin only)
 */
export async function deleteUser(id: number): Promise<void> {
  try {
    await api.delete(`/api/users/${id}`);
  } catch (error) {
    console.error(`Failed to delete user ${id}:`, error);
    throw error;
  }
}

/**
 * Invite user (admin only)
 */
export async function inviteUser(
  data: InviteUserData
): Promise<InviteUserResponse> {
  try {
    const response = await api.post<InviteUserResponse>(
      "/api/users/invite",
      data
    );
    return response;
  } catch (error) {
    console.error("Failed to invite user:", error);
    throw error;
  }
}
