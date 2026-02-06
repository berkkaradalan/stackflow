package models

import "time"

// Task status constants
const (
	TaskStatusOpen       = "open"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
	TaskStatusClosed     = "closed"
	TaskStatusWontDo     = "wont_do"
)

// Task priority constants
const (
	TaskPriorityLow      = "low"
	TaskPriorityMedium   = "medium"
	TaskPriorityHigh     = "high"
	TaskPriorityCritical = "critical"
)

// Creator type constants
const (
	CreatorTypeUser  = "user"
	CreatorTypeAgent = "agent"
)

// Task activity action constants
const (
	TaskActionCreated       = "created"
	TaskActionStatusChanged = "status_changed"
	TaskActionAssigned      = "assigned"
	TaskActionReviewerSet   = "reviewer_set"
	TaskActionCommented     = "commented"
	TaskActionProgress      = "progress"
)

// Task represents a task in the system
type Task struct {
	ID              int       `json:"id"`
	ProjectID       int       `json:"project_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	Priority        string    `json:"priority"`
	AssignedAgentID *int      `json:"assigned_agent_id,omitempty"`
	ReviewerID      *int      `json:"reviewer_id,omitempty"`
	CreatedBy       int       `json:"created_by"`
	CreatorType     string    `json:"creator_type"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TaskWithDetails includes related entity names for display
type TaskWithDetails struct {
	Task
	AssignedAgentName *string `json:"assigned_agent_name,omitempty"`
	ReviewerName      *string `json:"reviewer_name,omitempty"`
	CreatorName       string  `json:"creator_name"`
	ProjectName       string  `json:"project_name"`
}

// TaskActivity represents an activity log entry for a task
type TaskActivity struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	ActorID   int       `json:"actor_id"`
	ActorType string    `json:"actor_type"`
	Action    string    `json:"action"`
	OldValue  *string   `json:"old_value,omitempty"`
	NewValue  *string   `json:"new_value,omitempty"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskActivityWithDetails includes actor name for display
type TaskActivityWithDetails struct {
	TaskActivity
	ActorName string `json:"actor_name"`
}

// CreateTaskRequest is the request model for creating a task
type CreateTaskRequest struct {
	Title           string   `json:"title" binding:"required,min=3,max=200"`
	Description     string   `json:"description" binding:"omitempty,max=2000"`
	Priority        string   `json:"priority" binding:"omitempty,oneof=low medium high critical"`
	AssignedAgentID *int     `json:"assigned_agent_id" binding:"omitempty"`
	ReviewerID      *int     `json:"reviewer_id" binding:"omitempty"`
	Tags            []string `json:"tags" binding:"omitempty"`
}

// UpdateTaskRequest is the request model for updating a task
type UpdateTaskRequest struct {
	Title       *string  `json:"title" binding:"omitempty,min=3,max=200"`
	Description *string  `json:"description" binding:"omitempty,max=2000"`
	Priority    *string  `json:"priority" binding:"omitempty,oneof=low medium high critical"`
	Tags        []string `json:"tags" binding:"omitempty"`
}

// AssignAgentRequest is the request model for assigning an agent to a task
type AssignAgentRequest struct {
	AgentID int `json:"agent_id" binding:"required"`
}

// SetReviewerRequest is the request model for setting a reviewer
type SetReviewerRequest struct {
	ReviewerID int `json:"reviewer_id" binding:"required"`
}

// TaskStatusChangeRequest is the request model for status changes with optional message
type TaskStatusChangeRequest struct {
	Message string `json:"message" binding:"omitempty,max=1000"`
}

// CreateActivityRequest is the request model for adding a progress update
type CreateActivityRequest struct {
	Message string `json:"message" binding:"required,min=1,max=2000"`
}

// TaskFilters is the model for filtering tasks
type TaskFilters struct {
	ProjectID       *int    `form:"project_id"`
	Status          *string `form:"status"`
	Priority        *string `form:"priority"`
	AssignedAgentID *int    `form:"assigned_agent_id"`
	ReviewerID      *int    `form:"reviewer_id"`
}

// TaskListResponse is the response model for listing tasks
type TaskListResponse struct {
	Tasks      []TaskWithDetails `json:"tasks"`
	TotalCount int               `json:"total_count"`
}

// TaskActivityListResponse is the response model for listing task activities
type TaskActivityListResponse struct {
	Activities []TaskActivityWithDetails `json:"activities"`
	TotalCount int                       `json:"total_count"`
}
