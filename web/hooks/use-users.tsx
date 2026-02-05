"use client";

import { useState, useCallback } from "react";
import type { User, InviteUserData, UpdateUserData } from "@/lib/users/types";
import {
  getAllUsers as getAllUsersService,
  getUserById as getUserByIdService,
  updateUser as updateUserService,
  deleteUser as deleteUserService,
  inviteUser as inviteUserService,
} from "@/lib/users/service";

export function useUsers() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const fetchedUsers = await getAllUsersService();
      setUsers(fetchedUsers);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to fetch users";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const getUserById = useCallback(async (id: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const user = await getUserByIdService(id);
      return user;
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to fetch user";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateUser = useCallback(async (id: number, data: UpdateUserData) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedUser = await updateUserService(id, data);
      setUsers((prev) =>
        prev.map((user) => (user.id === id ? updatedUser : user))
      );
      return updatedUser;
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to update user";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteUser = useCallback(async (id: number) => {
    setIsLoading(true);
    setError(null);
    try {
      await deleteUserService(id);
      setUsers((prev) => prev.filter((user) => user.id !== id));
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to delete user";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const inviteUser = useCallback(async (data: InviteUserData) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await inviteUserService(data);
      // Note: Invited user will appear in the list after they complete registration
      return response;
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to invite user";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    users,
    isLoading,
    error,
    fetchUsers,
    getUserById,
    updateUser,
    deleteUser,
    inviteUser,
  };
}
