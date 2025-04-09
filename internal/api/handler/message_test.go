package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"guestbook-example/internal/api/handler/mocks"
	"guestbook-example/internal/domain"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageHandler_Get(t *testing.T) {
	gin.DefaultWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name           string
		messageService MessageService
		id             string
		expectedStatus int
	}{
		{
			name: "success",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Get", mock.Anything, mock.AnythingOfType("int64")).Return(&domain.Message{}, nil)
				return mockService
			}(),
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name: "failed to parse id",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				return mockService
			}(),
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to get message",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Get", mock.Anything, mock.AnythingOfType("int64")).Return(nil, domain.ErrNotFound)
				return mockService
			}(),
			id:             "1",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMessageHandler(logger, tt.messageService)

			router := gin.Default()
			router.GET("/messages/:id", handler.Get)

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/messages/%s", tt.id), nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestMessageHandler_GetAll(t *testing.T) {
	gin.DefaultWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name           string
		messageService MessageService
		expectedStatus int
	}{
		{
			name: "success",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("GetAll", mock.Anything).Return([]*domain.Message{}, nil)
				return mockService
			}(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "failed to get all messages",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("GetAll", mock.Anything).Return(nil, fmt.Errorf("failed to get all messages"))
				return mockService
			}(),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMessageHandler(logger, tt.messageService)

			router := gin.Default()
			router.GET("/messages", handler.GetAll)

			req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestMessageHandler_Create(t *testing.T) {
	gin.DefaultWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	type requestBody struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}
	tests := []struct {
		name           string
		messageService MessageService
		requestBody    requestBody
		expectedStatus int
	}{
		{
			name: "success",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(int64(1), nil)
				return mockService
			}(),
			requestBody: requestBody{
				Author:  "John Doe",
				Content: "Hello everybody!",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "failed to create message with empty request body",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(int64(1), nil)
				return mockService
			}(),
			requestBody:    requestBody{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to create message with empty author",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(int64(1), nil)
				return mockService
			}(),
			requestBody: requestBody{
				Author:  "",
				Content: "Hello everybody!",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to create message with empty content",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(int64(1), nil)
				return mockService
			}(),
			requestBody: requestBody{
				Author:  "John Doe",
				Content: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMessageHandler(logger, tt.messageService)

			router := gin.Default()
			router.POST("/messages", handler.Create)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to encode messageRequest: %v", err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/messages", &buf)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestMessageHandler_Update(t *testing.T) {
	gin.DefaultWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	type requestBody struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}
	tests := []struct {
		name           string
		messageService MessageService
		requestBody    requestBody
		expectedStatus int
	}{
		{
			name: "success",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Update", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(nil)
				return mockService
			}(),
			requestBody: requestBody{
				Author:  "John Doe",
				Content: "Hello everybody!",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failed to update message",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Update", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(fmt.Errorf("failed to update message"))
				return mockService
			}(),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMessageHandler(logger, tt.messageService)

			router := gin.Default()
			router.PUT("/messages/:id", handler.Update)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to encode messageRequest: %v", err)
			}

			req, _ := http.NewRequest(http.MethodPut, "/messages/1", &buf)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestMessageHandler_Delete(t *testing.T) {
	gin.DefaultWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name           string
		messageService MessageService
		id             string
		expectedStatus int
	}{
		{
			name: "success",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)
				return mockService
			}(),
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name: "failed to parse id",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				return mockService
			}(),
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to delete message",
			messageService: func() MessageService {
				mockService := new(mocks.MessageService)
				mockService.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(fmt.Errorf("failed to delete message"))
				return mockService
			}(),
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMessageHandler(logger, tt.messageService)

			router := gin.Default()
			router.DELETE("/messages/:id", handler.Delete)

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/messages/%s", tt.id), nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}
