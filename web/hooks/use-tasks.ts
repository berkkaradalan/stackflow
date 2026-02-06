import { useState, useCallback } from "react";
import {
  tasksApi,
  type TaskWithDetails,
  type TaskActivityWithDetails,
  type CreateTaskRequest,
  type UpdateTaskRequest,
  type TaskStatus,
} from "@/lib/api/tasks";

export function useTasks(projectId: number) {
  const [tasks, setTasks] = useState<TaskWithDetails[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchTasks = useCallback(async () => {
    if (!projectId) return;
    setIsLoading(true);
    setError(null);
    try {
      const response = await tasksApi.getByProject(projectId);
      setTasks(response.tasks || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch tasks");
      setTasks([]);
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);

  const createTask = useCallback(
    async (data: CreateTaskRequest) => {
      setIsLoading(true);
      setError(null);
      try {
        const newTask = await tasksApi.create(projectId, data);
        setTasks((prev) => [newTask, ...prev]);
        return newTask;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Failed to create task";
        setError(errorMessage);
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [projectId]
  );

  const updateTask = useCallback(async (id: number, data: UpdateTaskRequest) => {
    setError(null);
    try {
      const updatedTask = await tasksApi.update(id, data);
      setTasks((prev) => prev.map((task) => (task.id === id ? updatedTask : task)));
      return updatedTask;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update task";
      setError(errorMessage);
      throw err;
    }
  }, []);

  const deleteTask = useCallback(async (id: number) => {
    setError(null);
    try {
      await tasksApi.delete(id);
      setTasks((prev) => prev.filter((task) => task.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete task";
      setError(errorMessage);
      throw err;
    }
  }, []);

  const assignAgent = useCallback(async (taskId: number, agentId: number) => {
    setError(null);
    try {
      const updatedTask = await tasksApi.assignAgent(taskId, { agent_id: agentId });
      setTasks((prev) => prev.map((task) => (task.id === taskId ? updatedTask : task)));
      return updatedTask;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to assign agent";
      setError(errorMessage);
      throw err;
    }
  }, []);

  const setReviewer = useCallback(async (taskId: number, reviewerId: number) => {
    setError(null);
    try {
      const updatedTask = await tasksApi.setReviewer(taskId, { reviewer_id: reviewerId });
      setTasks((prev) => prev.map((task) => (task.id === taskId ? updatedTask : task)));
      return updatedTask;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to set reviewer";
      setError(errorMessage);
      throw err;
    }
  }, []);

  // Status change helper
  const changeStatus = useCallback(
    async (taskId: number, action: "start" | "complete" | "close" | "wontdo" | "reopen", message?: string) => {
      setError(null);
      try {
        const apiMethod = {
          start: tasksApi.start,
          complete: tasksApi.complete,
          close: tasksApi.close,
          wontdo: tasksApi.wontDo,
          reopen: tasksApi.reopen,
        }[action];

        const updatedTask = await apiMethod(taskId, message ? { message } : undefined);
        setTasks((prev) => prev.map((task) => (task.id === taskId ? updatedTask : task)));
        return updatedTask;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Failed to change status";
        setError(errorMessage);
        throw err;
      }
    },
    []
  );

  // Get tasks grouped by status for kanban
  const getTasksByStatus = useCallback(() => {
    const grouped: Record<TaskStatus, TaskWithDetails[]> = {
      open: [],
      in_progress: [],
      done: [],
      closed: [],
      wont_do: [],
    };

    tasks.forEach((task) => {
      if (grouped[task.status]) {
        grouped[task.status].push(task);
      }
    });

    return grouped;
  }, [tasks]);

  return {
    tasks,
    isLoading,
    error,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    assignAgent,
    setReviewer,
    changeStatus,
    getTasksByStatus,
  };
}

// Hook for task activities
export function useTaskActivities(taskId: number) {
  const [activities, setActivities] = useState<TaskActivityWithDetails[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchActivities = useCallback(async () => {
    if (!taskId) return;
    setIsLoading(true);
    setError(null);
    try {
      const response = await tasksApi.getActivities(taskId);
      setActivities(response.activities || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch activities");
      setActivities([]);
    } finally {
      setIsLoading(false);
    }
  }, [taskId]);

  const addProgress = useCallback(
    async (message: string) => {
      setError(null);
      try {
        const newActivity = await tasksApi.addProgress(taskId, { message });
        setActivities((prev) => [newActivity, ...prev]);
        return newActivity;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Failed to add progress";
        setError(errorMessage);
        throw err;
      }
    },
    [taskId]
  );

  return {
    activities,
    isLoading,
    error,
    fetchActivities,
    addProgress,
  };
}
