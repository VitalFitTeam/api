package config

import (
	"time"

	env "github.com/vitalfit/api/pkg/Env"
	"github.com/vitalfit/api/pkg/ratelimiter"
)

type Config struct {
	Addrs       string
	ApiUrl      string
	Db          dbConfig
	Env         string
	Mail        MailConfig
	Auth        AuthConfig
	RateLimiter ratelimiter.Config
}

type dbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type MailConfig struct {
	FromEmail string
	Exp       time.Duration
	Resend    ResendConfig
}

type ResendConfig struct {
	ApiKey string
}

type AuthConfig struct {
	Token TokenConfig
}
type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
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
		Mail: MailConfig{
			Exp:       time.Hour * 24 * 3, //3 days
			FromEmail: env.GetString("FROM_RESEND_EMAIL", ""),
			Resend: ResendConfig{
				ApiKey: env.GetString("RESEND_API_KEY", ""),
			},
		},
		Auth: AuthConfig{
			Token: TokenConfig{
				Secret: env.GetString("JWT_SECRET", ""),
				Exp:    time.Hour * 24 * 3, //3 days
				Iss:    env.GetString("JWT_ISS", ""),
			},
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_PER_TIME_FRAME", 150),
			TimeFrame:            time.Minute * 1,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}
}
