package config

type Config struct {
	ServerAddr string `env:"CHAT_SERVER_ADDR" envDefault:"localhost:12345"`
	RedisAddr  string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPass  string `env:"REDIS_PASS" envDefault:""`
}
