//go:generate mockery --name MessageRepo --output ./mocks/

package service

import (
	"context"
	"fmt"
	"guestbook-example/internal/domain"
	"log/slog"
)

type MessageRepo interface {
	Get(context.Context, int64) (*domain.Message, error)
	GetAll(context.Context) ([]*domain.Message, error)
	Create(context.Context, *domain.Message) (int64, error)
	Update(context.Context, *domain.Message) error
	Delete(context.Context, int64) error
}

// MessageService is the interface that provides message methods.
type MessageService struct {
	logger      *slog.Logger
	messageRepo MessageRepo
}

// NewMessageService returns a new MessageService instance.
func NewMessageService(logger *slog.Logger, messageRepo MessageRepo) *MessageService {
	return &MessageService{
		logger:      logger,
		messageRepo: messageRepo,
	}
}

// Get returns a message.
func (s *MessageService) Get(ctx context.Context, id int64) (*domain.Message, error) {
	msg, err := s.messageRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return msg, nil
}

// GetAll returns all messages.
func (s *MessageService) GetAll(ctx context.Context) ([]*domain.Message, error) {
	msgs, err := s.messageRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all messages: %w", err)
	}

	return msgs, nil
}

// Create creates a message.
func (s *MessageService) Create(ctx context.Context, message *domain.Message) (int64, error) {
	id, err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return id, nil
}

// Update updates a message.
func (s *MessageService) Update(ctx context.Context, message *domain.Message) error {
	if err := s.messageRepo.Update(ctx, message); err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

// Delete deletes a message.
func (s *MessageService) Delete(ctx context.Context, id int64) error {
	if err := s.messageRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}
