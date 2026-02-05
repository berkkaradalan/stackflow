package service

import (
	"github.com/berkkaradalan/stackflow/config"
	"github.com/berkkaradalan/stackflow/models"
)

type ProviderService struct{}

func NewProviderService() *ProviderService {
	return &ProviderService{}
}

// GetAllProviders returns all available providers
func (s *ProviderService) GetAllProviders() *models.ProviderListResponse {
	providers := config.GetProviderRegistry()
	return &models.ProviderListResponse{
		Providers: providers,
	}
}

// GetProviderModels returns all models for a specific provider
func (s *ProviderService) GetProviderModels(providerName string) (*models.ModelListResponse, error) {
	provider := config.GetProviderByName(providerName)
	if provider == nil {
		return nil, &ProviderNotFoundError{Provider: providerName}
	}

	return &models.ModelListResponse{
		Provider: providerName,
		Models:   provider.Models,
	}, nil
}

// GetProviderByName returns a specific provider configuration
func (s *ProviderService) GetProviderByName(providerName string) (*models.ProviderConfig, error) {
	provider := config.GetProviderByName(providerName)
	if provider == nil {
		return nil, &ProviderNotFoundError{Provider: providerName}
	}
	return provider, nil
}

// ProviderNotFoundError represents a provider not found error
type ProviderNotFoundError struct {
	Provider string
}

func (e *ProviderNotFoundError) Error() string {
	return "provider not found: " + e.Provider
}
