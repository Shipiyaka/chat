package logging

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{}
}
