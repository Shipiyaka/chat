package hub

import (
	"chat/app/logging"

	"github.com/gorilla/websocket"
)

var (
	chatParticipants = make(map[string]*client)

	NewConnCh = make(chan *websocket.Conn, 1)

	deleteClientCh = make(chan string, 1)

	newMessageCh = make(chan message, 1)
)

func EventHandler() {
	for {
		select {
		case newConn := <-NewConnCh:
			randomUsername := randomString()

			client := &client{conn: newConn, ChatUsername: randomUsername, incoming: make(chan sendable)}

			chatParticipants[randomUsername] = client

			go client.read()
			go client.write()

			client.incoming <- client

			newConn.SetCloseHandler(func(code int, text string) error {
				logging.Logger.Infof("Closed by %s. Code: %d, text: %s", randomUsername, code, text)

				deleteClientCh <- randomUsername

				return nil
			})
		case clientToRemove := <-deleteClientCh:
			delete(chatParticipants, clientToRemove)
		case newMessage := <-newMessageCh:
			for username, client := range chatParticipants {
				if username != newMessage.FromUser {
					client.incoming <- &newMessage
				}
			}
		}
	}
}
