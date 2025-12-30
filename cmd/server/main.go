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

	log.Fatal(http.ListenAndServe(":8080", r))
}
