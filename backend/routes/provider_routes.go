package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/gin-gonic/gin"
)

func setupProviderRoutes(r *gin.RouterGroup, providerHandler *handler.ProviderHandler) {
	providers := r.Group("/providers")
	{
		providers.GET("", providerHandler.GetAllProviders)
		providers.GET("/:name", providerHandler.GetProviderByName)
		providers.GET("/:name/models", providerHandler.GetProviderModels)
	}
}
