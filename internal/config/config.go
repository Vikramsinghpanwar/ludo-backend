package config

import (
	"os"
)

type Config struct {
	SMS SMSConfig
}

func Load() *Config {
	return &Config{
		SMS: SMSConfig{
			APIKey:   os.Getenv("SMS_API_KEY"),
			UserID:   os.Getenv("SMS_USER_ID"),
			Password: os.Getenv("SMS_PASSWORD"),
			SenderID: os.Getenv("SMS_SENDER_ID"),
		},
	}
}
