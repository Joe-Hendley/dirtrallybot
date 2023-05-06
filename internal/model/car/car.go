package car

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
)

type Model struct {
	name  string
	class class.Model
}

func New(name string, class class.Model) Model {
	return Model{
		name:  name,
		class: class,
	}
}

func (m Model) String() string {
	return m.name + " (" + m.class.String() + ")"
}

func (m Model) Class() class.Model {
	return m.class
}
