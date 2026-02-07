package models

import "time"

// Execution plan status constants
const (
	ExecutionPlanStatusActive    = "active"
	ExecutionPlanStatusCompleted = "completed"
	ExecutionPlanStatusCancelled = "cancelled"
	ExecutionPlanStatusDraft     = "draft"
)

// Agent assignment status constants
const (
	AssignmentStatusPending    = "pending"
	AssignmentStatusInProgress = "in_progress"
	AssignmentStatusCompleted  = "completed"
	AssignmentStatusFailed     = "failed"
	AssignmentStatusSkipped    = "skipped"
)

// Execution report type constants
const (
	ReportTypeDaily   = "daily"
	ReportTypeWeekly  = "weekly"
	ReportTypeCustom  = "custom"
	ReportTypeSummary = "summary"
)

// ExecutionPlan represents a PM-generated execution plan for a project
type ExecutionPlan struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	CreatedBy   int       `json:"created_by"`
	CreatorType string    `json:"creator_type"`
	PlanData    PlanData  `json:"plan_data"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlanData holds the JSONB plan structure
type PlanData struct {
	PriorityOrder []TaskPriorityItem `json:"priority_order"`
	Constraints   PlanConstraints    `json:"constraints"`
	FocusAreas    []string           `json:"focus_areas"`
	Notes         string             `json:"notes,omitempty"`
}

// TaskPriorityItem represents a single task's priority and assignment within the plan
type TaskPriorityItem struct {
	TaskID          int               `json:"task_id"`
	Title           string            `json:"title"`
	AssignedAgentID *int              `json:"assigned_agent_id,omitempty"`
	Priority        string            `json:"priority"`
	Dependencies    []int             `json:"dependencies"`
	Constraints     map[string]any    `json:"constraints,omitempty"`
	EstimatedEffort string            `json:"estimated_effort"`
	Notes           string            `json:"notes,omitempty"`
}

// PlanConstraints holds global constraints for the execution plan
type PlanConstraints struct {
	MaxParallelTasks   int  `json:"max_parallel_tasks"`
	CodeReviewRequired bool `json:"code_review_required"`
	TestCoverageMin    int  `json:"test_coverage_min"`
}

// ExecutionPlanWithDetails includes related entity names for display
type ExecutionPlanWithDetails struct {
	ExecutionPlan
	ProjectName string `json:"project_name"`
	CreatorName string `json:"creator_name"`
}

// AgentAssignment represents a task assignment to an agent within a plan
type AgentAssignment struct {
	ID          int        `json:"id"`
	PlanID      int        `json:"plan_id"`
	AgentID     int        `json:"agent_id"`
	TaskID      int        `json:"task_id"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ReportData  any        `json:"report_data,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AgentAssignmentWithDetails includes related entity names
type AgentAssignmentWithDetails struct {
	AgentAssignment
	AgentName string `json:"agent_name"`
	TaskTitle string `json:"task_title"`
}

// ExecutionReport represents a generated execution report
type ExecutionReport struct {
	ID            int       `json:"id"`
	ProjectID     int       `json:"project_id"`
	ReportType    string    `json:"report_type"`
	GeneratedBy   int       `json:"generated_by"`
	GeneratorType string    `json:"generator_type"`
	ReportData    any       `json:"report_data"`
	CreatedAt     time.Time `json:"created_at"`
}

// ExecutionReportWithDetails includes related entity names
type ExecutionReportWithDetails struct {
	ExecutionReport
	ProjectName   string `json:"project_name"`
	GeneratorName string `json:"generator_name"`
}

// --- Request DTOs ---

// CreateExecutionPlanRequest is the request model for creating an execution plan
type CreateExecutionPlanRequest struct {
	PlanData PlanData `json:"plan_data" binding:"required"`
}

// UpdateExecutionPlanRequest is the request model for updating an execution plan
type UpdateExecutionPlanRequest struct {
	PlanData *PlanData `json:"plan_data" binding:"omitempty"`
	Status   *string   `json:"status" binding:"omitempty,oneof=active completed cancelled draft"`
}

// TaskCompleteRequest is the request model for an agent reporting task completion
type TaskCompleteRequest struct {
	TaskID     int    `json:"task_id" binding:"required"`
	ReportData any    `json:"report_data" binding:"omitempty"`
	Message    string `json:"message" binding:"omitempty,max=2000"`
}

// GenerateReportRequest is the request model for generating a report
type GenerateReportRequest struct {
	ReportType string `json:"report_type" binding:"required,oneof=daily weekly custom summary"`
	ReportData any    `json:"report_data" binding:"required"`
}

// --- Response DTOs ---

// ExecutionPlanListResponse is the response model for listing execution plans
type ExecutionPlanListResponse struct {
	Plans      []ExecutionPlanWithDetails `json:"plans"`
	TotalCount int                        `json:"total_count"`
}

// AgentAssignmentListResponse is the response model for listing agent assignments
type AgentAssignmentListResponse struct {
	Assignments []AgentAssignmentWithDetails `json:"assignments"`
	TotalCount  int                          `json:"total_count"`
}

// NextTaskResponse is the response model for an agent asking what to do next
type NextTaskResponse struct {
	Assignment *AgentAssignmentWithDetails `json:"assignment,omitempty"`
	Context    *AgentContextResponse       `json:"context,omitempty"`
	Message    string                      `json:"message"`
}

// AgentContextResponse provides PM constraints and notes for the agent
type AgentContextResponse struct {
	PlanID      int             `json:"plan_id"`
	Constraints PlanConstraints `json:"constraints"`
	FocusAreas  []string        `json:"focus_areas"`
	Notes       string          `json:"notes,omitempty"`
}

// ExecutionReportListResponse is the response model for listing reports
type ExecutionReportListResponse struct {
	Reports    []ExecutionReportWithDetails `json:"reports"`
	TotalCount int                          `json:"total_count"`
}
