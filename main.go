package main

import (
	"embed"
	"fmt"
	"guestbook-example/internal/api"
	"guestbook-example/internal/api/handler"
	"guestbook-example/internal/infra/repository"
	"guestbook-example/internal/service"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

//go:embed static/*
var embeddedFiles embed.FS

// GORM with glebarez/sqlite
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&repository.Message{})

	return db, err
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&repository.Message{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := initDB()
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	err = migrateDB(db)
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	staticFiles, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	messageRepo := repository.NewMessageRepo(logger, db)
	messageService := service.NewMessageService(logger, messageRepo)
	messageHandler := handler.NewMessageHandler(logger, messageService)
	staticFileHandler := handler.NewStaticFileHandler(logger, http.FS(staticFiles))

	router := api.SetupRouter(messageHandler, staticFileHandler)

	router.Run(":8080")

}
