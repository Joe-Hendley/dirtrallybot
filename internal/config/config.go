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

// RandomiserType selects how a generated challenge's blanks are filled in.
type RandomiserType string

const (
	// RandomiserDeterministic replays a fixed sequence - reproducible across
	// restarts, mainly useful for debugging.
	RandomiserDeterministic RandomiserType = "deterministic"
	// RandomiserRandom picks uniformly at random.
	RandomiserRandom RandomiserType = "random"
	// RandomiserBiased leans towards items with good feedback.
	RandomiserBiased RandomiserType = "biased"
)

const defaultRandomiser = RandomiserBiased

// defaultWebAddr is where the local challenge viewer listens when WEBADDR is
// unset - loopback only. Set WEBADDR to ":8080" to expose it on the network, or
// to "off" to disable the viewer.
const defaultWebAddr = "localhost:8080"

type Config struct {
	App          string
	Token        string
	Store        StoreType
	Randomiser   RandomiserType
	TestServerID string
	// WebAddr is the listen address for the local challenge viewer, or "" when
	// it is disabled.
	WebAddr string
}

// Load reads configuration from a .env file in the working directory, with the
// process environment taking precedence. It is an error for the .env file to be
// missing.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("loading .env file: %w", err)
	}

	randomiser, err := randomiserFromEnv()
	if err != nil {
		return Config{}, err
	}

	return Config{
		App:          os.Getenv("APP"),
		Token:        os.Getenv("TOKEN"),
		Store:        defaultStore,
		Randomiser:   randomiser,
		TestServerID: os.Getenv("TESTSERVER"),
		WebAddr:      webAddrFromEnv(),
	}, nil
}

// webAddrFromEnv reads the WEBADDR key, defaulting when unset and treating "off"
// as a request to disable the viewer.
func webAddrFromEnv() string {
	switch addr := os.Getenv("WEBADDR"); addr {
	case "":
		return defaultWebAddr
	case "off":
		return ""
	default:
		return addr
	}
}

// randomiserFromEnv reads the RANDOMISER key, defaulting when unset and
// rejecting an unrecognised value rather than silently falling back.
func randomiserFromEnv() (RandomiserType, error) {
	switch value := RandomiserType(os.Getenv("RANDOMISER")); value {
	case "":
		return defaultRandomiser, nil
	case RandomiserDeterministic, RandomiserRandom, RandomiserBiased:
		return value, nil
	default:
		return "", fmt.Errorf("invalid randomiser %q", value)
	}
}
