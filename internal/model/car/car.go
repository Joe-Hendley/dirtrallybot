package car

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
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

func InClass(c class.Model, g game.Model) []Model {
	switch g {
	case game.DR2:
		return inClassDR2(c)
	case game.WRC:
		return inClassWRC(c)
	}
	return []Model{}
}

func inClassDR2(c class.Model) []Model {
	switch c {
	// DR2
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

// TODO - this
func inClassWRC(c class.Model) []Model {
	switch c {
	case class.WRC_WRC:
		return []Model{
			New("Ford Puma Rally1 HYBRID", c),
			New("Hyundai i20 N Rally1 HYBRID", c),
			New("Toyota GR Yaris Rally1 HYBRID", c),
		}
	case class.WRC2:
		return []Model{
			New("Citroën C3 Rally2", c),
			New("Ford Fiesta Rally2", c),
			New("Hyundai i20 N Rally2", c),
			New("Škoda Fabia Rally2 Evo", c),
			New("Škoda Fabia RS Rally2", c),
			New("Volkswagen Polo GTI R5", c),
		}
	case class.JuniorWRC:
		return []Model{
			New("Ford Fiesta Rally3", c),
		}
	case class.WRC2017to2021:
		return []Model{
			New("Ford Fiesta WRC", c),
			New("Volkswagen Polo 2017", c),
		}
	case class.WRC1997to2011:
		return []Model{
			New("Citroën C4 WRC", c),
			New("Citroën Xsara WRC", c),
			New("Ford Focus RS Rally 2001", c),
			New("Ford Focus RS Rally 2008", c),
			New("MINI Countryman Rally Edition", c),
			New("Mitsubishi Lancer Evolution VI", c),
			New("Peugeot 206 Rally", c),
			New("Seat C\u00f3rdoba WRC", c),
			New("Škoda Fabia WRC", c),
			New("SUBARU Impreza 1998", c),
			New("SUBARU Impreza 2001", c),
			New("SUBARU Impreza 2008", c),
		}
	case class.Rally2:
		return []Model{
			New("Ford Fiesta R5 MK7 Evo 2", c),
			New("Peugeot 208 T16 R5", c),
		}
	case class.Rally4:
		return []Model{
			New("Ford Fiesta MK8 Rally4", c),
			New("Opel Adam R2", c),
			New("Peugeot 208 Rally4", c),
			New("Renault Twingo II", c),
		}
	case class.NR4_WRC:
		return []Model{
			New("McRae R4", c),
			New("Mitsubishi Lancer Evolution X", c),
			New("SUBARU WRX STI NR4", c),
		}
	case class.S2000:
		return []Model{
			New("Fiat Grande Punto Abarth S2000", c),
			New("Opel Corsa S2000", c),
			New("Peugeot 207 S2000", c),
		}
	case class.S1600:
		return []Model{
			New("Citroën C2 Super 1600", c),
			New("Citroën Saxo Super 1600", c),
			New("Ford Puma S1600", c),
			New("Renault Clio S1600", c),
		}
	case class.F2_WRC:
		return []Model{
			New("Ford Escort Mk 6 Maxi", c),
			New("Peugeot 306 Maxi", c),
			New("Renault Maxi Mégane", c),
			New("Seat Ibiza Kit Car", c),
			New("Vauxhall Astra Rally Car", c),
			New("Volkswagen Golf IV Kit Car", c),
		}
	case class.GroupA_WRC:
		return []Model{
			New("Ford Escort RS Cosworth", c),
			New("Lancia Delta HF Integrale", c),
			New("Mitsubishi Galant VR4", c),
			New("SUBARU Impreza 1995", c),
			New("SUBARU Legacy RS", c),
		}
	case class.GroupB4WD_WRC:
		return []Model{
			New("Audi Sport quattro S1 (E2)", c),
			New("Ford RS200", c),
			New("Lancia Delta S4", c),
			New("MG Metro 6R4", c),
			New("Peugeot 205 T16 Evo 2", c),
		}
	case class.GroupBRWD_WRC:
		return []Model{
			New("BMW M1 Procar Rally", c),
			New("Lancia 037 Evo 2", c),
			New("Opel Manta 400", c),
			New("Porsche 911 SC RS", c),
		}
	case class.H3RWD_WRC:
		return []Model{
			New("BMW M3 Evo Rally", c),
			New("Ford Escort MK2 McRae Motorsport", c),
			New("Ford Sierra Cosworth RS500", c),
			New("Lancia Stratos", c),
			New("Opel Ascona 400", c),
			New("Renault 5 Turbo", c),
		}
	case class.H2RWD_WRC:
		return []Model{
			New("Alpine Renault A110 1600 S", c),
			New("Fiat 131 Abarth Rally", c),
			New("Ford Escort MK2", c),
			New("Hillman Avenger", c),
			New("Opel Kadett C GT/E", c),
			New("Talbot Sunbeam Lotus", c),
		}
	case class.H2FWD_WRC:
		return []Model{
			New("Peugeot 205 GTI", c),
			New("Peugeot 309 GTI", c),
			New("Volkswagen Golf GTI", c),
		}
	case class.H1_WRC:
		return []Model{
			New("Lancia Fulvia HF", c),
			New("MINI Cooper S", c),
			New("Vauxhall Nova Sport", c),
		}
	}
	return []Model{}
}
