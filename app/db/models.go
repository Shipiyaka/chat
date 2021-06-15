package db

type Message struct {
	ID          int `gorm:"primaryKey"`
	Content     string
	ContentType string
	FromUser    string
	Date        string
}
