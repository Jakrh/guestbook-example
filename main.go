package main

import (
	"fmt"
	"guestbook-example/internal/api"
	"guestbook-example/internal/api/handler"
	"guestbook-example/internal/domain"
	"guestbook-example/internal/infra/repository"
	"guestbook-example/internal/service"
	"log/slog"
	"os"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

// GORM with modernc.org/sqlite
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&domain.Message{})

	return db, err
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Message{})
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

	messageRepo := repository.NewMessageRepo(logger, db)
	messageService := service.NewMessageService(logger, messageRepo)
	messageHandler := handler.NewMessageHandler(logger, messageService)

	router := api.SetupRouter(messageHandler)

	router.Run(":8080")

}
