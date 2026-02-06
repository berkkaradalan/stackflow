package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// getActorInfo extracts actor ID and type from context
func (h *TaskHandler) getActorInfo(c *gin.Context) (int, string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, "", false
	}
	// For now, all requests come from users (agents will use a different auth mechanism)
	return userID.(int), models.CreatorTypeUser, true
}

// GetAllTasks handles GET /api/tasks
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters for filtering
	var filters models.TaskFilters

	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := strconv.Atoi(projectID); err == nil {
			filters.ProjectID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	if priority := c.Query("priority"); priority != "" {
		filters.Priority = &priority
	}

	if agentID := c.Query("assigned_agent_id"); agentID != "" {
		if id, err := strconv.Atoi(agentID); err == nil {
			filters.AssignedAgentID = &id
		}
	}

	if reviewerID := c.Query("reviewer_id"); reviewerID != "" {
		if id, err := strconv.Atoi(reviewerID); err == nil {
			filters.ReviewerID = &id
		}
	}

	tasks, err := h.taskService.GetAllTasks(ctx, &filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTask handles POST /api/projects/:id/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.CreateTask(ctx, projectID, &req, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrAgentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Agent not found or does not belong to project"})
			return
		}
		if errors.Is(err, service.ErrReviewerNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reviewer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasksByProject handles GET /api/projects/:id/tasks
func (h *TaskHandler) GetTasksByProject(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tasks, err := h.taskService.GetTasksByProjectID(ctx, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID handles GET /api/tasks/:id
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskService.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask handles PUT /api/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.UpdateTask(ctx, id, &req, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /api/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = h.taskService.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// AssignAgent handles POST /api/tasks/:id/assign
func (h *TaskHandler) AssignAgent(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.AssignAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.AssignAgent(ctx, id, req.AgentID, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrAgentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Agent not found or does not belong to project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign agent"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// SetReviewer handles POST /api/tasks/:id/reviewer
func (h *TaskHandler) SetReviewer(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.SetReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.SetReviewer(ctx, id, req.ReviewerID, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrReviewerNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reviewer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set reviewer"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// StartTask handles POST /api/tasks/:id/start
func (h *TaskHandler) StartTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskStatusChangeRequest
	_ = c.ShouldBindJSON(&req) // Optional message

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.StartTask(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task can only be started from open status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CompleteTask handles POST /api/tasks/:id/done
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskStatusChangeRequest
	_ = c.ShouldBindJSON(&req)

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.CompleteTask(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task can only be completed from in_progress status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CloseTask handles POST /api/tasks/:id/close
func (h *TaskHandler) CloseTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskStatusChangeRequest
	_ = c.ShouldBindJSON(&req)

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.CloseTask(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task can only be closed from done status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// WontDoTask handles POST /api/tasks/:id/wontdo
func (h *TaskHandler) WontDoTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskStatusChangeRequest
	_ = c.ShouldBindJSON(&req)

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.WontDoTask(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task is already marked as won't do"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ReopenTask handles POST /api/tasks/:id/reopen
func (h *TaskHandler) ReopenTask(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskStatusChangeRequest
	_ = c.ShouldBindJSON(&req)

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task, err := h.taskService.ReopenTask(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task can only be reopened from closed, done, or wont_do status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reopen task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetTaskActivities handles GET /api/tasks/:id/activities
func (h *TaskHandler) GetTaskActivities(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	activities, err := h.taskService.GetTaskActivities(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// AddProgress handles POST /api/tasks/:id/activities
func (h *TaskHandler) AddProgress(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	activity, err := h.taskService.AddProgress(ctx, id, req.Message, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add progress"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}
