package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type StoreType string

const (
	MEMORY StoreType = "memory"
	BOLT   StoreType = "bolt"
)

const (
	DEFAULTAPP        = "NOTSET"
	DEFAULTTESTSERVER = "NOTSET"
	DEFAULTTOKEN      = "NOTSET"
	DEFAULTSTORE      = BOLT
)

var (
	app        string    = DEFAULTAPP
	token      string    = DEFAULTTOKEN
	store      StoreType = DEFAULTSTORE
	testServer string    = DEFAULTTESTSERVER
)

func init() {
	err := loadEnv()
	if err != nil {
		slog.Error("loading env file", "err", err)
		os.Exit(1)
	}
}

type Config struct {
	App          string
	Token        string
	Store        StoreType
	TestServerID string
}

func New() Config {
	return Config{
		App:          app,
		Token:        token,
		Store:        store,
		TestServerID: testServer,
	}
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	token = os.Getenv("token")
	app = os.Getenv("app")
	testServer = os.Getenv("testserver")
	return nil
}
