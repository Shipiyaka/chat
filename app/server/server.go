package server

import (
	"chat/app/hub"
	"chat/app/logging"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  16384,
	WriteBufferSize: 16384,
}

func upgradeConnToWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Logger.Error(err)
		return
	}

	hub.NewConnCh <- conn
}

func chatPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error while parsing html file", http.StatusInternalServerError)
	}
}

func GetServerInstance(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", chatPage)
	mux.HandleFunc("/ws", upgradeConnToWS)

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	return srv
}
