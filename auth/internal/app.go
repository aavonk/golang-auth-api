package internal

import (
	"github.com/todo-app/internal/config"
	"github.com/todo-app/internal/logger"
)

type App struct {
	DataStore *DataStore
	Confg     *config.Confg
}

func BootstrapApp() (*App, error) {
	cfg := config.Get()

	db, err := GetDataStore(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	logger.Info.Print("successfully connected to database")
	return &App{
		DataStore: db,
		Confg:     cfg,
	}, nil
}
