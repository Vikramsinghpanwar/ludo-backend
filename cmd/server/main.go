package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vikramsinghpanwar/ludo-backend/internal/app/router"
	"github.com/vikramsinghpanwar/ludo-backend/internal/auth"
	"github.com/vikramsinghpanwar/ludo-backend/internal/infra/sms"
	"github.com/vikramsinghpanwar/ludo-backend/pkg/database"

	"github.com/vikramsinghpanwar/ludo-backend/internal/config"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Println("HTTP_PORT not set, defaulting to 8080")
		port = "8080"
	}

	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Load()

	smsProvider := sms.NewBulkSMSLab(
		cfg.SMS.APIKey,
		cfg.SMS.UserID,
		cfg.SMS.Password,
		cfg.SMS.SenderID,
	)

	authRepo := auth.NewPostgresAuthRepo(db)

	authService := auth.NewService(authRepo, smsProvider)

	authHandler := auth.NewHandler(authService)

	r := router.NewRouter(&router.Dependencies{
		AuthHandler: authHandler,
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}
