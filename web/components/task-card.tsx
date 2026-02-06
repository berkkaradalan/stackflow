"use client";

import { Bot, User, Tag, MessageSquare } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  type TaskWithDetails,
  getPriorityColor,
  getPriorityLabel,
} from "@/lib/api/tasks";

interface TaskCardProps {
  task: TaskWithDetails;
  onClick?: (task: TaskWithDetails) => void;
}

export function TaskCard({ task, onClick }: TaskCardProps) {
  return (
    <Card
      className="cursor-pointer transition-shadow hover:shadow-md"
      onClick={() => onClick?.(task)}
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
