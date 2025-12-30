package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vikramsinghpanwar/ludo-backend/internal/app/router"
	"github.com/vikramsinghpanwar/ludo-backend/internal/auth"
	"github.com/vikramsinghpanwar/ludo-backend/pkg/database"
)

func main() {
	log.Println("Starting Ludo backend...")
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Println("HTTP_PORT not set, defaulting to 8080")
		port = "8080"
	}
	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("database connected!!")
	}

	authRepo := auth.NewPostgresAuthRepo(db)

	authService := auth.NewService(authRepo)

	authHandler := auth.NewHandler(authService)

	r := router.NewRouter(&router.Dependencies{
		AuthHandler: authHandler,
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}
