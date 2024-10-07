package car

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
)

const Emoji = "🏎️"

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

func (m Model) Name() string {
	return m.name
}

func (m Model) String() string {
	return m.name + " (" + m.class.String() + ")"
}

func (m Model) Class() class.Model {
	return m.class
}

func (m Model) FancyString() string {
	return fmt.Sprintf("%s **%s » %s**", Emoji, m.class.String(), m.name)
}

func WeightedMap() map[Model]int {
	cars := make(map[Model]int)
	for _, class := range class.List() {
		for _, car := range InClass(class) {
			cars[car] = 0
		}
	}

	return cars
}

func InClass(c class.Model) []Model {
	switch c {
	case class.H1:
		return []Model{
			New("Mini Cooper S", c),
			New("Lancia Fulva HF", c),
			New("DS Automobiles DS21", c)}
	case class.H2FWD:
		return []Model{
			New("Volkswagen Golf GTI 16v", c),
			New("Peugeot 205 GTI", c)}
	case class.H2RWD:
		return []Model{
			New("Ford Escort Mk II", c),
			New("Alpine Renault A110 1600 S", c),
			New("Fiat 131 Abarth", c),
			New("Opel Kadett C GT/E", c)}
	case class.H3:
		return []Model{
			New("BMW E30 M3 Evo Rally", c),
			New("Opel Ascona 400", c),
			New("Lancia Stratos", c),
			New("Renault 5 Turbo", c),
			New("Datsun 240Z", c),
			New("Ford Sierra Cosworth RS500", c)}
	case class.GroupBRWD:
		return []Model{
			New("Lancia 037 Evo 2", c),
			New("Opel Manta 400", c),
			New("BMW M1 Procar", c),
			New("Porsche 911 SC RS", c)}
	case class.GroupB4WD:
		return []Model{
			New("Audi Sport quattro S1 E2", c),
			New("Peugeot 205 T16 Evo 2", c),
			New("Lancia Delta S4", c),
			New("Ford RS200", c),
			New("MG Metro 6R4", c)}
	case class.R2:
		return []Model{
			New("Ford Fiesta R2", c),
			New("Opel Adam R2", c),
			New("Peugeot 208 R2", c)}
	case class.F2:
		return []Model{
			New("Peugeot 306 Maxi", c),
			New("SEAT Ibiza Kitcar", c),
			New("Volkswagen Golf Kitcar", c)}
	case class.GroupA:
		return []Model{
			New("Mitsubishi Lancer Evolution VI", c),
			New("Subaru Impreza 1995", c),
			New("Lancia Delta HF Integrale", c),
			New("Ford Escort RS Cosworth", c),
			New("Subaru Legacy RS", c)}
	case class.NR4:
		return []Model{
			New("Subaru WRX STI NR4", c),
			New("Mitsubishi Lancer Evolution X", c)}
	case class.WRC:
		return []Model{
			New("Ford Focus RS Rally 2001", c),
			New("Subaru Impreza (2001)", c),
			New("Citroën C4 Rally", c),
			New("Škoda Fabia Rally 2005", c),
			New("Ford Focus RS Rally 2007", c),
			New("Subaru Impreza", c),
			New("Peugeot 206 Rally", c),
			New("Subaru Impreza S4 Rally", c)}
	case class.R5:
		return []Model{
			New("Ford Fiesta R5", c),
			New("Ford Fiesta R5 MKII", c),
			New("Peugeot 208 R5 T16", c),
			New("Mitsubishi Space Star R5", c),
			New("ŠKODA Fabia R5", c),
			New("Citroën C3 R5", c),
			New("Volkswagen Polo GTI R5", c)}
	case class.RGT:
		return []Model{
			New("Chevrolet Camaro GT4-R", c),
			New("Porsche 911 RGT Rally Spec", c),
			New("Aston Martin V8 Vantage GT4", c),
			New("Ford Mustang GT4", c),
			New("BMW M2 Competition", c)}
	}
	return []Model{}
}
