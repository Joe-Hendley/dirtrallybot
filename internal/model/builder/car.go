package builder

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
)

func Cars() map[car.Model]int {
	cars := make(map[car.Model]int)
	for _, class := range class.List() {
		for _, car := range CarsInClass(class) {
			cars[car] = 0
		}
	}

	return cars
}

func CarsInClass(c class.Model) []car.Model {
	switch c {
	case class.H1:
		return []car.Model{
			car.New("Mini Cooper S", c),
			car.New("Lancia Fulva HF", c),
			car.New("DS Automobiles DS21", c)}
	case class.H2FWD:
		return []car.Model{
			car.New("Volkswagen Golf GTI 16v", c),
			car.New("Peugeot 205 GTI", c)}
	case class.H2RWD:
		return []car.Model{
			car.New("Ford Escort Mk II", c),
			car.New("Alpine Renault A110 1600 S", c),
			car.New("Fiat 131 Abarth", c),
			car.New("Opel Kadett C GT/E", c)}
	case class.H3:
		return []car.Model{
			car.New("BMW E30 M3 Evo Rally", c),
			car.New("Opel Ascona 400", c),
			car.New("Lancia Stratos", c),
			car.New("Renault 5 Turbo", c),
			car.New("Datsun 240Z", c),
			car.New("Ford Sierra Cosworth RS500", c)}
	case class.GroupBRWD:
		return []car.Model{
			car.New("Lancia 037 Evo 2", c),
			car.New("Opel Manta 400", c),
			car.New("BMW M1 Procar", c),
			car.New("Porsche 911 SC RS", c)}
	case class.GroupB4WD:
		return []car.Model{
			car.New("Audi Sport quattro S1 E2", c),
			car.New("Peugeot 205 T16 Evo 2", c),
			car.New("Lancia Delta S4", c),
			car.New("Ford RS200", c),
			car.New("MG Metro 6R4", c)}
	case class.R2:
		return []car.Model{
			car.New("Ford Fiesta R2", c),
			car.New("Opel Adam R2", c),
			car.New("Peugeot 208 R2", c)}
	case class.F2:
		return []car.Model{
			car.New("Peugeot 306 Maxi", c),
			car.New("SEAT Ibiza Kitcar", c),
			car.New("Volkswagen Golf Kitcar", c)}
	case class.GroupA:
		return []car.Model{
			car.New("Mitsubishi Lancer Evolution VI", c),
			car.New("Subaru Impreza 1995", c),
			car.New("Lancia Delta HF Integrale", c),
			car.New("Ford Escort RS Cosworth", c),
			car.New("Subaru Legacy RS", c)}
	case class.NR4:
		return []car.Model{
			car.New("Subaru WRX STI NR4", c),
			car.New("Mitsubishi Lancer Evolution X", c)}
	case class.WRC:
		return []car.Model{
			car.New("Ford Focus RS Rally 2001", c),
			car.New("Subaru Impreza (2001)", c),
			car.New("Citroën C4 Rally", c),
			car.New("Škoda Fabia Rally 2005", c),
			car.New("Ford Focus RS Rally 2007", c),
			car.New("Subaru Impreza", c),
			car.New("Peugeot 206 Rally", c),
			car.New("Subaru Impreza S4 Rally", c)}
	case class.R5:
		return []car.Model{
			car.New("Ford Fiesta R5", c),
			car.New("Ford Fiesta R5 MKII", c),
			car.New("Peugeot 208 R5 T16", c),
			car.New("Mitsubishi Space Star R5", c),
			car.New("ŠKODA Fabia R5", c),
			car.New("Citroën C3 R5", c),
			car.New("Volkswagen Polo GTI R5", c)}
	case class.RGT:
		return []car.Model{
			car.New("Chevrolet Camaro GT4-R", c),
			car.New("Porsche 911 RGT Rally Spec", c),
			car.New("Aston Martin V8 Vantage GT4", c),
			car.New("Ford Mustang GT4", c),
			car.New("BMW M2 Competition", c)}
	}
	return []car.Model{}
}
