package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/juandabar/taskflow-go/internal/infrastructure/config"
	"github.com/juandabar/taskflow-go/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
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

	ctx := context.Background()

	var count int
	db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM users
		WHERE email = ?
	`, "admin@taskflow.com").Scan(&count)
	if count > 0 {
		log.Println("admin user already exists, skipping")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("Betplay2026*"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		"00000000-0000-0000-0000-000000000001",
		"Admin",
		"admin@taskflow.com",
		string(hash),
		"ADMIN",
		time.Now().UTC().String(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("admin user created: admin@taskflow.com / Betplay2026*")
}
