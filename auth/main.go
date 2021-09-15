package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/todo-app/api/router"
	"github.com/todo-app/internal"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/pkg/config"
	"github.com/todo-app/pkg/exithandler"
	"github.com/todo-app/pkg/logger"
	"github.com/todo-app/pkg/server"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Info.Panic("failed to load env vars")
	}
	cfg := config.Get()

	db, err := internal.GetDataStore(cfg.GetDBConnStr())
	if err != nil {
		logger.Error.Fatal(err.Error())

	}
	logger.Info.Print("successfully connected to database")

	app, err := application.BootstrapApp(db, cfg)

	if err != nil {
		logger.Error.Fatalf("Failed bootstrapping app -- error: %s", err.Error())
	}

	srv := server.Get().
		WithAddr(app.Confg.GetAPIPort()).
		WithRouter(router.Get(app)).
		WithErrorLogger(logger.Error)

	go func() {
		logger.Info.Printf("Starting server on port %s", app.Confg.GetAPIPort())

		if err := srv.Listen(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		if err := srv.Close(); err != nil {
			log.Println(err.Error())
		}

		if err := app.CloseDBConn(); err != nil {
			logger.Error.Println(err.Error())
		}
	})

}
