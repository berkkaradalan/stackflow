package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	repository "github.com/berkkaradalan/stackflow/repository/postgres"
)

var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrUnauthorized          = errors.New("unauthorized to perform this action")
	ErrAgentNotFound         = errors.New("agent not found")
	ErrReviewerNotFound      = errors.New("reviewer not found")
)

type TaskService struct {
	taskRepo    *repository.TaskRepository
	agentRepo   *repository.AgentRepository
	userRepo    *repository.UserRepository
	projectRepo *repository.ProjectRepository
}

func NewTaskService(
	taskRepo *repository.TaskRepository,
	agentRepo *repository.AgentRepository,
	userRepo *repository.UserRepository,
	projectRepo *repository.ProjectRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		agentRepo:   agentRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, projectID int, req *models.CreateTaskRequest, creatorID int, creatorType string) (*models.TaskWithDetails, error) {
	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Verify assigned agent if provided
	if req.AssignedAgentID != nil {
		agent, err := s.agentRepo.GetByID(ctx, *req.AssignedAgentID)
		if err != nil || agent.ProjectID != projectID {
			return nil, ErrAgentNotFound
		}
	}

	// Verify reviewer if provided
	if req.ReviewerID != nil {
		_, err := s.userRepo.GetByID(ctx, *req.ReviewerID)
		if err != nil {
			return nil, ErrReviewerNotFound
		}
	}

	task := &models.Task{
		ProjectID:       projectID,
		Title:           req.Title,
		Description:     req.Description,
		Status:          models.TaskStatusOpen,
		Priority:        req.Priority,
		AssignedAgentID: req.AssignedAgentID,
		ReviewerID:      req.ReviewerID,
		CreatedBy:       creatorID,
		CreatorType:     creatorType,
		Tags:            req.Tags,
	}

	if task.Priority == "" {
		task.Priority = models.TaskPriorityMedium
	}

	if task.Tags == nil {
		task.Tags = []string{}
	}

	err = s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Log activity
	activity := &models.TaskActivity{
		TaskID:    task.ID,
		ActorID:   creatorID,
		ActorType: creatorType,
		Action:    models.TaskActionCreated,
		Message:   fmt.Sprintf("Task '%s' created", task.Title),
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, task.ID)
}

// GetAllTasks retrieves all tasks with optional filters
func (s *TaskService) GetAllTasks(ctx context.Context, filters *models.TaskFilters) (*models.TaskListResponse, error) {
	tasks, err := s.taskRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	if tasks == nil {
		tasks = []models.TaskWithDetails{}
	}

	return &models.TaskListResponse{
		Tasks:      tasks,
		TotalCount: len(tasks),
	}, nil
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetTasksByProjectID retrieves all tasks for a project
func (s *TaskService) GetTasksByProjectID(ctx context.Context, projectID int) (*models.TaskListResponse, error) {
	tasks, err := s.taskRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	if tasks == nil {
		tasks = []models.TaskWithDetails{}
	}

	return &models.TaskListResponse{
		Tasks:      tasks,
		TotalCount: len(tasks),
	}, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(ctx context.Context, id int, req *models.UpdateTaskRequest, actorID int, actorType string) (*models.TaskWithDetails, error) {
	// Check if task exists
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	_, err = s.taskRepo.UpdatePartial(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return s.taskRepo.GetByIDWithDetails(ctx, id)
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTaskNotFound
	}

	return s.taskRepo.Delete(ctx, id)
}

// AssignAgent assigns an agent to a task
func (s *TaskService) AssignAgent(ctx context.Context, taskID int, agentID int, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Verify agent exists and belongs to same project
	agent, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil || agent.ProjectID != task.ProjectID {
		return nil, ErrAgentNotFound
	}

	err = s.taskRepo.AssignAgent(ctx, taskID, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign agent: %w", err)
	}

	// Log activity
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionAssigned,
		NewValue:  &agent.Name,
		Message:   fmt.Sprintf("Agent '%s' assigned to task", agent.Name),
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// SetReviewer sets the reviewer for a task
func (s *TaskService) SetReviewer(ctx context.Context, taskID int, reviewerID int, actorID int, actorType string) (*models.TaskWithDetails, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Verify reviewer exists
	reviewer, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil {
		return nil, ErrReviewerNotFound
	}

	err = s.taskRepo.SetReviewer(ctx, taskID, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to set reviewer: %w", err)
	}

	// Log activity
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionReviewerSet,
		NewValue:  &reviewer.Username,
		Message:   fmt.Sprintf("Reviewer set to '%s'", reviewer.Username),
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// StartTask moves task to in_progress status
func (s *TaskService) StartTask(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Validate status transition: only open -> in_progress
	if task.Status != models.TaskStatusOpen {
		return nil, ErrInvalidStatusTransition
	}

	oldStatus := task.Status
	newStatus := models.TaskStatusInProgress

	err = s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Log activity
	if message == "" {
		message = "Task started"
	}
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionStatusChanged,
		OldValue:  &oldStatus,
		NewValue:  &newStatus,
		Message:   message,
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// CompleteTask moves task to done status
func (s *TaskService) CompleteTask(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Validate status transition: in_progress -> done
	if task.Status != models.TaskStatusInProgress {
		return nil, ErrInvalidStatusTransition
	}

	oldStatus := task.Status
	newStatus := models.TaskStatusDone

	err = s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Log activity
	if message == "" {
		message = "Task completed"
	}
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionStatusChanged,
		OldValue:  &oldStatus,
		NewValue:  &newStatus,
		Message:   message,
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// CloseTask moves task to closed status (after review)
func (s *TaskService) CloseTask(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Validate status transition: done -> closed
	if task.Status != models.TaskStatusDone {
		return nil, ErrInvalidStatusTransition
	}

	oldStatus := task.Status
	newStatus := models.TaskStatusClosed

	err = s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Log activity
	if message == "" {
		message = "Task closed after review"
	}
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionStatusChanged,
		OldValue:  &oldStatus,
		NewValue:  &newStatus,
		Message:   message,
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// WontDoTask moves task to wont_do status
func (s *TaskService) WontDoTask(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Can mark as won't do from any status except already wont_do
	if task.Status == models.TaskStatusWontDo {
		return nil, ErrInvalidStatusTransition
	}

	oldStatus := task.Status
	newStatus := models.TaskStatusWontDo

	err = s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Log activity
	if message == "" {
		message = "Task marked as won't do"
	}
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionStatusChanged,
		OldValue:  &oldStatus,
		NewValue:  &newStatus,
		Message:   message,
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// ReopenTask moves task back to open status
func (s *TaskService) ReopenTask(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskWithDetails, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Can reopen from closed, done, or wont_do
	if task.Status != models.TaskStatusClosed && task.Status != models.TaskStatusDone && task.Status != models.TaskStatusWontDo {
		return nil, ErrInvalidStatusTransition
	}

	oldStatus := task.Status
	newStatus := models.TaskStatusOpen

	err = s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Log activity
	if message == "" {
		message = "Task reopened"
	}
	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionStatusChanged,
		OldValue:  &oldStatus,
		NewValue:  &newStatus,
		Message:   message,
	}
	_ = s.taskRepo.CreateActivity(ctx, activity)

	return s.taskRepo.GetByIDWithDetails(ctx, taskID)
}

// GetTaskActivities retrieves all activities for a task
func (s *TaskService) GetTaskActivities(ctx context.Context, taskID int) (*models.TaskActivityListResponse, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	activities, err := s.taskRepo.GetActivitiesByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	if activities == nil {
		activities = []models.TaskActivityWithDetails{}
	}

	return &models.TaskActivityListResponse{
		Activities: activities,
		TotalCount: len(activities),
	}, nil
}

// AddProgress adds a progress update to a task
func (s *TaskService) AddProgress(ctx context.Context, taskID int, message string, actorID int, actorType string) (*models.TaskActivityWithDetails, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	activity := &models.TaskActivity{
		TaskID:    taskID,
		ActorID:   actorID,
		ActorType: actorType,
		Action:    models.TaskActionProgress,
		Message:   message,
	}

	err = s.taskRepo.CreateActivity(ctx, activity)
	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	// Return with actor name
	activities, err := s.taskRepo.GetActivitiesByTaskID(ctx, taskID)
	if err != nil || len(activities) == 0 {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return &activities[0], nil
}
