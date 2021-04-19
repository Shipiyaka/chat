package main

import (
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

func unmarshalMessage(b []byte) (Message, error) {
	var message Message
	
	err := json.Unmarshal(b, &message)
	if err != nil {
		return message, err
	}
	
	return message, nil
}
