package internal

import (
	"github.com/todo-app/pkg/config"
	"github.com/todo-app/pkg/logger"
)

type App struct {
	DataStore      *DataStore
	Confg          *config.Confg
	UserRepository UserRepositoryInterface
	// IdentityService service.IdentityServiceInterface
}

func BootstrapApp() (*App, error) {
	cfg := config.Get()

	db, err := GetDataStore(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	logger.Info.Print("successfully connected to database")
	return &App{
		DataStore:      db, // TODO: take out -- db should only be accessed in repositories
		Confg:          cfg,
		UserRepository: NewUserRepository(db.Client),
		// IdentityService: service.NewIdentityService(db.Client),
	}, nil
}
