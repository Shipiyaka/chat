package hub

import (
	"chat/app/db"
	"encoding/json"
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStringLen = 8
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomString() string {
	b := make([]byte, randomStringLen)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func unmarshalMessage(b []byte) (*message, error) {
	m := new(message)

	err := json.Unmarshal(b, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func sendOldMessages(c *client) error {
	oldMessages := make([]db.Message, 0)
	err := db.ReturnValues(map[string]interface{}{}, &oldMessages)
	if err != nil {
		return err
	}

	for _, oldMessage := range oldMessages {
		var forSending message

		forSending.FromUser = oldMessage.FromUser
		forSending.Date = oldMessage.Date

		if oldMessage.ContentType == "image" {
			forSending.Img = oldMessage.Content
		} else if oldMessage.ContentType == "text" {
			forSending.Text = oldMessage.Content
		}

		c.incoming <- &forSending
	}

	return nil
}

func saveMessage(m *message) error {
	var (
		content     string
		contentType string
	)
	if m.Text != "" {
		content = m.Text
		contentType = "text"
	}
	if m.Img != "" {
		content = m.Img
		contentType = "image"
	}

	err := db.Insert(&db.Message{
		Content:     content,
		ContentType: contentType,
		FromUser:    m.FromUser,
		Date:        m.Date,
	})

	return err
}
