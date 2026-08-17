package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (Config, error) {
	var cfg Config
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Err in openning .env file:%s\n", err.Error())
		return cfg, err
	}
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.SSLMode = os.Getenv("SSL_MODE")
	cfg.Sign = os.Getenv("SIGN")

	return cfg, nil

}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	Sign       string
}
