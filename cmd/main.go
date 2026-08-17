package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"avito/internal/config"
	"avito/internal/handler"
	"avito/internal/middleware"
	"avito/internal/repository"
	"avito/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error in getting .env: %s\n", err.Error())
	}
	dbConnectionLink := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dbConnectionLink)
	if err != nil {
		log.Printf("Error in openning database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error in openning database: %s", err)
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := handler.NewService(service)

	http.HandleFunc("POST /register", handler.Register)
	http.HandleFunc("POST /login", handler.Login)
	http.HandleFunc("GET /dummyLogin/", handler.DummyHandler)
	http.HandleFunc("POST /house/create/", middleware.AuthMiddleware(handler.CreateHouse))
	http.HandleFunc("POST /flat/create/", middleware.AuthMiddleware(handler.CreateFlat))
	http.HandleFunc("POST /flat/update/", middleware.AuthMiddleware(handler.ModUpdate))
	http.HandleFunc("GET /house/{id}", middleware.AuthMiddleware(handler.GetFlatsByHouseId))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error in booting server: %s\n", err.Error())
	}
}
