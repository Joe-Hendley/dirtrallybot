package stage

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
)

type distance bool

const (
	Short distance = true
	Long  distance = false
)

type Model struct {
	name     string
	location location.Model
	distance distance
}

func New(name string, location location.Model, distance distance) Model {
	return Model{
		name:     name,
		location: location,
		distance: distance,
	}
}

func (m Model) String() string {
	return m.location.String() + ": " + m.name
}

func (m Model) FancyString() string {
	return fmt.Sprintf("%s **%s » %s**", m.location.Flag(), m.location.String(), m.name)
}

func (m Model) Location() location.Model {
	return m.location
}
