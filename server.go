package main

import (
	"html/template"
	"net/http"

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

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", chatPage)
	http.HandleFunc("/ws", upgradeToWebSocketConn)

	go eventHandler()

	logger.Error(http.ListenAndServe(":12345", nil))
}
