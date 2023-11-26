package repository

import "guestbook-example/internal/domain"

type Message struct {
	ID      int64  `gorm:"primaryKey"`
	Author  string `gorm:"not null"`
	Message string `gorm:"not null"`
}

func (m *Message) ToEntity() *domain.Message {
	return &domain.Message{
		ID:      m.ID,
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
