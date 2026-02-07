package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	repository "github.com/berkkaradalan/stackflow/repository/postgres"
)

var (
	ErrExecutionPlanNotFound = errors.New("execution plan not found")
	ErrNoActivePlan          = errors.New("no active execution plan found")
	ErrAssignmentNotFound    = errors.New("assignment not found")
	ErrNoTaskAvailable       = errors.New("no pending task available for agent")
)

type ExecutionPlanService struct {
	planRepo    *repository.ExecutionPlanRepository
	projectRepo *repository.ProjectRepository
	agentRepo   *repository.AgentRepository
	taskRepo    *repository.TaskRepository
}

func NewExecutionPlanService(
	planRepo *repository.ExecutionPlanRepository,
	projectRepo *repository.ProjectRepository,
	agentRepo *repository.AgentRepository,
	taskRepo *repository.TaskRepository,
) *ExecutionPlanService {
	return &ExecutionPlanService{
		planRepo:    planRepo,
		projectRepo: projectRepo,
		agentRepo:   agentRepo,
		taskRepo:    taskRepo,
	}
}

// CreatePlan creates a new execution plan for a project
func (s *ExecutionPlanService) CreatePlan(ctx context.Context, projectID int, req *models.CreateExecutionPlanRequest, creatorID int, creatorType string) (*models.ExecutionPlanWithDetails, error) {
	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	plan := &models.ExecutionPlan{
		ProjectID:   projectID,
		CreatedBy:   creatorID,
		CreatorType: creatorType,
		PlanData:    req.PlanData,
		Status:      models.ExecutionPlanStatusActive,
	}

	if plan.PlanData.FocusAreas == nil {
		plan.PlanData.FocusAreas = []string{}
	}

	err = s.planRepo.CreatePlan(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution plan: %w", err)
	}

	return s.planRepo.GetPlanByIDWithDetails(ctx, plan.ID)
}

// GetActivePlan retrieves the active execution plan for a project
func (s *ExecutionPlanService) GetActivePlan(ctx context.Context, projectID int) (*models.ExecutionPlanWithDetails, error) {
	plan, err := s.planRepo.GetActivePlanByProjectID(ctx, projectID)
	if err != nil {
		return nil, ErrNoActivePlan
	}
	return plan, nil
}

// GetPlanByID retrieves an execution plan by ID
func (s *ExecutionPlanService) GetPlanByID(ctx context.Context, id int) (*models.ExecutionPlanWithDetails, error) {
	plan, err := s.planRepo.GetPlanByIDWithDetails(ctx, id)
	if err != nil {
		return nil, ErrExecutionPlanNotFound
	}
	return plan, nil
}

// GetPlansByProject retrieves all plans for a project
func (s *ExecutionPlanService) GetPlansByProject(ctx context.Context, projectID int) (*models.ExecutionPlanListResponse, error) {
	plans, err := s.planRepo.GetPlansByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	if plans == nil {
		plans = []models.ExecutionPlanWithDetails{}
	}

	return &models.ExecutionPlanListResponse{
		Plans:      plans,
		TotalCount: len(plans),
	}, nil
}

// UpdatePlan updates an execution plan
func (s *ExecutionPlanService) UpdatePlan(ctx context.Context, projectID int, req *models.UpdateExecutionPlanRequest, actorID int, actorType string) (*models.ExecutionPlanWithDetails, error) {
	// Get active plan for this project
	plan, err := s.planRepo.GetActivePlanByProjectID(ctx, projectID)
	if err != nil {
		return nil, ErrNoActivePlan
	}

	updates := make(map[string]interface{})

	if req.PlanData != nil {
		updates["plan_data"] = req.PlanData
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	err = s.planRepo.UpdatePlan(ctx, plan.ID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	return s.planRepo.GetPlanByIDWithDetails(ctx, plan.ID)
}

// --- Agent Task Flow ---

// GetNextTask finds and returns the next pending task for an agent
func (s *ExecutionPlanService) GetNextTask(ctx context.Context, agentID int) (*models.NextTaskResponse, error) {
	// Verify agent exists
	_, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		return nil, ErrAgentNotFound
	}

	assignment, err := s.planRepo.GetNextAssignmentForAgent(ctx, agentID)
	if err != nil {
		return &models.NextTaskResponse{
			Message: "No pending tasks available",
		}, nil
	}

	// Mark the assignment as in_progress
	_ = s.planRepo.StartAssignment(ctx, assignment.ID)
	assignment.Status = models.AssignmentStatusInProgress

	// Get the plan context for this assignment
	plan, err := s.planRepo.GetPlanByID(ctx, assignment.PlanID)
	if err != nil {
		return &models.NextTaskResponse{
			Assignment: assignment,
			Message:    "Task assigned",
		}, nil
	}

	return &models.NextTaskResponse{
		Assignment: assignment,
		Context: &models.AgentContextResponse{
			PlanID:      plan.ID,
			Constraints: plan.PlanData.Constraints,
			FocusAreas:  plan.PlanData.FocusAreas,
			Notes:       plan.PlanData.Notes,
		},
		Message: "Task assigned",
	}, nil
}

// CompleteTask handles an agent reporting task completion
func (s *ExecutionPlanService) CompleteTask(ctx context.Context, agentID int, req *models.TaskCompleteRequest) (*models.AgentAssignment, error) {
	// Verify agent exists
	_, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		return nil, ErrAgentNotFound
	}

	// Find the active assignment for this agent and task
	assignment, err := s.planRepo.GetAssignmentByAgentAndTask(ctx, agentID, req.TaskID)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}

	// Complete the assignment
	err = s.planRepo.CompleteAssignment(ctx, assignment.ID, req.ReportData)
	if err != nil {
		return nil, fmt.Errorf("failed to complete assignment: %w", err)
	}

	// Also update the task status to done
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err == nil && task.Status == models.TaskStatusInProgress {
		_ = s.taskRepo.UpdateStatus(ctx, req.TaskID, models.TaskStatusDone)

		// Log activity on the task
		if req.Message == "" {
			req.Message = "Task completed by agent"
		}
		newStatus := models.TaskStatusDone
		activity := &models.TaskActivity{
			TaskID:    req.TaskID,
			ActorID:   agentID,
			ActorType: models.CreatorTypeAgent,
			Action:    models.TaskActionStatusChanged,
			OldValue:  &task.Status,
			NewValue:  &newStatus,
			Message:   req.Message,
		}
		_ = s.taskRepo.CreateActivity(ctx, activity)
	}

	return s.planRepo.GetAssignmentByID(ctx, assignment.ID)
}

// GetAgentContext retrieves PM constraints/notes for an agent from the active plan
func (s *ExecutionPlanService) GetAgentContext(ctx context.Context, agentID int) (*models.AgentContextResponse, error) {
	// Verify agent exists and get their project
	agent, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		return nil, ErrAgentNotFound
	}

	// Find active plan for the agent's project
	plan, err := s.planRepo.GetActivePlanByProjectID(ctx, agent.ProjectID)
	if err != nil {
		return nil, ErrNoActivePlan
	}

	return &models.AgentContextResponse{
		PlanID:      plan.ID,
		Constraints: plan.PlanData.Constraints,
		FocusAreas:  plan.PlanData.FocusAreas,
		Notes:       plan.PlanData.Notes,
	}, nil
}

// --- Reports ---

// GenerateReport creates a new execution report
func (s *ExecutionPlanService) GenerateReport(ctx context.Context, projectID int, req *models.GenerateReportRequest, generatorID int, generatorType string) (*models.ExecutionReport, error) {
	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	report := &models.ExecutionReport{
		ProjectID:     projectID,
		ReportType:    req.ReportType,
		GeneratedBy:   generatorID,
		GeneratorType: generatorType,
		ReportData:    req.ReportData,
	}

	err = s.planRepo.CreateReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

// GetReports retrieves all reports for a project
func (s *ExecutionPlanService) GetReports(ctx context.Context, projectID int) (*models.ExecutionReportListResponse, error) {
	reports, err := s.planRepo.GetReportsByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	if reports == nil {
		reports = []models.ExecutionReportWithDetails{}
	}

	return &models.ExecutionReportListResponse{
		Reports:    reports,
		TotalCount: len(reports),
	}, nil
}

// GetReportsByType retrieves reports filtered by type
func (s *ExecutionPlanService) GetReportsByType(ctx context.Context, projectID int, reportType string) (*models.ExecutionReportListResponse, error) {
	reports, err := s.planRepo.GetReportsByType(ctx, projectID, reportType)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	if reports == nil {
		reports = []models.ExecutionReportWithDetails{}
	}

	return &models.ExecutionReportListResponse{
		Reports:    reports,
		TotalCount: len(reports),
	}, nil
}
