package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ouz/gobackend/errors"
)

type Config struct {
	DB_HOST                string
	DB_USER                string
	DB_PASSWORD            string
	DB_NAME                string
	DB_PORT                string
	JWT_SECRET             string
	JWT_ACCESS_EXPIRATION  time.Duration
	JWT_REFRESH_EXPIRATION time.Duration
}

var Conf Config

func LoadConfig() error {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	accessExpiration, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRATION"))
	if err != nil {
		return errors.ValidationError("Invalid JWT_ACCESS_EXPIRATION", err)
	}

	refreshExpiration, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRATION"))
	if err != nil {
		return errors.ValidationError("Invalid JWT_REFRESH_EXPIRATION", err)
	}

	Conf = Config{
		DB_HOST:                os.Getenv("DB_HOST"),
		DB_USER:                os.Getenv("DB_USER"),
		DB_PASSWORD:            os.Getenv("DB_PASSWORD"),
		DB_NAME:                os.Getenv("DB_NAME"),
		DB_PORT:                os.Getenv("DB_PORT"),
		JWT_SECRET:             os.Getenv("JWT_SECRET"),
		JWT_ACCESS_EXPIRATION:  accessExpiration,
		JWT_REFRESH_EXPIRATION: refreshExpiration,
	}

	if err := validateConfig(); err != nil {
		return err
	}

	return nil
}

func validateConfig() error {
	if Conf.DB_HOST == "" {
		return errors.ValidationError("DB_HOST is not set", nil)
	}
	if Conf.DB_USER == "" {
		return errors.ValidationError("DB_USER is not set", nil)
	}
	if Conf.DB_PASSWORD == "" {
		return errors.ValidationError("DB_PASSWORD is not set", nil)
	}
	if Conf.DB_NAME == "" {
		return errors.ValidationError("DB_NAME is not set", nil)
	}
	if Conf.DB_PORT == "" {
		return errors.ValidationError("DB_PORT is not set", nil)
	}
	if Conf.JWT_SECRET == "" {
		return errors.ValidationError("JWT_SECRET is not set", nil)
	}
	return nil
}

func GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		Conf.DB_HOST, Conf.DB_USER, Conf.DB_PASSWORD, Conf.DB_NAME, Conf.DB_PORT)
}
