import type { User } from "@/lib/auth/types";

export type { User };

export interface InviteUserData {
  email: string;
  username: string;
  role?: string;
}

export interface InviteUserResponse {
  invite_link: string;
  token: string;
  expires_at: string;
  email: string;
}

export interface UpdateUserData {
  username?: string;
  email?: string;
  role?: string;
  is_active?: boolean;
  avatar_url?: string;
}

export interface GetAllUsersResponse {
  users: User[];
  total_count: number;
}
