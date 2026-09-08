package challenge

import (
	"context"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorsBuildTheirRandomiser(t *testing.T) {
	store := memorystore.New()

	assert.IsType(t, &randomiser.Deterministic{}, DeterministicGenerator(context.Background(), store, game.DR2))
	assert.IsType(t, &randomiser.Random{}, RandomGenerator(context.Background(), store, game.DR2))
	assert.IsType(t, &randomiser.Biased{}, BiasedGenerator(context.Background(), store, game.DR2))
}

func TestBiasedGeneratorToleratesAStoreError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // memorystore refuses a cancelled context

	r := BiasedGenerator(ctx, memorystore.New(), game.DR2)

	// Falls back to an empty snapshot rather than failing; still usable.
	assert.Contains(t, location.List(game.DR2), r.Loc())
}

func TestSetGenerator(t *testing.T) {
	t.Cleanup(func() { newRandomiser = BiasedGenerator })

	SetGenerator(RandomGenerator)

	assert.IsType(t, &randomiser.Random{}, newRandomiser(context.Background(), memorystore.New(), game.WRC))
}
