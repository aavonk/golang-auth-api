package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/todo-app/api/router"
	"github.com/todo-app/internal"
	"github.com/todo-app/internal/exithandler"
	"github.com/todo-app/internal/logger"
	"github.com/todo-app/internal/server"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Info.Panic("failed to load env vars")
	}

	app, err := internal.BootstrapApp()

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

		if err := app.DataStore.Close(); err != nil {
			logger.Error.Println(err.Error())
		}
	})

}
