package main

import (
	"chat/app/db"
	"chat/app/hub"
	"chat/app/logging"
	"chat/app/server"
	"chat/config"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	cfg config.Config
)

func init() {
	logging.InitLogger()
	logging.Logger.Info("Logging started")

	if err := env.Parse(&cfg); err != nil {
		logging.Logger.Fatal(err)
	}

	if err := db.InitDBInstance(cfg.DBPath); err != nil {
		logging.Logger.Fatal(err)
	}
}

func main() {
	srv := server.GetServerInstance(cfg.ServerAddr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logging.Logger.Error(err)
		}
	}()

	go hub.EventHandler()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		logging.Logger.Error(err)
	}
}