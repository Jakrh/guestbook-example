package api

import (
	"guestbook-example/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(messageHandler *handler.MessageHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/messages", messageHandler.Create)
	router.GET("/messages", messageHandler.GetAll)
	router.GET("/messages/:id", messageHandler.Get)
	router.PUT("/messages/:id", messageHandler.Update)
	router.DELETE("/messages/:id", messageHandler.Delete)

	return router
}
