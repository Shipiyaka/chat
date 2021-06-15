package config

type Config struct {
	ServerAddr string `env:"CHAT_SERVER_ADDR" envDefault:"localhost:12345"`
	DBPath     string `env:"CHAT_DB_PATH" envDefault:"main.db"`
}
