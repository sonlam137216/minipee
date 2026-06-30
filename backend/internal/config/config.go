package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort        string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiration  time.Duration
	FrontendOrigin string
}

func Load() (Config, error) {
	var missing []string
	getRequired := func(key string) string {
		value := os.Getenv(key)
		if value == "" {
			missing = append(missing, key)
		}
		return value
	}

	rawExpiration := getRequired("JWT_EXPIRATION_MINUTES")
	expirationMinutes, err := strconv.Atoi(rawExpiration)
	if rawExpiration != "" && (err != nil || expirationMinutes <= 0) {
		return Config{}, errors.New("JWT_EXPIRATION_MINUTES must be a positive integer")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return Config{
		AppPort:        getRequired("APP_PORT"),
		DatabaseURL:    getRequired("DATABASE_URL"),
		JWTSecret:      getRequired("JWT_SECRET"),
		JWTExpiration:  time.Duration(expirationMinutes) * time.Minute,
		FrontendOrigin: getRequired("FRONTEND_ORIGIN"),
	}, nil
}
