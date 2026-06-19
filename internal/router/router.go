package router

import (
	"github.com/ShamsiddinTS/subscriptions-service/internal/handler"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	subscriptionHandler *handler.SubscriptionHandler,
) *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	api := r.Group("/api")
	{
		api.POST("/subscriptions", subscriptionHandler.Create)
		api.GET("/subscriptions", subscriptionHandler.List)
		api.GET("/subscriptions/:id", subscriptionHandler.GetByID)
		api.PUT("/subscriptions/:id", subscriptionHandler.Update)
		api.DELETE("/subscriptions/:id", subscriptionHandler.Delete)
		api.GET("/subscriptions/total", subscriptionHandler.CalculateTotalCost)
	}

	return r
}
