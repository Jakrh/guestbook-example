package handler

import (
	"context"
	"errors"
	"guestbook-example/internal/api/model"
	"guestbook-example/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageService interface {
	Get(context.Context, int64) (*domain.Message, error)
	GetAll(context.Context) ([]*domain.Message, error)
	Create(context.Context, *domain.Message) (int64, error)
	Update(context.Context, *domain.Message) error
	Delete(context.Context, int64) error
}

// MessageHandler is the handler for message
type MessageHandler struct {
	logger         *slog.Logger
	messageService MessageService
}

// NewMessageHandler returns a new MessageHandler
func NewMessageHandler(logger *slog.Logger, messageService MessageService) *MessageHandler {
	return &MessageHandler{
		logger:         logger,
		messageService: messageService,
	}
}

// Get returns a message
func (h *MessageHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty of invalid id"})
		return
	}

	entity, err := h.messageService.Get(c, id)
	if err != nil {
		h.logger.Error("failed to get message", slog.String("error", err.Error()))

		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	message := &model.GetMessageResponse{
		ID:      entity.ID,
		Author:  entity.Author,
		Content: entity.Message,
	}

	c.JSON(http.StatusOK, message)
}

// GetAll returns all messages
func (h *MessageHandler) GetAll(c *gin.Context) {
	entities, err := h.messageService.GetAll(c)
	if err != nil {
		h.logger.Error("failed to get all messages", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	messages := model.NewListMessagesResponse(entities)

	c.JSON(http.StatusOK, messages)
}

// Create creates a message
func (h *MessageHandler) Create(c *gin.Context) {
	var req model.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind json", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.messageService.Create(c, &domain.Message{
		Author:  req.Author,
		Message: req.Content,
	})
	if err != nil {
		h.logger.Error("failed to create message", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Update updates a message
func (h *MessageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty of invalid id"})
		return
	}

	var req model.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind json", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = id

	err = h.messageService.Update(c, req.ToEntity())
	if err != nil {
		h.logger.Error("failed to update message", slog.String("error", err.Error()))

		// TODO: determine error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": req.ID})
}

// Delete deletes a message
func (h *MessageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))

		// TODO: determine error
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty of invalid id"})
		return
	}

	err = h.messageService.Delete(c, id)
	if err != nil {
		h.logger.Error("failed to delete message", slog.String("error", err.Error()))

		// TODO: determine error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
