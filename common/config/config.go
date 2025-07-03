package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadDBConfig() DBConfig {
	_ = godotenv.Load()

	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		log.Fatal("POSTGRES_URL not set")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("Invalid POSTGRES_URL format: %v", err)
	}

	password, _ := u.User.Password()
	return DBConfig{
		Host:     u.Hostname(),
		Port:     u.Port(),
		User:     u.User.Username(),
		Password: password,
		DBName:   strings.TrimPrefix(u.Path, "/"),
	}
}