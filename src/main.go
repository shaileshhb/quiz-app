package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shaileshhb/quiz/src/db"
	"github.com/shaileshhb/quiz/src/log"
	"github.com/shaileshhb/quiz/src/server"
)

func main() {
	logger := log.InitializeLogger()
	database := db.NewDatabase()
	ser := server.NewServer(logger, database)
	ser.InitializeRouter()

	ser.RegisterModuleRoutes()

	logger.Error().Err(ser.App.Listen(":8080")).Msg("")

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	os.Exit(0)
}
