"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useAgents } from "@/hooks/use-agents";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  Avatar,
  AvatarImage,
  AvatarFallback,
} from "@/components/ui/avatar";
import {
  Loader2,
  ArrowLeft,
  Activity,
  TrendingUp,
  Zap,
  CheckCircle2,
  XCircle,
  AlertCircle,
  DollarSign,
  Hash,
  Clock,
} from "lucide-react";
import type { Agent, AgentHealthResponse, AgentWorkloadResponse, AgentPerformanceResponse } from "@/lib/api/agents";

export default function AgentDetailPage() {
  const params = useParams();
  const router = useRouter();
  const agentId = parseInt(params.id as string);

  const { fetchAgentById, healthCheck, getWorkload, getPerformance } = useAgents();

  const [agent, setAgent] = useState<Agent | null>(null);
  const [health, setHealth] = useState<AgentHealthResponse | null>(null);
  const [workload, setWorkload] = useState<AgentWorkloadResponse | null>(null);
  const [performance, setPerformance] = useState<AgentPerformanceResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadAgentData = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const [agentData, healthData, workloadData, performanceData] = await Promise.all([
          fetchAgentById(agentId),
          healthCheck(agentId),
          getWorkload(agentId),
          getPerformance(agentId),
        ]);

        setAgent(agentData);
        setHealth(healthData);
        setWorkload(workloadData);
        setPerformance(performanceData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load agent data");
      } finally {
        setIsLoading(false);
      }
    };

    loadAgentData();
  }, [agentId, fetchAgentById, healthCheck, getWorkload, getPerformance]);

  const handleHealthCheck = async () => {
    try {
      const healthData = await healthCheck(agentId);
      setHealth(healthData);
    } catch (err) {
      console.error("Health check failed:", err);
    }
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

  if (isLoading) {
    return (
      <div className="flex min-h-[600px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !agent) {
    return (
      <div className="space-y-6">
        <Button variant="ghost" onClick={() => router.push("/agents")}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Agents
        </Button>
        <div className="rounded-lg bg-destructive/10 p-4 text-sm text-destructive">
          {error || "Agent not found"}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Button variant="ghost" onClick={() => router.push("/agents")}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Agents
        </Button>
      </div>

      {/* Agent Info */}
      <div className="space-y-3">
        <div className="flex items-start gap-4">
          <Avatar size="lg" className="h-16 w-16">
            <AvatarImage
              src={`https://api.dicebear.com/7.x/bottts/svg?seed=${agent.name}`}
              alt={agent.name}
            />
            <AvatarFallback className="bg-gradient-to-br from-blue-500 to-purple-600 text-white text-2xl">
              {agent.name.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 space-y-2">
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold tracking-tight">{agent.name}</h1>
              <Badge className={getStatusColor(agent.status)}>{agent.status}</Badge>
            </div>
            <p className="text-muted-foreground">
              {agent.description || "No description provided"}
            </p>
          </div>
        </div>
      </div>

      <Separator />

      {/* Agent Details Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Role & Level</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              <div className="text-2xl font-bold capitalize">
                {agent.role.replace(/_/g, " ")}
              </div>
              <p className="text-xs text-muted-foreground capitalize">
                {agent.level} level
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Provider</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              <div className="text-2xl font-bold capitalize">{agent.provider}</div>
              <p className="text-xs text-muted-foreground">{agent.model}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Total Requests</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <Hash className="h-4 w-4 text-muted-foreground" />
              <div className="text-2xl font-bold">{agent.total_requests.toLocaleString()}</div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium">Total Cost</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <DollarSign className="h-4 w-4 text-muted-foreground" />
              <div className="text-2xl font-bold">${agent.total_cost.toFixed(4)}</div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Health Status */}
      {health && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  {health.healthy ? (
                    <CheckCircle2 className="h-5 w-5 text-green-600" />
                  ) : (
                    <XCircle className="h-5 w-5 text-red-600" />
                  )}
                  Health Status
                </CardTitle>
                <CardDescription>{health.message}</CardDescription>
              </div>
              <Button onClick={handleHealthCheck} variant="outline" size="sm">
                <Activity className="mr-2 h-4 w-4" />
                Check Health
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Status</p>
                <div className="flex items-center gap-2">
                  <Badge className={getStatusColor(health.status)}>
                    {health.status}
                  </Badge>
                </div>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Active</p>
                <div className="flex items-center gap-2">
                  {health.is_active ? (
                    <Badge className="bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100">
                      Yes
                    </Badge>
                  ) : (
                    <Badge variant="secondary">No</Badge>
                  )}
                </div>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Last Active</p>
                <p className="text-sm font-medium">
                  {health.last_active_at
                    ? new Date(health.last_active_at).toLocaleString()
                    : "Never"}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Workload Metrics */}
      {workload && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Zap className="h-5 w-5" />
              Workload Metrics
            </CardTitle>
            <CardDescription>{workload.description}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Total Requests</p>
                <p className="text-2xl font-bold">{workload.total_requests.toLocaleString()}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Total Tokens Used</p>
                <p className="text-2xl font-bold">{workload.total_tokens_used.toLocaleString()}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Total Cost</p>
                <p className="text-2xl font-bold">${workload.total_cost.toFixed(4)}</p>
              </div>
            </div>
            <Separator className="my-4" />
            <div className="grid gap-2 md:grid-cols-2">
              <div className="flex items-center gap-2 text-sm">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">Last Active:</span>
                <span className="font-medium">
                  {workload.last_active_at
                    ? new Date(workload.last_active_at).toLocaleString()
                    : "Never"}
                </span>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">Updated:</span>
                <span className="font-medium">
                  {new Date(workload.updated_at).toLocaleString()}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Performance Metrics */}
      {performance && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Performance Metrics
            </CardTitle>
            <CardDescription>{performance.description}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Average Tokens per Request</p>
                <p className="text-2xl font-bold">
                  {performance.average_tokens_per_request.toFixed(2)}
                </p>
                <p className="text-xs text-muted-foreground">
                  {performance.total_tokens_used.toLocaleString()} tokens / {performance.total_requests.toLocaleString()} requests
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Average Cost per Request</p>
                <p className="text-2xl font-bold">
                  ${performance.average_cost_per_request.toFixed(6)}
                </p>
                <p className="text-xs text-muted-foreground">
                  ${performance.total_cost.toFixed(4)} total / {performance.total_requests.toLocaleString()} requests
                </p>
              </div>
            </div>
            <Separator className="my-4" />
            {performance.total_requests === 0 && (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <AlertCircle className="h-4 w-4" />
                No requests processed yet. Performance metrics will appear after first use.
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Configuration */}
      <Card>
        <CardHeader>
          <CardTitle>Configuration</CardTitle>
          <CardDescription>Agent settings and parameters</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Temperature</p>
              <p className="text-lg font-medium">{agent.config.temperature}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Max Tokens</p>
              <p className="text-lg font-medium">{agent.config.max_tokens}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Top P</p>
              <p className="text-lg font-medium">{agent.config.top_p}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Frequency Penalty</p>
              <p className="text-lg font-medium">{agent.config.frequency_penalty}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Presence Penalty</p>
              <p className="text-lg font-medium">{agent.config.presence_penalty}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Metadata */}
      <Card>
        <CardHeader>
          <CardTitle>Metadata</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Created At</p>
              <p className="text-sm font-medium">{new Date(agent.created_at).toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Last Updated</p>
              <p className="text-sm font-medium">{new Date(agent.updated_at).toLocaleString()}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
