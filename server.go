package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

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

func getServerInstance() *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", chatPage)
	mux.HandleFunc("/ws", upgradeToWebSocketConn)

	srv := &http.Server{
		Addr:         "localhost:12345",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	return srv
}
