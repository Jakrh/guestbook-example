package model

import "guestbook-example/internal/domain"

type CreateMessageRequest struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func (r *CreateMessageRequest) ToEntity() *domain.Message {
	return &domain.Message{
		Author:  r.Author,
		Message: r.Content,
	}
}

type CreateMessageResponse struct {
	ID int64 `json:"id"`
}

type GetMessageResponse struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func NewGetMessageResponse(entity *domain.Message) *GetMessageResponse {
	return &GetMessageResponse{
		ID:      entity.ID,
		Author:  entity.Author,
		Content: entity.Message,
	}
}

type ListMessagesResponse struct {
	Messages []GetMessageResponse `json:"messages"`
}

func NewListMessagesResponse(entities []*domain.Message) *ListMessagesResponse {
	messages := make([]GetMessageResponse, len(entities))
	for i, entity := range entities {
		messages[i] = *NewGetMessageResponse(entity)
	}
	return &ListMessagesResponse{
		Messages: messages,
	}
}

type UpdateMessageRequest struct {
	ID      int64
	Author  string `json:"author"`
	Content string `json:"content"`
}

func (r *UpdateMessageRequest) ToEntity() *domain.Message {
	return &domain.Message{
		ID:      r.ID,
		Author:  r.Author,
		Message: r.Content,
	}
}

type UpdateMessageResponse struct {
	ID int64 `json:"id"`
}

type DeleteMessageRequest struct {
	ID int64 `json:"id"`
}
