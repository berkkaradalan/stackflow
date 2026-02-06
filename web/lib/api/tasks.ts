import { api } from "../api-client";

// Task status types
export type TaskStatus = "open" | "in_progress" | "done" | "closed" | "wont_do";
export type TaskPriority = "low" | "medium" | "high" | "critical";
export type CreatorType = "user" | "agent";

export interface Task {
  id: number;
  project_id: number;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  assigned_agent_id: number | null;
  reviewer_id: number | null;
  created_by: number;
  creator_type: CreatorType;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface TaskWithDetails extends Task {
  assigned_agent_name: string | null;
  reviewer_name: string | null;
  creator_name: string;
  project_name: string;
}

export interface TaskActivity {
  id: number;
  task_id: number;
  actor_id: number;
  actor_type: CreatorType;
  action: string;
  old_value: string | null;
  new_value: string | null;
  message: string;
  created_at: string;
}

export interface TaskActivityWithDetails extends TaskActivity {
  actor_name: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
  priority?: TaskPriority;
  assigned_agent_id?: number;
  reviewer_id?: number;
  tags?: string[];
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  priority?: TaskPriority;
  tags?: string[];
}

export interface AssignAgentRequest {
  agent_id: number;
}

export interface SetReviewerRequest {
  reviewer_id: number;
}

export interface StatusChangeRequest {
  message?: string;
}

export interface CreateActivityRequest {
  message: string;
}

export interface TaskListResponse {
  tasks: TaskWithDetails[];
  total_count: number;
}

export interface TaskActivityListResponse {
  activities: TaskActivityWithDetails[];
  total_count: number;
}

export interface TaskFilters {
  project_id?: number;
  status?: TaskStatus;
  priority?: TaskPriority;
  assigned_agent_id?: number;
  reviewer_id?: number;
}

export const tasksApi = {
  // Get all tasks with optional filters
  getAll: async (filters?: TaskFilters): Promise<TaskListResponse> => {
    const params = new URLSearchParams();
    if (filters?.project_id) params.append("project_id", String(filters.project_id));
    if (filters?.status) params.append("status", filters.status);
    if (filters?.priority) params.append("priority", filters.priority);
    if (filters?.assigned_agent_id) params.append("assigned_agent_id", String(filters.assigned_agent_id));
    if (filters?.reviewer_id) params.append("reviewer_id", String(filters.reviewer_id));

    const queryString = params.toString();
    return api.get<TaskListResponse>(`/api/tasks${queryString ? `?${queryString}` : ""}`);
  },

  // Get all tasks for a project
  getByProject: async (projectId: number): Promise<TaskListResponse> => {
    return api.get<TaskListResponse>(`/api/projects/${projectId}/tasks`);
  },

  // Create a new task
  create: async (projectId: number, data: CreateTaskRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/projects/${projectId}/tasks`, data);
  },

  // Get a single task by ID
  getById: async (id: number): Promise<TaskWithDetails> => {
    return api.get<TaskWithDetails>(`/api/tasks/${id}`);
  },

  // Update a task
  update: async (id: number, data: UpdateTaskRequest): Promise<TaskWithDetails> => {
    return api.put<TaskWithDetails>(`/api/tasks/${id}`, data);
  },

  // Delete a task
  delete: async (id: number): Promise<{ message: string }> => {
    return api.delete<{ message: string }>(`/api/tasks/${id}`);
  },

  // Assign an agent to a task
  assignAgent: async (id: number, data: AssignAgentRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/assign`, data);
  },

  // Set reviewer for a task
  setReviewer: async (id: number, data: SetReviewerRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/reviewer`, data);
  },

  // Status transitions
  start: async (id: number, data?: StatusChangeRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/start`, data || {});
  },

  complete: async (id: number, data?: StatusChangeRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/done`, data || {});
  },

  close: async (id: number, data?: StatusChangeRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/close`, data || {});
  },

  wontDo: async (id: number, data?: StatusChangeRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/wontdo`, data || {});
  },

  reopen: async (id: number, data?: StatusChangeRequest): Promise<TaskWithDetails> => {
    return api.post<TaskWithDetails>(`/api/tasks/${id}/reopen`, data || {});
  },

  // Activities
  getActivities: async (id: number): Promise<TaskActivityListResponse> => {
    return api.get<TaskActivityListResponse>(`/api/tasks/${id}/activities`);
  },

  addProgress: async (id: number, data: CreateActivityRequest): Promise<TaskActivityWithDetails> => {
    return api.post<TaskActivityWithDetails>(`/api/tasks/${id}/activities`, data);
  },
};

// Helper functions
export const getStatusLabel = (status: TaskStatus): string => {
  const labels: Record<TaskStatus, string> = {
    open: "Open",
    in_progress: "In Progress",
    done: "Done",
    closed: "Closed",
    wont_do: "Won't Do",
  };
  return labels[status] || status;
};

export const getStatusColor = (status: TaskStatus): string => {
  const colors: Record<TaskStatus, string> = {
    open: "bg-blue-100 text-blue-800",
    in_progress: "bg-yellow-100 text-yellow-800",
    done: "bg-green-100 text-green-800",
    closed: "bg-gray-100 text-gray-800",
    wont_do: "bg-red-100 text-red-800",
  };
  return colors[status] || "bg-gray-100 text-gray-800";
};

export const getPriorityColor = (priority: TaskPriority): string => {
  const colors: Record<TaskPriority, string> = {
    low: "bg-slate-100 text-slate-800",
    medium: "bg-blue-100 text-blue-800",
    high: "bg-orange-100 text-orange-800",
    critical: "bg-red-100 text-red-800",
  };
  return colors[priority] || "bg-gray-100 text-gray-800";
};

export const getPriorityLabel = (priority: TaskPriority): string => {
  const labels: Record<TaskPriority, string> = {
    low: "Low",
    medium: "Medium",
    high: "High",
    critical: "Critical",
  };
  return labels[priority] || priority;
};
