"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import {
  ListTodo,
  Loader2,
  RefreshCw,
  Filter,
  X,
  Bot,
  User,
  Tag,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  tasksApi,
  type TaskWithDetails,
  type TaskFilters,
  type TaskStatus,
  getStatusLabel,
  getStatusColor,
  getPriorityLabel,
  getPriorityColor,
} from "@/lib/api/tasks";
import { projectsApi, type Project } from "@/lib/api/projects";
import { agentsApi, type Agent } from "@/lib/api/agents";
import { TaskDetailDialog } from "@/components/task-detail-dialog";

const KANBAN_COLUMNS: { status: TaskStatus; label: string; color: string }[] = [
  { status: "open", label: "Open", color: "bg-blue-500" },
  { status: "in_progress", label: "In Progress", color: "bg-yellow-500" },
  { status: "done", label: "Done", color: "bg-green-500" },
  { status: "closed", label: "Closed", color: "bg-gray-500" },
];

interface TaskCardProps {
  task: TaskWithDetails;
  onClick: (task: TaskWithDetails) => void;
  showProject?: boolean;
}

function TaskCard({ task, onClick, showProject = true }: TaskCardProps) {
  const router = useRouter();

  return (
    <Card
      className="cursor-pointer transition-shadow hover:shadow-md"
      onClick={() => onClick(task)}
    >
      <CardHeader className="pb-2 pt-3 px-3">
        <div className="flex items-start justify-between gap-2">
          <h4 className="text-sm font-medium leading-tight line-clamp-2">
            {task.title}
          </h4>
          <Badge
            variant="secondary"
            className={`shrink-0 text-xs ${getPriorityColor(task.priority)}`}
          >
            {getPriorityLabel(task.priority)}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="pb-3 px-3 pt-0">
        {/* Project Link */}
        {showProject && (
          <div
            className="text-xs text-blue-600 hover:underline mb-2 cursor-pointer"
            onClick={(e) => {
              e.stopPropagation();
              router.push(`/projects/${task.project_id}`);
            }}
          >
            {task.project_name}
          </div>
        )}

        {task.description && (
          <p className="text-xs text-muted-foreground line-clamp-2 mb-2">
            {task.description}
          </p>
        )}

        {/* Tags */}
        {task.tags && task.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-2">
            {task.tags.slice(0, 3).map((tag, index) => (
              <Badge
                key={index}
                variant="outline"
                className="text-xs px-1.5 py-0"
              >
                <Tag className="h-2.5 w-2.5 mr-1" />
                {tag}
              </Badge>
            ))}
            {task.tags.length > 3 && (
              <Badge variant="outline" className="text-xs px-1.5 py-0">
                +{task.tags.length - 3}
              </Badge>
            )}
          </div>
        )}

        {/* Footer */}
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <div className="flex items-center gap-2">
            {task.assigned_agent_name && (
              <div className="flex items-center gap-1" title="Assigned Agent">
                <Bot className="h-3 w-3" />
                <span className="truncate max-w-[80px]">
                  {task.assigned_agent_name}
                </span>
              </div>
            )}
            {task.reviewer_name && (
              <div className="flex items-center gap-1" title="Reviewer">
                <User className="h-3 w-3" />
                <span className="truncate max-w-[60px]">
                  {task.reviewer_name}
                </span>
              </div>
            )}
          </div>
          <div className="flex items-center gap-1" title="Created by">
            {task.creator_type === "agent" ? (
              <Bot className="h-3 w-3" />
            ) : (
              <User className="h-3 w-3" />
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default function TasksPage() {
  const [tasks, setTasks] = useState<TaskWithDetails[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filtersOpen, setFiltersOpen] = useState(true);

  // Filters
  const [filters, setFilters] = useState<TaskFilters>({});
  const [selectedTask, setSelectedTask] = useState<TaskWithDetails | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  const fetchTasks = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await tasksApi.getAll(filters);
      setTasks(response.tasks || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch tasks");
      setTasks([]);
    } finally {
      setIsLoading(false);
    }
  }, [filters]);

  const fetchFiltersData = useCallback(async () => {
    try {
      const [projectsRes, agentsRes] = await Promise.all([
        projectsApi.getAll(),
        agentsApi.getAll(),
      ]);
      setProjects(projectsRes.projects || []);
      setAgents(agentsRes.agents || []);
    } catch (err) {
      console.error("Failed to fetch filter data:", err);
    }
  }, []);

  useEffect(() => {
    fetchFiltersData();
  }, [fetchFiltersData]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const handleFilterChange = (key: keyof TaskFilters, value: string | undefined) => {
    setFilters((prev) => {
      const newFilters = { ...prev };
      if (value === undefined || value === "all") {
        delete newFilters[key];
      } else if (key === "project_id" || key === "assigned_agent_id" || key === "reviewer_id") {
        (newFilters as Record<string, unknown>)[key] = parseInt(value);
      } else {
        (newFilters as Record<string, unknown>)[key] = value;
      }
      return newFilters;
    });
  };

  const clearFilters = () => {
    setFilters({});
  };

  const hasActiveFilters = Object.keys(filters).length > 0;

  const handleTaskClick = (task: TaskWithDetails) => {
    setSelectedTask(task);
    setIsDetailOpen(true);
  };

  const handleStatusChange = async (
    taskId: number,
    action: "start" | "complete" | "close" | "wontdo" | "reopen"
  ) => {
    try {
      const apiMethod = {
        start: tasksApi.start,
        complete: tasksApi.complete,
        close: tasksApi.close,
        wontdo: tasksApi.wontDo,
        reopen: tasksApi.reopen,
      }[action];

      const updatedTask = await apiMethod(taskId);
      setTasks((prev) => prev.map((t) => (t.id === taskId ? updatedTask : t)));
      if (selectedTask?.id === taskId) {
        setSelectedTask(updatedTask);
      }
    } catch (err) {
      console.error("Failed to change status:", err);
    }
  };

  // Group tasks by status
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

  const tasksByStatus = getTasksByStatus();

  // Stats
  const stats = {
    total: tasks.length,
    open: tasksByStatus.open.length,
    inProgress: tasksByStatus.in_progress.length,
    done: tasksByStatus.done.length,
    closed: tasksByStatus.closed.length,
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Tasks</h1>
          <p className="text-muted-foreground">
            Manage and track all tasks across projects
          </p>
        </div>
        <Button variant="outline" onClick={fetchTasks} disabled={isLoading}>
          <RefreshCw className={`mr-2 h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
          Refresh
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total</CardDescription>
            <CardTitle className="text-2xl">{stats.total}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Open</CardDescription>
            <CardTitle className="text-2xl text-blue-600">{stats.open}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>In Progress</CardDescription>
            <CardTitle className="text-2xl text-yellow-600">{stats.inProgress}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Done</CardDescription>
            <CardTitle className="text-2xl text-green-600">{stats.done}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Closed</CardDescription>
            <CardTitle className="text-2xl text-gray-600">{stats.closed}</CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Collapsible Filters */}
      <Collapsible open={filtersOpen} onOpenChange={setFiltersOpen}>
        <Card>
          <CollapsibleTrigger asChild>
            <CardHeader className="pb-3 cursor-pointer hover:bg-muted/50 transition-colors">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Filter className="h-4 w-4" />
                  <CardTitle className="text-base">Filters</CardTitle>
                  {hasActiveFilters && (
                    <Badge variant="secondary" className="ml-2">
                      {Object.keys(filters).length} active
                    </Badge>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  {hasActiveFilters && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        clearFilters();
                      }}
                    >
                      <X className="mr-1 h-4 w-4" />
                      Clear
                    </Button>
                  )}
                  {filtersOpen ? (
                    <ChevronUp className="h-4 w-4 text-muted-foreground" />
                  ) : (
                    <ChevronDown className="h-4 w-4 text-muted-foreground" />
                  )}
                </div>
              </div>
            </CardHeader>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <CardContent className="pt-0">
              <div className="grid gap-4 md:grid-cols-4">
                {/* Project Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Project</label>
                  <Select
                    value={filters.project_id?.toString() || "all"}
                    onValueChange={(value) => handleFilterChange("project_id", value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="All Projects" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Projects</SelectItem>
                      {projects.map((project) => (
                        <SelectItem key={project.id} value={project.id.toString()}>
                          {project.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {/* Priority Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Priority</label>
                  <Select
                    value={filters.priority || "all"}
                    onValueChange={(value) => handleFilterChange("priority", value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="All Priorities" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Priorities</SelectItem>
                      <SelectItem value="critical">Critical</SelectItem>
                      <SelectItem value="high">High</SelectItem>
                      <SelectItem value="medium">Medium</SelectItem>
                      <SelectItem value="low">Low</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Agent Filter */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Assigned Agent</label>
                  <Select
                    value={filters.assigned_agent_id?.toString() || "all"}
                    onValueChange={(value) => handleFilterChange("assigned_agent_id", value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="All Agents" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Agents</SelectItem>
                      {agents.map((agent) => (
                        <SelectItem key={agent.id} value={agent.id.toString()}>
                          {agent.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {/* Status Filter - hidden since kanban shows all statuses, but can filter to single */}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Status</label>
                  <Select
                    value={filters.status || "all"}
                    onValueChange={(value) => handleFilterChange("status", value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="All Statuses" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Statuses</SelectItem>
                      <SelectItem value="open">Open</SelectItem>
                      <SelectItem value="in_progress">In Progress</SelectItem>
                      <SelectItem value="done">Done</SelectItem>
                      <SelectItem value="closed">Closed</SelectItem>
                      <SelectItem value="wont_do">Won&apos;t Do</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardContent>
          </CollapsibleContent>
        </Card>
      </Collapsible>

      {/* Kanban Board */}
      {isLoading ? (
        <div className="flex min-h-[400px] items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <Card>
          <CardContent className="flex min-h-[200px] flex-col items-center justify-center text-center p-6">
            <p className="text-sm text-muted-foreground mb-4">{error}</p>
            <Button variant="outline" onClick={fetchTasks}>
              <RefreshCw className="mr-2 h-4 w-4" />
              Retry
            </Button>
          </CardContent>
        </Card>
      ) : tasks.length === 0 ? (
        <Card>
          <CardContent className="flex min-h-[200px] flex-col items-center justify-center text-center p-6">
            <ListTodo className="mb-4 h-12 w-12 text-muted-foreground" />
            <h3 className="mb-2 text-lg font-semibold">No Tasks Found</h3>
            <p className="text-sm text-muted-foreground">
              {hasActiveFilters
                ? "Try adjusting your filters"
                : "Create your first task in a project"}
            </p>
          </CardContent>
        </Card>
      ) : (
        <>
          {/* Main Kanban Columns */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {KANBAN_COLUMNS.map(({ status, label, color }) => (
              <div key={status} className="space-y-3">
                {/* Column Header */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${color}`} />
                    <h3 className="font-medium text-sm">{label}</h3>
                  </div>
                  <span className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                    {tasksByStatus[status]?.length || 0}
                  </span>
                </div>

                {/* Column Content */}
                <div className="space-y-2 min-h-[200px] p-2 bg-muted/30 rounded-lg">
                  {tasksByStatus[status]?.length > 0 ? (
                    tasksByStatus[status].map((task) => (
                      <TaskCard
                        key={task.id}
                        task={task}
                        onClick={handleTaskClick}
                        showProject={!filters.project_id}
                      />
                    ))
                  ) : (
                    <div className="flex items-center justify-center h-[100px] text-sm text-muted-foreground">
                      No tasks
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>

          {/* Won't Do section */}
          {tasksByStatus.wont_do?.length > 0 && (
            <div className="mt-6">
              <div className="flex items-center gap-2 mb-3">
                <div className="w-2 h-2 rounded-full bg-red-500" />
                <h3 className="font-medium text-sm">Won&apos;t Do</h3>
                <span className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                  {tasksByStatus.wont_do.length}
                </span>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2">
                {tasksByStatus.wont_do.map((task) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    onClick={handleTaskClick}
                    showProject={!filters.project_id}
                  />
                ))}
              </div>
            </div>
          )}
        </>
      )}

      {/* Task Detail Dialog */}
      <TaskDetailDialog
        task={selectedTask}
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        onStatusChange={handleStatusChange}
      />
    </div>
  );
}
