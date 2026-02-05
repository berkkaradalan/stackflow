import { api } from "../api-client";

export interface AgentConfig {
  temperature: number;
  max_tokens: number;
  top_p: number;
  frequency_penalty: number;
  presence_penalty: number;
}

export interface Agent {
  id: number;
  name: string;
  description: string;
  project_id: number;
  created_by: number;
  role: string;
  level: string;
  provider: string;
  model: string;
  config: AgentConfig;
  status: string;
  is_active: boolean;
  last_active_at?: string;
  total_tokens_used: number;
  total_cost: number;
  total_requests: number;
  created_at: string;
  updated_at: string;
}

export interface CreateAgentRequest {
  name: string;
  description?: string;
  project_id: number;
  role: "backend_developer" | "frontend_developer" | "fullstack_developer" | "tester" | "devops" | "project_manager";
  level: "junior" | "mid" | "senior";
  provider: "openrouter" | "openai" | "anthropic" | "gemini" | "groq" | "together" | "glm" | "claude" | "kimi";
  model: string;
  api_key: string;
  config?: AgentConfig;
}

export interface UpdateAgentRequest {
  name?: string;
  description?: string;
  role?: "backend_developer" | "frontend_developer" | "fullstack_developer" | "tester" | "devops" | "project_manager";
  level?: "junior" | "mid" | "senior";
  provider?: "openrouter" | "openai" | "anthropic" | "gemini" | "groq" | "together" | "glm" | "claude" | "kimi";
  model?: string;
  api_key?: string;
  config?: AgentConfig;
  status?: "idle" | "active" | "busy" | "error" | "disabled" | "initializing";
  is_active?: boolean;
}

export interface AgentListResponse {
  agents: Agent[];
  total_count: number;
}

export interface AgentStatusResponse {
  status: string;
  is_active: boolean;
  last_active_at?: string;
}

export interface AgentWorkloadResponse {
  total_requests: number;
  total_tokens_used: number;
  total_cost: number;
  last_active_at?: string;
  updated_at: string;
  description: string;
}

export interface AgentPerformanceResponse {
  total_requests: number;
  total_tokens_used: number;
  total_cost: number;
  average_tokens_per_request: number;
  average_cost_per_request: number;
  last_active_at?: string;
  updated_at: string;
  description: string;
}

export interface AgentHealthResponse {
  id: number;
  name: string;
  status: string;
  is_active: boolean;
  last_active_at?: string;
  updated_at: string;
  healthy: boolean;
  message: string;
}

export const agentsApi = {
  /**
   * Get all agents
   */
  getAll: async (): Promise<AgentListResponse> => {
    return api.get<AgentListResponse>("/api/agents");
  },

  /**
   * Get a single agent by ID
   */
  getById: async (id: number): Promise<Agent> => {
    return api.get<Agent>(`/api/agents/${id}`);
  },

  /**
   * Get agents by project ID
   */
  getByProjectId: async (projectId: number): Promise<AgentListResponse> => {
    return api.get<AgentListResponse>(`/api/projects/${projectId}/agents`);
  },

  /**
   * Create a new agent
   */
  create: async (data: CreateAgentRequest): Promise<Agent> => {
    return api.post<Agent>("/api/agents", data);
  },

  /**
   * Update an existing agent
   */
  update: async (id: number, data: UpdateAgentRequest): Promise<Agent> => {
    return api.put<Agent>(`/api/agents/${id}`, data);
  },

  /**
   * Delete an agent
   */
  delete: async (id: number): Promise<{ message: string }> => {
    return api.delete<{ message: string }>(`/api/agents/${id}`);
  },

  /**
   * Get agent status
   */
  getStatus: async (id: number): Promise<AgentStatusResponse> => {
    return api.get<AgentStatusResponse>(`/api/agents/${id}/status`);
  },

  /**
   * Get agent workload metrics
   */
  getWorkload: async (id: number): Promise<AgentWorkloadResponse> => {
    return api.get<AgentWorkloadResponse>(`/api/agents/${id}/workload`);
  },

  /**
   * Get agent performance metrics
   */
  getPerformance: async (id: number): Promise<AgentPerformanceResponse> => {
    return api.get<AgentPerformanceResponse>(`/api/agents/${id}/performance`);
  },

  /**
   * Perform health check on agent
   */
  healthCheck: async (id: number): Promise<AgentHealthResponse> => {
    return api.get<AgentHealthResponse>(`/api/agents/${id}/health`);
  },
};
