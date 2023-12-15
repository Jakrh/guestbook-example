package repository

import (
	"context"
	"errors"
	"fmt"
	"guestbook-example/internal/domain"
	"log/slog"

	"gorm.io/gorm"
)

type MessageRepo struct {
	logger *slog.Logger
	db     *gorm.DB
}

func NewMessageRepo(logger *slog.Logger, db *gorm.DB) *MessageRepo {
	return &MessageRepo{
		logger: logger,
		db:     db,
	}
}

func (r *MessageRepo) Create(ctx context.Context, m *domain.Message) (int64, error) {
	po := &Message{
		Author:  m.Author,
		Message: m.Message,
	}

	tx := r.db.Create(po)
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to create message from repository: %w", tx.Error)
	}

	return int64(po.ID), nil
}

func (r *MessageRepo) Get(ctx context.Context, id int64) (*domain.Message, error) {
	var m Message
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(domain.ErrNotFound, err)
		}
		return nil, err
	}

	return m.ToEntity(), nil
}

func (r *MessageRepo) GetAll(ctx context.Context) ([]*domain.Message, error) {
	var ms *Messages
	if err := r.db.Find(&ms).Error; err != nil {
		return nil, err
	}

	return ms.ToEntity(), nil
}

func (r *MessageRepo) Update(ctx context.Context, m *domain.Message) error {
	po := &Message{
		Model: gorm.Model{
			ID: uint(m.ID),
		},
		Author:  m.Author,
		Message: m.Message,
	}

	return r.db.Save(po).Error
}

func (r *MessageRepo) Delete(ctx context.Context, id int64) error {
	return r.db.Delete(&Message{}, id).Error
}
