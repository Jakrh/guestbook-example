package repository

import (
	"guestbook-example/internal/domain"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Author  string `gorm:"not null"`
	Message string `gorm:"not null"`
}

func (m *Message) ToEntity() *domain.Message {
	return &domain.Message{
		Author:  m.Author,
		Message: m.Message,
	}
}

type Messages []*Message

func (ms Messages) ToEntity() []*domain.Message {
	entities := make([]*domain.Message, len(ms))
	for i, m := range ms {
		entities[i] = m.ToEntity()
	}
	return entities
}
