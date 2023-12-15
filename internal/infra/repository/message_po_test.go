package repository

import (
	"guestbook-example/internal/domain"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestMessage_ToEntity(t *testing.T) {
	tests := []struct {
		name string
		m    *Message
		want *domain.Message
	}{
		{
			name: "success",
			m: &Message{
				Model: gorm.Model{
					ID: 1,
				},
				Author:  "Arthur Morgan",
				Message: "Hey, Dutch!",
			},
			want: &domain.Message{
				ID:      1,
				Author:  "Arthur Morgan",
				Message: "Hey, Dutch!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ToEntity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.ToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessages_ToEntity(t *testing.T) {
	tests := []struct {
		name string
		ms   Messages
		want []*domain.Message
	}{
		{
			name: "success",
			ms: Messages{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
				{
					Model: gorm.Model{
						ID: 2,
					},
					Author:  "Dutch van der Linde",
					Message: "I have a plan!",
				},
			},
			want: []*domain.Message{
				{
					ID:      1,
					Author:  "Arthur Morgan",
					Message: "Hey, Dutch!",
				},
				{
					ID:      2,
					Author:  "Dutch van der Linde",
					Message: "I have a plan!",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ToEntity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Messages.ToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
