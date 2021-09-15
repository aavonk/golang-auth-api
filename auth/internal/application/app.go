package application

import (
	"github.com/todo-app/internal"
	"github.com/todo-app/internal/mailer"
	"github.com/todo-app/internal/repositories"
	"github.com/todo-app/internal/services"
	"github.com/todo-app/pkg/config"
)

type App struct {
	dataStore       *internal.DataStore
	Confg           *config.Confg
	Mailer          mailer.Mailer
	UserRepository  repositories.UserRepositoryInterface
	TokenRepository repositories.TokenRepositoryInterface
	IdentityService services.IdentityServiceInterface
}

func BootstrapApp(db *internal.DataStore, cfg *config.Confg) (*App, error) {

	return &App{
		dataStore:       db,
		Confg:           cfg,
		UserRepository:  repositories.NewUserRepository(db.Client),
		TokenRepository: repositories.NewTokenRepository(db.Client),
		IdentityService: services.NewIdentityService(db.Client),
		Mailer: mailer.New(
			cfg.Smtp.Host,
			cfg.Smtp.Port,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Sender,
		),
	}, nil
}

func (a *App) CloseDBConn() error {
	return a.dataStore.Close()
}
