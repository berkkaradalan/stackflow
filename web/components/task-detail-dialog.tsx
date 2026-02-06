"use client";

import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import {
  Bot,
  User,
  Tag,
  Calendar,
  Play,
  CheckCircle,
  XCircle,
  RotateCcw,
  Ban,
  Loader2,
  Send,
} from "lucide-react";
import {
  type TaskWithDetails,
  type TaskActivityWithDetails,
  type TaskStatus,
  getStatusColor,
  getStatusLabel,
  getPriorityColor,
  getPriorityLabel,
} from "@/lib/api/tasks";
import { useTaskActivities } from "@/hooks/use-tasks";

interface TaskDetailDialogProps {
  task: TaskWithDetails | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onStatusChange?: (
    taskId: number,
    action: "start" | "complete" | "close" | "wontdo" | "reopen"
  ) => Promise<void>;
}

// Get available actions based on current status
function getAvailableActions(status: TaskStatus) {
  const actions: {
    action: "start" | "complete" | "close" | "wontdo" | "reopen";
    label: string;
    icon: React.ReactNode;
    variant: "default" | "secondary" | "destructive" | "outline";
  }[] = [];

  switch (status) {
    case "open":
      actions.push({
        action: "start",
        label: "Start",
        icon: <Play className="h-4 w-4" />,
        variant: "default",
      });
      actions.push({
        action: "wontdo",
        label: "Won't Do",
        icon: <Ban className="h-4 w-4" />,
        variant: "destructive",
      });
      break;
    case "in_progress":
      actions.push({
        action: "complete",
        label: "Done",
        icon: <CheckCircle className="h-4 w-4" />,
        variant: "default",
      });
      actions.push({
        action: "wontdo",
        label: "Won't Do",
        icon: <Ban className="h-4 w-4" />,
        variant: "destructive",
      });
      break;
    case "done":
      actions.push({
        action: "close",
        label: "Close",
        icon: <XCircle className="h-4 w-4" />,
        variant: "secondary",
      });
      actions.push({
        action: "reopen",
        label: "Reopen",
        icon: <RotateCcw className="h-4 w-4" />,
        variant: "outline",
      });
      break;
    case "closed":
    case "wont_do":
      actions.push({
        action: "reopen",
        label: "Reopen",
        icon: <RotateCcw className="h-4 w-4" />,
        variant: "outline",
      });
      break;
  }

  return actions;
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function ActivityItem({ activity }: { activity: TaskActivityWithDetails }) {
  const getActionIcon = (action: string) => {
    switch (action) {
      case "created":
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case "status_changed":
        return <RotateCcw className="h-4 w-4 text-blue-500" />;
      case "assigned":
        return <Bot className="h-4 w-4 text-purple-500" />;
      case "reviewer_set":
        return <User className="h-4 w-4 text-orange-500" />;
      case "progress":
        return <Send className="h-4 w-4 text-cyan-500" />;
      default:
        return <Calendar className="h-4 w-4 text-gray-500" />;
    }
  };

  return (
    <div className="flex gap-3 py-2">
      <div className="mt-0.5">{getActionIcon(activity.action)}</div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 text-sm">
          <span className="font-medium">{activity.actor_name}</span>
          {activity.actor_type === "agent" && (
            <Bot className="h-3 w-3 text-muted-foreground" />
          )}
          <span className="text-muted-foreground text-xs">
            {formatDate(activity.created_at)}
          </span>
        </div>
        <p className="text-sm text-muted-foreground mt-0.5">{activity.message}</p>
        {activity.old_value && activity.new_value && (
          <div className="flex items-center gap-2 mt-1 text-xs">
            <Badge variant="outline" className="text-xs">
              {activity.old_value}
            </Badge>
            <span>→</span>
            <Badge variant="outline" className="text-xs">
              {activity.new_value}
            </Badge>
          </div>
        )}
      </div>
    </div>
  );
}

export function TaskDetailDialog({
  task,
  open,
  onOpenChange,
  onStatusChange,
}: TaskDetailDialogProps) {
  const [isChangingStatus, setIsChangingStatus] = useState(false);
  const [progressMessage, setProgressMessage] = useState("");
  const [isAddingProgress, setIsAddingProgress] = useState(false);

  const { activities, isLoading: activitiesLoading, fetchActivities, addProgress } =
    useTaskActivities(task?.id || 0);

  useEffect(() => {
    if (open && task) {
      fetchActivities();
    }
  }, [open, task, fetchActivities]);

  const handleStatusChange = async (
    action: "start" | "complete" | "close" | "wontdo" | "reopen"
  ) => {
    if (!task || !onStatusChange) return;
    setIsChangingStatus(true);
    try {
      await onStatusChange(task.id, action);
    } finally {
      setIsChangingStatus(false);
    }
  };

  const handleAddProgress = async () => {
    if (!progressMessage.trim()) return;
    setIsAddingProgress(true);
    try {
      await addProgress(progressMessage);
      setProgressMessage("");
    } finally {
      setIsAddingProgress(false);
    }
  };

  if (!task) return null;

  const availableActions = getAvailableActions(task.status);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[85vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <div className="flex items-start gap-3">
            <div className="flex-1">
              <DialogTitle className="text-xl">{task.title}</DialogTitle>
              <DialogDescription className="mt-1">
                {task.project_name} • Created by {task.creator_name}
              </DialogDescription>
            </div>
            <div className="flex gap-2">
              <Badge className={getStatusColor(task.status)}>
                {getStatusLabel(task.status)}
              </Badge>
              <Badge className={getPriorityColor(task.priority)}>
                {getPriorityLabel(task.priority)}
              </Badge>
            </div>
          </div>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto space-y-4">
          {/* Description */}
          {task.description && (
            <div>
              <h4 className="text-sm font-medium mb-1">Description</h4>
              <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                {task.description}
              </p>
            </div>
          )}

          {/* Tags */}
          {task.tags && task.tags.length > 0 && (
            <div>
              <h4 className="text-sm font-medium mb-2">Tags</h4>
              <div className="flex flex-wrap gap-1">
                {task.tags.map((tag, index) => (
                  <Badge key={index} variant="outline">
                    <Tag className="h-3 w-3 mr-1" />
                    {tag}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {/* Assignment Info */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h4 className="text-sm font-medium mb-1">Assigned Agent</h4>
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Bot className="h-4 w-4" />
                {task.assigned_agent_name || "Not assigned"}
              </div>
            </div>
            <div>
              <h4 className="text-sm font-medium mb-1">Reviewer</h4>
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <User className="h-4 w-4" />
                {task.reviewer_name || "Not assigned"}
              </div>
            </div>
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-muted-foreground">Created: </span>
              {formatDate(task.created_at)}
            </div>
            <div>
              <span className="text-muted-foreground">Updated: </span>
              {formatDate(task.updated_at)}
            </div>
          </div>

          <Separator />

          {/* Actions */}
          {availableActions.length > 0 && (
            <div>
              <h4 className="text-sm font-medium mb-2">Actions</h4>
              <div className="flex flex-wrap gap-2">
                {availableActions.map(({ action, label, icon, variant }) => (
                  <Button
                    key={action}
                    variant={variant}
                    size="sm"
                    onClick={() => handleStatusChange(action)}
                    disabled={isChangingStatus}
                  >
                    {isChangingStatus ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      icon
                    )}
                    <span className="ml-1">{label}</span>
                  </Button>
                ))}
              </div>
            </div>
          )}

          <Separator />

          {/* Activity Log */}
          <div>
            <h4 className="text-sm font-medium mb-2">Activity</h4>

            {/* Add Progress */}
            <div className="flex gap-2 mb-4">
              <Textarea
                placeholder="Add a progress update..."
                value={progressMessage}
                onChange={(e) => setProgressMessage(e.target.value)}
                className="min-h-[60px] text-sm"
              />
              <Button
                size="sm"
                onClick={handleAddProgress}
                disabled={isAddingProgress || !progressMessage.trim()}
              >
                {isAddingProgress ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Send className="h-4 w-4" />
                )}
              </Button>
            </div>

            {/* Activity List */}
            {activitiesLoading ? (
              <div className="flex justify-center py-4">
                <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
              </div>
            ) : activities.length > 0 ? (
              <div className="space-y-1 divide-y">
                {activities.map((activity) => (
                  <ActivityItem key={activity.id} activity={activity} />
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground text-center py-4">
                No activity yet
              </p>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
