package randomiser_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/stretchr/testify/assert"
)

// The fixed seed is deliberate: two randomisers for the same game must produce
// the same sequence of stages.
func TestDeterministicIsReproducible(t *testing.T) {
	a := randomiser.NewDeterministic(game.DR2)
	b := randomiser.NewDeterministic(game.DR2)

	for range 20 {
		loc := a.Loc()
		assert.Equal(t, a.Stage(loc), b.Stage(b.Loc()))
	}
}
