package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gobwas/ws"
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
		log.Fatalf("Error opening log file: %s", err.Error())
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	logger.Info("Logging started")
}

func main() {
	defer logFile.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", chatPage)
	http.HandleFunc("/ws", upgradeToWebSocketConn)

	go eventHandler()

	logger.Panic(http.ListenAndServe(":12345", nil))
}

func upgradeToWebSocketConn(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("New connection")

	newParticipantCh <- conn
}

func chatPage(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Error while parsing html file", http.StatusInternalServerError)
	}
}
