package main

type User struct {
	Username string `json:"username"`
}

type Message struct {
	Text     string `json:"text,omitempty"`
	Img      string `json:"img,omitempty"`
	FromUser string `json:"from_user"`
	Date     string `json:"date"`
}
