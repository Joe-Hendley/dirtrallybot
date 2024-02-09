package challenge

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Randomiser interface {
	Car() car.Model
	Loc() location.Model
	Weather(loc location.Model) weather.Model
	Stage(loc location.Model) stage.Model
}

type Model struct {
	Stage   stage.Model
	Weather weather.Model
	Car     car.Model
}

func New(r Randomiser) Model {
	loc := r.Loc()

	return Model{
		Stage:   r.Stage(loc),
		Weather: r.Weather(loc),
		Car:     r.Car(),
	}
}

func (m Model) FancyString() string {

	return m.Stage.FancyString() + "\n" +
		location.WeatherStrings()[m.Stage.Location()][m.Weather] + "\n" +
		m.Car.FancyString() + "\n"
}
