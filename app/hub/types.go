package hub

import (
	"chat/app/logging"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type sendable interface {
	serialize() []byte
}

type client struct {
	conn          *websocket.Conn

	ChatUsername  string `json:"username"`
	UsernameColor string `json:"username_color"`

	incoming chan sendable
}

func (c *client) read() {
	defer func() {
		deleteClientCh <- c
	}()

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			logging.Logger.Error(err)
			break
		}

		message, err := unmarshalMessage(rawMessage)
		if err != nil {
			logging.Logger.Error(err)
			continue
		}

		newMessageCh <- message
	}
}

func (c *client) write() {
	ticker := time.NewTicker(30 * time.Second)

	defer func() {
		deleteClientCh <- c
		ticker.Stop()
	}()

	for {
		select {
		case sendable, ok := <-c.incoming:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logging.Logger.Error(err)
				return
			}

			w.Write(sendable.serialize())

			if err := w.Close(); err != nil {
				logging.Logger.Error(err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logging.Logger.Error(err)
				return
			}
		}
	}
}

func (c *client) serialize() []byte {
	b, _ := json.Marshal(c)
	return b
}

type message struct {
	Text          string `json:"text,omitempty"`
	Img           string `json:"img,omitempty"`
	FromUser      string `json:"from_user"`
	UsernameColor string `json:"username_color"`
	Date          string `json:"date"`
}

func (m *message) serialize() []byte {
	b, _ := json.Marshal(m)
	return b
}