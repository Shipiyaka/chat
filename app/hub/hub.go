package hub

import (
	"chat/app/db"
	"chat/app/logging"

	"github.com/gorilla/websocket"
)

var (
	chatParticipants = make(map[string]*client)

	NewConnCh = make(chan *websocket.Conn, 1)

	deleteClientCh = make(chan *client, 1)

	newMessageCh = make(chan message, 1)
)

func EventHandler() {
	for {
		select {
		case newConn := <-NewConnCh:
			randomUsername := randomString()

			client := &client{conn: newConn, ChatUsername: randomUsername, incoming: make(chan sendable)}

			go client.read()
			go client.write()

			history := make([]db.Message, 0)
			err := db.ReturnValues(map[string]interface{}{}, &history)
			if err == nil {
				for _, oldMessage := range history {
					var forSending message

					forSending.FromUser = oldMessage.FromUser
					forSending.Date = oldMessage.Date

					if oldMessage.ContentType == "image" {
						forSending.Img = oldMessage.Content
					} else if oldMessage.ContentType == "text" {
						forSending.Text = oldMessage.Content
					}

					client.incoming <- &forSending
				}
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
					client.incoming <- &newMessage
				}
			}

			messageContent := newMessage.Text
			contentType := "text"
			if messageContent == "" {
				messageContent = newMessage.Img
				contentType = "image"
			}

			err := db.Insert(&db.Message{
				Content:     messageContent,
				ContentType: contentType,
				FromUser:    newMessage.FromUser,
				Date:        newMessage.Date,
			})
			if err != nil {
				logging.Logger.Error(err)
			}
		}
	}
}
