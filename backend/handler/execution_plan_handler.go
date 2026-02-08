package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/service"
	"github.com/gin-gonic/gin"
)

type ExecutionPlanHandler struct {
	planService *service.ExecutionPlanService
}

func NewExecutionPlanHandler(planService *service.ExecutionPlanService) *ExecutionPlanHandler {
	return &ExecutionPlanHandler{
		planService: planService,
	}
}

// getActorInfo extracts actor ID and type from context
func (h *ExecutionPlanHandler) getActorInfo(c *gin.Context) (int, string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, "", false
	}
	return userID.(int), models.CreatorTypeUser, true
}

// --- Execution Plan Endpoints ---

// CreatePlan handles POST /api/projects/:id/execution-plan
func (h *ExecutionPlanHandler) CreatePlan(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.CreateExecutionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	plan, err := h.planService.CreatePlan(ctx, projectID, &req, actorID, actorType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create execution plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// GetActivePlan handles GET /api/projects/:id/execution-plan
func (h *ExecutionPlanHandler) GetActivePlan(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	plan, err := h.planService.GetActivePlan(ctx, projectID)
	if err != nil {
		if errors.Is(err, service.ErrNoActivePlan) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active execution plan found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch execution plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetAllPlans handles GET /api/projects/:id/execution-plans
func (h *ExecutionPlanHandler) GetAllPlans(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	plans, err := h.planService.GetPlansByProject(ctx, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch execution plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// UpdatePlan handles PUT /api/projects/:id/execution-plan
func (h *ExecutionPlanHandler) UpdatePlan(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.UpdateExecutionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	plan, err := h.planService.UpdatePlan(ctx, projectID, &req, actorID, actorType)
	if err != nil {
		if errors.Is(err, service.ErrNoActivePlan) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active execution plan found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update execution plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// --- Agent Task Flow Endpoints ---

// GetNextTask handles GET /api/agents/:id/next-task
func (h *ExecutionPlanHandler) GetNextTask(c *gin.Context) {
	ctx := c.Request.Context()

	agentIDParam := c.Param("id")
	agentID, err := strconv.Atoi(agentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	response, err := h.planService.GetNextTask(ctx, agentID)
	if err != nil {
		if errors.Is(err, service.ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch next task"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// TaskComplete handles POST /api/agents/:id/task-complete
func (h *ExecutionPlanHandler) TaskComplete(c *gin.Context) {
	ctx := c.Request.Context()

	agentIDParam := c.Param("id")
	agentID, err := strconv.Atoi(agentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	var req models.TaskCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.planService.CompleteTask(ctx, agentID, &req)
	if err != nil {
		if errors.Is(err, service.ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		if errors.Is(err, service.ErrAssignmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active assignment found for this agent and task"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// GetAgentContext handles GET /api/agents/:id/context
func (h *ExecutionPlanHandler) GetAgentContext(c *gin.Context) {
	ctx := c.Request.Context()

	agentIDParam := c.Param("id")
	agentID, err := strconv.Atoi(agentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	context, err := h.planService.GetAgentContext(ctx, agentID)
	if err != nil {
		if errors.Is(err, service.ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		if errors.Is(err, service.ErrNoActivePlan) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active execution plan for agent's project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agent context"})
		return
	}

	c.JSON(http.StatusOK, context)
}

// --- Report Endpoints ---

// GetDailyReport handles GET /api/projects/:id/reports/daily
func (h *ExecutionPlanHandler) GetDailyReport(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	reports, err := h.planService.GetReportsByType(ctx, projectID, models.ReportTypeDaily)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch daily reports"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetWeeklyReport handles GET /api/projects/:id/reports/weekly
func (h *ExecutionPlanHandler) GetWeeklyReport(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	reports, err := h.planService.GetReportsByType(ctx, projectID, models.ReportTypeWeekly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly reports"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GenerateReport handles POST /api/projects/:id/reports/generate
func (h *ExecutionPlanHandler) GenerateReport(c *gin.Context) {
	ctx := c.Request.Context()

	projectIDParam := c.Param("id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, actorType, ok := h.getActorInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	report, err := h.planService.GenerateReport(ctx, projectID, &req, actorID, actorType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	c.JSON(http.StatusCreated, report)
}
