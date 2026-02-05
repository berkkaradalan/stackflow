import { api } from "../api-client";

export interface Project {
  id: number;
  name: string;
  description: string;
  status: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
  status?: "active" | "inactive" | "archived";
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  status?: "active" | "inactive" | "archived";
}

export interface ProjectListResponse {
  projects: Project[];
  total_count: number;
}

export interface ProjectStats {
  total_tasks: number;
  completed_tasks: number;
  pending_tasks: number;
  total_agents: number;
  total_workflows: number;
}

export const projectsApi = {
  /**
   * Get all projects
   */
  getAll: async (): Promise<ProjectListResponse> => {
    return api.get<ProjectListResponse>("/api/projects");
  },

  /**
   * Get a single project by ID
   */
  getById: async (id: number): Promise<Project> => {
    return api.get<Project>(`/api/projects/${id}`);
  },

  /**
   * Create a new project
   */
  create: async (data: CreateProjectRequest): Promise<Project> => {
    return api.post<Project>("/api/projects", data);
  },

  /**
   * Update an existing project
   */
  update: async (id: number, data: UpdateProjectRequest): Promise<Project> => {
    return api.put<Project>(`/api/projects/${id}`, data);
  },

  /**
   * Delete a project
   */
  delete: async (id: number): Promise<{ message: string }> => {
    return api.delete<{ message: string }>(`/api/projects/${id}`);
  },

  /**
   * Get project statistics
   */
  getStats: async (id: number): Promise<ProjectStats> => {
    return api.get<ProjectStats>(`/api/projects/${id}/stats`);
  },
};
