package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"guestbook-example/internal/api/mocks"
)

// Helper function to perform a request and get response
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSetupRouter(t *testing.T) {
	// Set Gin to test mode and disable logs
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard

	// Create instances of the mocks
	mockMessageHandler := &mocks.MessageHandler{}
	mockStaticFileHandler := &mocks.StaticFileHandler{}

	// Table-driven mock setup
	mockSetups := []struct {
		handler    any
		methodName string
		returnCode int
	}{
		{mockMessageHandler, "Create", http.StatusCreated},
		{mockMessageHandler, "GetAll", http.StatusOK},
		{mockMessageHandler, "Get", http.StatusOK},
		{mockMessageHandler, "Update", http.StatusOK},
		{mockMessageHandler, "Delete", http.StatusNoContent},
		{mockStaticFileHandler, "Get", http.StatusOK},
	}

	// Configure all mocks using the table-driven approach
	for _, setup := range mockSetups {
		switch h := setup.handler.(type) {
		case *mocks.MessageHandler:
			h.On(setup.methodName, mock.Anything).Run(func(args mock.Arguments) {
				c := args.Get(0).(*gin.Context)
				c.Status(setup.returnCode)
			})
		case *mocks.StaticFileHandler:
			h.On(setup.methodName, mock.Anything).Run(func(args mock.Arguments) {
				c := args.Get(0).(*gin.Context)
				c.Status(setup.returnCode)
			})
		}
	}

	// Call SetupRouter with the mocks
	router := SetupRouter(mockMessageHandler, mockStaticFileHandler)

	// Table-driven test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		handlerMethod  string
		mockHandler    *mock.Mock
	}{
		{
			name:           "POST /api/v1/messages",
			method:         "POST",
			path:           "/api/v1/messages",
			expectedStatus: http.StatusCreated,
			handlerMethod:  "Create",
			mockHandler:    &mockMessageHandler.Mock,
		},
		{
			name:           "GET /api/v1/messages",
			method:         "GET",
			path:           "/api/v1/messages",
			expectedStatus: http.StatusOK,
			handlerMethod:  "GetAll",
			mockHandler:    &mockMessageHandler.Mock,
		},
		{
			name:           "GET /api/v1/messages/123",
			method:         "GET",
			path:           "/api/v1/messages/123",
			expectedStatus: http.StatusOK,
			handlerMethod:  "Get",
			mockHandler:    &mockMessageHandler.Mock,
		},
		{
			name:           "PUT /api/v1/messages/123",
			method:         "PUT",
			path:           "/api/v1/messages/123",
			expectedStatus: http.StatusOK,
			handlerMethod:  "Update",
			mockHandler:    &mockMessageHandler.Mock,
		},
		{
			name:           "DELETE /api/v1/messages/123",
			method:         "DELETE",
			path:           "/api/v1/messages/123",
			expectedStatus: http.StatusNoContent,
			handlerMethod:  "Delete",
			mockHandler:    &mockMessageHandler.Mock,
		},
		{
			name:           "NoRoute handler",
			method:         "GET",
			path:           "/non-existent-route",
			expectedStatus: http.StatusOK,
			handlerMethod:  "Get",
			mockHandler:    &mockStaticFileHandler.Mock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := performRequest(router, tc.method, tc.path)
			assert.Equal(t, tc.expectedStatus, w.Code)
			tc.mockHandler.AssertCalled(t, tc.handlerMethod, mock.Anything)
		})
	}

	// Verify all expectations were met
	mockMessageHandler.AssertExpectations(t)
	mockStaticFileHandler.AssertExpectations(t)
}
