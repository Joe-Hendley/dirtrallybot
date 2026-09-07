package location_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	dr2 := location.List(game.DR2)
	wrc := location.List(game.WRC)

	assert.NotEmpty(t, dr2)
	assert.NotEmpty(t, wrc)

	for _, g := range []struct {
		name string
		list []location.Model
	}{{"DR2", dr2}, {"WRC", wrc}} {
		seen := map[location.Model]bool{}
		for _, loc := range g.list {
			assert.Falsef(t, seen[loc], "%s lists %s twice", g.name, loc)
			seen[loc] = true
			assert.NotEqualf(t, "invalid location", loc.String(), "%s lists an unnamed location", g.name)
		}
	}
}

func TestListUnknownGameIsEmpty(t *testing.T) {
	assert.Empty(t, location.List(game.NotSet))
}
