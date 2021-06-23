package hub

import (
	"chat/app/db"
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	colorCharset    = "0123456789ABCDEF"
	randomStringLen = 5
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomColor() string {
	var strBuilder strings.Builder
	strBuilder.Grow(6)

	strBuilder.WriteString("#")

	for i := 0; i < 6; i++ {
		strBuilder.WriteByte(colorCharset[seededRand.Intn(16)])
	}

	return strBuilder.String()
}

func randomString() string {
	var strBuilder strings.Builder
	strBuilder.Grow(randomStringLen)

	for i := 0; i < randomStringLen; i++ {
		strBuilder.WriteByte(charset[seededRand.Intn(len(charset))])
	}

	return strBuilder.String()
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
		forSending := message{
			FromUser:      oldMessage.FromUser,
			Date:          oldMessage.Date,
			UsernameColor: oldMessage.UsernameColor,
		}

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
		Content:       content,
		ContentType:   contentType,
		UsernameColor: m.UsernameColor,
		FromUser:      m.FromUser,
		Date:          m.Date,
	})

	return err
}
