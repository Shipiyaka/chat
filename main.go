package main

import (
	"html/template"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	chatUsers = make(map[net.Conn]string)
)

func upgradeToWebSocketConn(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Cannot upgrade connection to ws", http.StatusInternalServerError)
		return
	}

	logger.Info("New connection")

	go func() {
		defer conn.Close()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				logger.Error(err)
				break
			}
			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				logger.Error(err)
				break
			}
		}
	}()
}

func chatPage(w http.ResponseWriter, r *http.Request) {
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
	logger.Error(http.ListenAndServe(":12345", nil))
}
