package class_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
)

func TestWithDrivetrain(t *testing.T) {
	for _, dt := range []drivetrain.Model{drivetrain.FWD, drivetrain.AWD, drivetrain.RWD} {
		for _, c := range class.WithDrivetrain(dt) {
			if c.Drivetrain() != dt {
				t.Errorf("got %s want %s for class %s", c.Drivetrain().String(), dt.String(), c.String())
			}
		}
	}
}
