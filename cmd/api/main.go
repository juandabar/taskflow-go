package main

import (
	"log"
	"net/http"

	httpAdapter "github.com/juandabar/taskflow-go/internal/adapter/driving/http"
	"github.com/juandabar/taskflow-go/internal/infrastructure/config"
	"github.com/juandabar/taskflow-go/internal/infrastructure/container"
	"github.com/juandabar/taskflow-go/internal/infrastructure/database"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSQLiteConnection(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := container.NewContainer(db, cfg.JWTSecret)

	router := httpAdapter.NewRouter(c.AuthHandler)

	log.Printf("server runing on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
