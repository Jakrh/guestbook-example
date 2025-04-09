package api

import (
	"github.com/gin-gonic/gin"
)

type MessageHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type StaticFileHandler interface {
	Get(c *gin.Context)
}

func SetupRouter(messageHandler MessageHandler, staticFileHandler StaticFileHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/messages", messageHandler.Create)
		api.GET("/messages", messageHandler.GetAll)
		api.GET("/messages/:id", messageHandler.Get)
		api.PUT("/messages/:id", messageHandler.Update)
		api.DELETE("/messages/:id", messageHandler.Delete)
	}

	router.NoRoute(staticFileHandler.Get)

	return router
}
