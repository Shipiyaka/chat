package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logger  *logrus.Logger
	logFile *os.File
)

func init() {
	var err error
	logFile, err = os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %s", err.Error())
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	logger.Info("logging started")
}

func main() {
	defer logFile.Close()

	srv := getServerInstance()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	go eventHandler()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}
}
