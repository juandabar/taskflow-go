package container

import (
	"database/sql"

	sqliteRepo "github.com/juandabar/taskflow-go/internal/adapter/driven/persistence/sqlite/repository"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/handler"
	"github.com/juandabar/taskflow-go/internal/domain/usecase/auth"
	"github.com/juandabar/taskflow-go/internal/domain/usecase/user"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
}

func NewContainer(db *sql.DB, jwtSecret string) *Container {
	userRepo := sqliteRepo.NewSQLiteUserRepository(db)

	registerUseCase := auth.NewRegisterUserUseCase(userRepo)
	loginUseCase := auth.NewLoginUserUseCase(userRepo, jwtSecret)

	authHandler := handler.NewAuthHandler(registerUseCase, loginUseCase)

	getUserByIdUseCase := user.NewGetUserByIdUseCase(userRepo)

	userHandler := handler.NewUserHandler(getUserByIdUseCase)

	return &Container{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
}
