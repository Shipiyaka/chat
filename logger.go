package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	f, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %s", err.Error())
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(io.MultiWriter(os.Stdout, f))
	
	logger.Info("Logging started")
}
