package main

import (
	"encoding/json"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	chatParticipants = make(map[string]net.Conn)
	newParticipantCh = make(chan net.Conn, 10)
	deleteConnCh     = make(chan net.Conn, 10)
)

func handleMessagesFromConn(conn net.Conn) {
	for {
		msgInBytes, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			deleteConnCh <- conn
			logger.Error(err)
			break
		}

		logger.Infof("New message: %s", string(msgInBytes))

		msg, err := unmarshalMessage(msgInBytes)
		if err != nil {
			logger.Error(err)
			continue
		}

		for username, conn := range chatParticipants {
			if username != msg.FromUser {
				err = wsutil.WriteServerMessage(conn, op, msgInBytes)
				if err != nil {
					deleteConnCh <- conn
					logger.Error(err)
					break
				}
			}
		}
	}
}

func eventHandler() {
	for {
		select {
		case newConn := <-newParticipantCh:
			username := randomString()
			chatParticipants[username] = newConn

			b, _ := json.Marshal(map[string]string{"username": username})
			err := wsutil.WriteServerMessage(newConn, ws.OpText, b)
			if err != nil {
				deleteConnCh <- newConn
				logger.Error(err)
				continue
			}

			go handleMessagesFromConn(chatParticipants[username])
		case connToDelete := <-deleteConnCh:
			connToDelete.Close()
			for username, conn := range chatParticipants {
				if conn == connToDelete {
					delete(chatParticipants, username)
					break
				}
			}
		}
	}
}
