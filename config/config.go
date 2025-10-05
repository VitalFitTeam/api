package config

import env "github.com/vitalfit/api/pkg/Env"

type Config struct {
	Addrs  string
	ApiUrl string
	Db     dbConfig
	Env    string
}

type dbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func LoadConfig() *Config {
	return &Config{
		Addrs: env.GetString("ADDRS", ":8080"),
		Db: dbConfig{
			Dsn:          env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/vitalfit?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env:    env.GetString("ENV", "dev"),
		ApiUrl: env.GetString("API_URL", "localhost:8080"),
	}
}
