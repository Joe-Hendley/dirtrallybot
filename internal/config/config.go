package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type StoreType string

const (
	MEMORY StoreType = "memory"
	BOLT   StoreType = "bolt"
)

const defaultStore = BOLT

type Config struct {
	App          string
	Token        string
	Store        StoreType
	TestServerID string
}

// Load reads configuration from a .env file in the working directory, with the
// process environment taking precedence. It is an error for the .env file to be
// missing.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("loading .env file: %w", err)
	}

	return Config{
		App:          os.Getenv("app"),
		Token:        os.Getenv("token"),
		Store:        defaultStore,
		TestServerID: os.Getenv("testserver"),
	}, nil
}
