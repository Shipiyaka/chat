package main

type Message struct {
	Text     string `json:"text"`
	FromUser string `json:"from_user"`
	Date     string `json:"date"`
}
