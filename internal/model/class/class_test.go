package class_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
)

func TestWithDrivetrain(t *testing.T) {
	t.Run("Dirt Rally 2", func(t *testing.T) {
		for _, dt := range []drivetrain.Model{drivetrain.FWD, drivetrain.AWD, drivetrain.RWD} {
			for _, c := range class.WithDrivetrain(dt, game.DR2) {
				if c.Drivetrain() != dt {
					t.Errorf("got %s want %s for class %s", c.Drivetrain().String(), dt.String(), c.String())
				}
			}
		}
	})
	t.Run("WRC", func(t *testing.T) {
		for _, dt := range []drivetrain.Model{drivetrain.FWD, drivetrain.AWD, drivetrain.RWD} {
			for _, c := range class.WithDrivetrain(dt, game.WRC) {
				if c.Drivetrain() != dt {
					t.Errorf("got %s want %s for class %s", c.Drivetrain().String(), dt.String(), c.String())
				}
			}
		}
	})
}
