package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var ErrEnv = errors.New("environment variables failed to load")

type Config struct {
	MongodbURI  string
	MongodbName string
	JwtSecret   string
	ServerPort  string
	BackendURL string
	MailtrapHost string
	MailtrapPort string
	MailtrapUsername string
	MailtrapPassword string
	MailtrapFrom string 
}

func LoadEnv() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		MongodbURI:  os.Getenv("MONGODB_URI"),
		MongodbName: os.Getenv("MONGODB_NAME"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		BackendURL: os.Getenv("BACKEND_BASE_URL"),
		MailtrapHost: os.Getenv("MAILTRAP_HOST"),
		MailtrapPort: os.Getenv("MAILTRAP_PORT"),
		MailtrapUsername: os.Getenv("MAILTRAP_USERNAME"),
		MailtrapPassword: os.Getenv("MAILTRAP_PASSWORD"),
		MailtrapFrom: os.Getenv("MAILTRAP_FROM"),
	}

	var missing []string
	if cfg.MongodbURI == "" {
		missing = append(missing, "MONGODB_URI")
	}
	if cfg.MongodbName == "" {
		missing = append(missing, "MONGODB_NAME")
	}
	if cfg.JwtSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if cfg.ServerPort == "" {
		missing = append(missing, "SERVER_PORT")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %v", strings.Join(missing, ", "))
	}

	return cfg, nil
}