"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAgents } from "@/hooks/use-agents";
import { useProjects } from "@/hooks/use-projects";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  Avatar,
  AvatarImage,
  AvatarFallback,
} from "@/components/ui/avatar";
import {
  Loader2,
  Plus,
  Pencil,
  Trash2,
  Search,
  Bot,
  Eye,
  Activity,
} from "lucide-react";
import type { Agent, CreateAgentRequest } from "@/lib/api/agents";

const ROLES = [
  { value: "backend_developer", label: "Backend Developer" },
  { value: "frontend_developer", label: "Frontend Developer" },
  { value: "fullstack_developer", label: "Fullstack Developer" },
  { value: "tester", label: "Tester" },
  { value: "devops", label: "DevOps" },
  { value: "project_manager", label: "Project Manager" },
];

const LEVELS = [
  { value: "junior", label: "Junior" },
  { value: "mid", label: "Mid" },
  { value: "senior", label: "Senior" },
];

const PROVIDERS = [
  { value: "openrouter", label: "OpenRouter" },
  { value: "openai", label: "OpenAI" },
  { value: "anthropic", label: "Anthropic" },
  { value: "gemini", label: "Gemini" },
  { value: "groq", label: "Groq" },
  { value: "together", label: "Together" },
  { value: "glm", label: "GLM" },
  { value: "claude", label: "Claude" },
  { value: "kimi", label: "Kimi" },
];

export default function AgentsPage() {
  const router = useRouter();
  const {
    agents,
    isLoading,
    error,
    fetchAgents,
    createAgent,
    deleteAgent,
  } = useAgents();

  const { projects, fetchProjects } = useProjects();

  const [searchQuery, setSearchQuery] = useState("");
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);

  const [createForm, setCreateForm] = useState<CreateAgentRequest>({
    name: "",
    description: "",
    project_id: 0,
    role: "backend_developer",
    level: "mid",
    provider: "openai",
    model: "",
    api_key: "",
    config: {
      temperature: 0.7,
      max_tokens: 2000,
      top_p: 1.0,
      frequency_penalty: 0.0,
      presence_penalty: 0.0,
    },
  });

  useEffect(() => {
    fetchAgents();
    fetchProjects();
  }, [fetchAgents, fetchProjects]);

  const filteredAgents = (agents || []).filter(
    (agent) =>
      agent &&
      (agent.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.role?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.level?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.provider?.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const handleCreateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createAgent(createForm);
      setCreateDialogOpen(false);
      setCreateForm({
        name: "",
        description: "",
        project_id: 0,
        role: "backend_developer",
        level: "mid",
        provider: "openai",
        model: "",
        api_key: "",
        config: {
          temperature: 0.7,
          max_tokens: 2000,
          top_p: 1.0,
          frequency_penalty: 0.0,
          presence_penalty: 0.0,
        },
      });
    } catch (err) {
      console.error("Failed to create agent:", err);
    }
  };

  const handleDeleteClick = (agent: Agent) => {
    setSelectedAgent(agent);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedAgent) return;

    try {
      await deleteAgent(selectedAgent.id);
      setDeleteDialogOpen(false);
      setSelectedAgent(null);
    } catch (err) {
      console.error("Failed to delete agent:", err);
    }
  };

  const handleViewAgent = (agentId: number) => {
    router.push(`/agents/${agentId}`);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100";
      case "idle":
        return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-100";
      case "busy":
        return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-100";
      case "error":
        return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-100";
      case "disabled":
        return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
    }
  };

  const getLevelColor = (level: string) => {
    switch (level) {
      case "senior":
        return "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-100";
      case "mid":
        return "bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-100";
      case "junior":
        return "bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-100";
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">AI Agents</h1>
        <p className="text-muted-foreground">
          Manage your AI agents and monitor their performance.
        </p>
      </div>

      <Separator />

      {/* Actions Bar */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative flex-1 sm:max-w-sm">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search agents..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          New Agent
        </Button>
      </div>

      {/* Agents Grid */}
      {isLoading ? (
        <div className="flex min-h-[400px] items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <div className="rounded-lg bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      ) : filteredAgents.length === 0 ? (
        <Card>
          <CardContent className="flex min-h-[400px] flex-col items-center justify-center text-center">
            <Bot className="mb-4 h-12 w-12 text-muted-foreground" />
            <h3 className="mb-2 text-lg font-semibold">No agents found</h3>
            <p className="mb-4 text-sm text-muted-foreground">
              {searchQuery
                ? "Try adjusting your search query"
                : "Get started by creating your first AI agent"}
            </p>
            {!searchQuery && (
              <Button onClick={() => setCreateDialogOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                New Agent
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {filteredAgents.map((agent) => (
            <Card key={agent.id} className="hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="flex items-start gap-3">
                  <Avatar size="lg">
                    <AvatarImage
                      src={`https://api.dicebear.com/7.x/bottts/svg?seed=${agent.name}`}
                      alt={agent.name}
                    />
                    <AvatarFallback className="bg-gradient-to-br from-blue-500 to-purple-600 text-white">
                      {agent.name.slice(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="space-y-1 flex-1 min-w-0">
                    <div className="flex items-center justify-between gap-2">
                      <CardTitle className="line-clamp-1 text-base">{agent.name}</CardTitle>
                      <Badge className={getStatusColor(agent.status)} variant="secondary">
                        {agent.status}
                      </Badge>
                    </div>
                    <CardDescription className="line-clamp-2">
                      {agent.description || "No description"}
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className={getLevelColor(agent.level)}>
                    {agent.level}
                  </Badge>
                  <Badge variant="outline">
                    {agent.role.replace(/_/g, " ")}
                  </Badge>
                </div>

                <div className="text-xs text-muted-foreground space-y-1">
                  <div>Provider: <span className="font-medium">{agent.provider}</span></div>
                  <div>Model: <span className="font-medium">{agent.model}</span></div>
                  <div>Requests: <span className="font-medium">{agent.total_requests}</span></div>
                  <div>Cost: <span className="font-medium">${agent.total_cost.toFixed(4)}</span></div>
                </div>

                <div className="flex items-center justify-between pt-2">
                  <div className="text-xs text-muted-foreground">
                    Created {new Date(agent.created_at).toLocaleDateString()}
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleViewAgent(agent.id)}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDeleteClick(agent)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Create Agent Dialog */}
      <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create New Agent</DialogTitle>
            <DialogDescription>
              Configure a new AI agent for your project.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleCreateSubmit}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={createForm.name}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, name: e.target.value })
                  }
                  placeholder="My Backend Agent"
                  required
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={createForm.description}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, description: e.target.value })
                  }
                  placeholder="Handles backend development tasks..."
                  rows={3}
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="project">Project *</Label>
                <Select
                  value={createForm.project_id?.toString() || ""}
                  onValueChange={(value) =>
                    setCreateForm({ ...createForm, project_id: parseInt(value) })
                  }
                  required
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a project" />
                  </SelectTrigger>
                  <SelectContent>
                    {projects.map((project) => (
                      <SelectItem key={project.id} value={project.id.toString()}>
                        {project.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="role">Role *</Label>
                  <Select
                    value={createForm.role}
                    onValueChange={(value: any) =>
                      setCreateForm({ ...createForm, role: value })
                    }
                    required
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {ROLES.map((role) => (
                        <SelectItem key={role.value} value={role.value}>
                          {role.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="level">Level *</Label>
                  <Select
                    value={createForm.level}
                    onValueChange={(value: any) =>
                      setCreateForm({ ...createForm, level: value })
                    }
                    required
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {LEVELS.map((level) => (
                        <SelectItem key={level.value} value={level.value}>
                          {level.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="provider">Provider *</Label>
                  <Select
                    value={createForm.provider}
                    onValueChange={(value: any) =>
                      setCreateForm({ ...createForm, provider: value })
                    }
                    required
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {PROVIDERS.map((provider) => (
                        <SelectItem key={provider.value} value={provider.value}>
                          {provider.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="model">Model *</Label>
                  <Input
                    id="model"
                    value={createForm.model}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, model: e.target.value })
                    }
                    placeholder="gpt-4, claude-3-opus, etc."
                    required
                  />
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="api_key">API Key *</Label>
                <Input
                  id="api_key"
                  type="password"
                  value={createForm.api_key}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, api_key: e.target.value })
                  }
                  placeholder="sk-..."
                  required
                />
              </div>

              <Separator />

              <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="temperature">Temperature</Label>
                  <Input
                    id="temperature"
                    type="number"
                    step="0.1"
                    min="0"
                    max="2"
                    value={createForm.config?.temperature || 0.7}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        config: { ...createForm.config!, temperature: parseFloat(e.target.value) },
                      })
                    }
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="max_tokens">Max Tokens</Label>
                  <Input
                    id="max_tokens"
                    type="number"
                    min="1"
                    value={createForm.config?.max_tokens || 2000}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        config: { ...createForm.config!, max_tokens: parseInt(e.target.value) },
                      })
                    }
                  />
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setCreateDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit">Create Agent</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete the agent "{selectedAgent?.name}".
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
