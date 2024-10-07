package store

import (
	"errors"

	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/boltstore"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
)

const defaultBoltPath string = "rallybot.db"

func New(cfg config.Config) (model.Store, error) {
	switch cfg.Store {
	case config.MEMORY:
		return memorystore.New(), nil
	case config.BOLT:
		return boltstore.New(defaultBoltPath)
	}
	return nil, errors.New("invalid store type " + string(cfg.Store))
}
