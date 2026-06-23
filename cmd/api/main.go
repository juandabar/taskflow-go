package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	httpAdapter "github.com/juandabar/taskflow-go/internal/adapter/driving/http"
	"github.com/juandabar/taskflow-go/internal/infrastructure/config"
	"github.com/juandabar/taskflow-go/internal/infrastructure/container"
	"github.com/juandabar/taskflow-go/internal/infrastructure/database"
)

func main() {
	config.LoadEnvFile(".env")

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

	router := httpAdapter.NewRouter(c.AuthHandler, c.UserHandler, cfg.JWTSecret)

	log.Printf("server runing on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading .env file: %v", err)
	}

}
