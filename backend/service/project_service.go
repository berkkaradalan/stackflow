package service

import (
	"context"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/repository/postgres"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *models.CreateProjectRequest, userID int) (*models.Project, error) {
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
		CreatedBy:   userID,
	}

	if req.Status != "" {
		project.Status = req.Status
	}

	err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) GetAllProjects(ctx context.Context) (*models.ProjectListResponse, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	return &models.ProjectListResponse{
		Projects:   projects,
		TotalCount: len(projects),
	}, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id int, req *models.UpdateProjectRequest) (*models.Project, error) {
	// First check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// Perform partial update
	updatedProject, err := s.projectRepo.UpdatePartial(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return updatedProject, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	err = s.projectRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (s *ProjectService) GetProjectStats(ctx context.Context, id int) (*models.ProjectStats, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	stats, err := s.projectRepo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project stats: %w", err)
	}

	return stats, nil
}
