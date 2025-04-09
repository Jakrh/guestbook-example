package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"guestbook-example/internal/api/handler/mocks"
)

func TestNewStaticFileHandler(t *testing.T) {
	// Arrange
	logger := slog.Default()
	mockFS := mocks.NewFileSystem(t)

	// Act
	handler := NewStaticFileHandler(logger, mockFS)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, logger, handler.logger)
	assert.Equal(t, mockFS, handler.staticFS)
}

func TestStaticFileHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		setupMock      func(fs *mocks.FileSystem)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successfully serve file",
			path: "/style.css",
			setupMock: func(fs *mocks.FileSystem) {
				mockFile := NewMockFile("body { color: red; }")
				mockFile.On("Close").Return(nil)
				mockFile.On("Stat").Return(mock.AnythingOfType("*os.fileStat"), nil)
				fs.On("Open", "style.css").Return(mockFile, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "body { color: red; }",
		},
		{
			name: "File not found, fallback to index.html",
			path: "/not-found.html",
			setupMock: func(fs *mocks.FileSystem) {
				fs.On("Open", "not-found.html").Return(nil, os.ErrNotExist)
				mockFile := NewMockFile("<html>Index Page</html>")
				mockFile.On("Close").Return(nil)
				mockFile.On("Stat").Return(mock.AnythingOfType("*os.fileStat"), nil)
				fs.On("Open", "index.html").Return(mockFile, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "<html>Index Page</html>",
		},
		{
			name: "Index.html not found",
			path: "/page.html",
			setupMock: func(fs *mocks.FileSystem) {
				fs.On("Open", "page.html").Return(nil, os.ErrNotExist)
				fs.On("Open", "index.html").Return(nil, os.ErrNotExist)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "index.html not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockFS := mocks.NewFileSystem(t)
			tt.setupMock(mockFS)

			handler := NewStaticFileHandler(slog.Default(), mockFS)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			c.Request = req

			// Act
			handler.Get(c)

			// Assert
			require.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

// MockFile implements http.File interface for testing
type MockFile struct {
	mock.Mock
	reader *strings.Reader
}

func NewMockFile(content string) *MockFile {
	return &MockFile{
		reader: strings.NewReader(content),
	}
}

func (m *MockFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *MockFile) Seek(offset int64, whence int) (int64, error) {
	return m.reader.Seek(offset, whence)
}

func (m *MockFile) Readdir(count int) ([]os.FileInfo, error) {
	args := m.Called(count)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

func (m *MockFile) Stat() (os.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(os.FileInfo), args.Error(1)
}
