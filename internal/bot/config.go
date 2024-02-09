package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	DEFAULTAPP   = "NOTSET"
	DEFAULTGUILD = "NOTSET"
	DEFAULTTOKEN = "NOTSET"
)

var (
	app   string = DEFAULTAPP
	token string = DEFAULTTOKEN
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
}

func NewConfig() Config {
	return Config{
		App:   app,
		Token: token,
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
