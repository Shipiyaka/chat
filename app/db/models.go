package db

type Message struct {
	ID            int `gorm:"primaryKey"`
	Content       string
	ContentType   string
	UsernameColor string
	FromUser      string
	Date          string
}
