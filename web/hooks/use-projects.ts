import { useState, useCallback } from "react";
import { projectsApi, type Project, type CreateProjectRequest, type UpdateProjectRequest } from "@/lib/api/projects";

export function useProjects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProjects = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await projectsApi.getAll();
      setProjects(response.projects || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch projects");
      setProjects([]); // Reset to empty array on error
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createProject = useCallback(async (data: CreateProjectRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      const newProject = await projectsApi.create(data);
      setProjects((prev) => [newProject, ...(prev || [])]);
      return newProject;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create project";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateProject = useCallback(async (id: number, data: UpdateProjectRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedProject = await projectsApi.update(id, data);
      setProjects((prev) =>
        (prev || []).map((project) => (project.id === id ? updatedProject : project))
      );
      return updatedProject;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update project";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteProject = useCallback(async (id: number) => {
    setIsLoading(true);
    setError(null);
    try {
      await projectsApi.delete(id);
      setProjects((prev) => (prev || []).filter((project) => project.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete project";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    projects,
    isLoading,
    error,
    fetchProjects,
    createProject,
    updateProject,
    deleteProject,
  };
}
