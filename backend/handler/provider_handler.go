package handler

import (
	"net/http"

	"github.com/berkkaradalan/stackflow/service"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerService *service.ProviderService
}

func NewProviderHandler(providerService *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// GetAllProviders handles GET /api/providers
func (h *ProviderHandler) GetAllProviders(c *gin.Context) {
	providers := h.providerService.GetAllProviders()
	c.JSON(http.StatusOK, providers)
}

// GetProviderModels handles GET /api/providers/:name/models
func (h *ProviderHandler) GetProviderModels(c *gin.Context) {
	providerName := c.Param("name")

	models, err := h.providerService.GetProviderModels(providerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

// GetProviderByName handles GET /api/providers/:name
func (h *ProviderHandler) GetProviderByName(c *gin.Context) {
	providerName := c.Param("name")

	provider, err := h.providerService.GetProviderByName(providerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}
