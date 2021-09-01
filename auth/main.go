package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/todo-app/api/router"
	"github.com/todo-app/internal/exithandler"
	"github.com/todo-app/internal/logger"
	"github.com/todo-app/internal/server"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Info.Panic("failed to load env vars")
	}

	srv := server.Get().
		WithAddr(":" + os.Getenv("API_PORT")).
		WithRouter(router.Get()).
		WithErrorLogger(logger.Error)

	go func() {
		logger.Info.Printf("Starting server on port %s", os.Getenv("API_PORT"))

		if err := srv.Listen(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		// TODO:Close the DB Connection here
		if err := srv.Close(); err != nil {
			log.Println(err.Error())
		}
	})

}
