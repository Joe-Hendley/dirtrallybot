package car_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/stretchr/testify/assert"
)

// Every class offered for a game must have at least one car, otherwise the
// builder can present a class the randomiser then cannot fill.
func TestEveryClassHasCars(t *testing.T) {
	for _, g := range []game.Model{game.DR2, game.WRC} {
		for _, c := range class.List(g) {
			cars := car.InClass(c, g)
			assert.NotEmptyf(t, cars, "%s / %s has no cars", g, c)

			for _, got := range cars {
				assert.Equalf(t, c, got.Class(), "car %q filed under the wrong class", got.Name())
			}
		}
	}
}

func TestInClassUnknownGameIsEmpty(t *testing.T) {
	assert.Empty(t, car.InClass(class.H1, game.NotSet))
}
