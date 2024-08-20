package store

import (
	"errors"

	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
)

func New(cfg config.Config) (model.Store, error) {
	switch cfg.Store {
	case config.MEMORY:
		return memorystore.New(), nil
	}
	return nil, errors.New("invalid store type " + string(cfg.Store))
}
