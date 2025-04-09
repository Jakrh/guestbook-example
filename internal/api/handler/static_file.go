package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StaticFileHandler struct {
	logger   *slog.Logger
	staticFS http.FileSystem
}

// NewStaticFileHandler creates a new StaticFileHandler
func NewStaticFileHandler(logger *slog.Logger, staticFS http.FileSystem) *StaticFileHandler {
	return &StaticFileHandler{
		logger:   logger,
		staticFS: staticFS,
	}
}

func (h *StaticFileHandler) Get(c *gin.Context) {
	path := c.Request.URL.Path

	// Try to serve the requested file from the static file system
	if f, err := h.staticFS.Open(path[1:]); err == nil {
		defer f.Close()
		http.ServeContent(c.Writer, c.Request, path, time.Now(), f)
		return
	}

	// fallback to index.html
	f, err := h.staticFS.Open("index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "index.html not found")
		return
	}
	defer f.Close()

	http.ServeContent(c.Writer, c.Request, "index.html", time.Now(), f)
}
