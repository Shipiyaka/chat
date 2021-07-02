package hub

import (
	"chat/app/logging"

	"github.com/gorilla/websocket"
)

var (
	chatParticipants = make(map[string]*client)

	NewConnCh = make(chan *websocket.Conn, 1)

	deleteClientCh = make(chan *client, 1)

	newMessageCh = make(chan *message, 1)
)

func EventHandler() {
	for {
		select {
		case newConn := <-NewConnCh:
			randomUsername := randomString()

			client := &client{conn: newConn, ChatUsername: randomUsername, UsernameColor: randomColor(), incoming: make(chan sendable)}

			go client.read()
			go client.write()

			err := sendOldMessages(client)
			if err != nil {
				logging.Logger.Error(err)
			}

			client.incoming <- client
			chatParticipants[randomUsername] = client

			newConn.SetCloseHandler(func(code int, text string) error {
				logging.Logger.Infof("Closed by %s. Code: %d, text: %s", randomUsername, code, text)

				deleteClientCh <- client

				return nil
			})
		case clientToRemove := <-deleteClientCh:
			clientToRemove.conn.Close()
			delete(chatParticipants, clientToRemove.ChatUsername)
		case newMessage := <-newMessageCh:
			for username, client := range chatParticipants {
				if username != newMessage.FromUser {
					client.incoming <- newMessage
				}
			}

			err := saveMessage(newMessage)
			if err != nil {
				logging.Logger.Error(err)
			}
		}
	}
}