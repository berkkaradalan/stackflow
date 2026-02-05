import { useState, useCallback } from "react";
import { agentsApi, type Agent, type CreateAgentRequest, type UpdateAgentRequest, type AgentHealthResponse, type AgentWorkloadResponse, type AgentPerformanceResponse } from "@/lib/api/agents";

export function useAgents() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAgents = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await agentsApi.getAll();
      setAgents(response.agents || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch agents");
      setAgents([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchAgentsByProject = useCallback(async (projectId: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await agentsApi.getByProjectId(projectId);
      setAgents(response.agents || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch agents");
      setAgents([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchAgentById = useCallback(async (id: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const agent = await agentsApi.getById(id);
      return agent;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch agent");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createAgent = useCallback(async (data: CreateAgentRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      const newAgent = await agentsApi.create(data);
      setAgents((prev) => [newAgent, ...(prev || [])]);
      return newAgent;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create agent";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateAgent = useCallback(async (id: number, data: UpdateAgentRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedAgent = await agentsApi.update(id, data);
      setAgents((prev) =>
        (prev || []).map((agent) => (agent.id === id ? updatedAgent : agent))
      );
      return updatedAgent;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update agent";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteAgent = useCallback(async (id: number) => {
    setIsLoading(true);
    setError(null);
    try {
      await agentsApi.delete(id);
      setAgents((prev) => (prev || []).filter((agent) => agent.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete agent";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const healthCheck = useCallback(async (id: number): Promise<AgentHealthResponse> => {
    setIsLoading(true);
    setError(null);
    try {
      const health = await agentsApi.healthCheck(id);
      return health;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to perform health check";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const getWorkload = useCallback(async (id: number): Promise<AgentWorkloadResponse> => {
    setIsLoading(true);
    setError(null);
    try {
      const workload = await agentsApi.getWorkload(id);
      return workload;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to fetch workload";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const getPerformance = useCallback(async (id: number): Promise<AgentPerformanceResponse> => {
    setIsLoading(true);
    setError(null);
    try {
      const performance = await agentsApi.getPerformance(id);
      return performance;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to fetch performance";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    agents,
    isLoading,
    error,
    fetchAgents,
    fetchAgentsByProject,
    fetchAgentById,
    createAgent,
    updateAgent,
    deleteAgent,
    healthCheck,
    getWorkload,
    getPerformance,
  };
}
