package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type StoreType string

const (
	MEMORY StoreType = "memory"
)

const (
	DEFAULTAPP   = "NOTSET"
	DEFAULTGUILD = "NOTSET"
	DEFAULTTOKEN = "NOTSET"
	DEFAULTSTORE = MEMORY
)

var (
	app   string    = DEFAULTAPP
	token string    = DEFAULTTOKEN
	store StoreType = DEFAULTSTORE
)

func init() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	App   string
	Token string
	Store StoreType
}

func New() Config {
	return Config{
		App:   app,
		Token: token,
		Store: store,
	}
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	token = os.Getenv("token")
	app = os.Getenv("app")
	return nil
}
