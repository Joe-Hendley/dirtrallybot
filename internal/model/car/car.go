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

var namesByGameClass = map[game.Model]map[class.Model][]string{
	game.DR2: {
		class.H1: {
			"Mini Cooper S",
			"Lancia Fulva HF",
			"DS Automobiles DS21",
		},
		class.H2FWD: {
			"Volkswagen Golf GTI 16v",
			"Peugeot 205 GTI",
		},
		class.H2RWD: {
			"Ford Escort Mk II",
			"Alpine Renault A110 1600 S",
			"Fiat 131 Abarth",
			"Opel Kadett C GT/E",
		},
		class.H3: {
			"BMW E30 M3 Evo Rally",
			"Opel Ascona 400",
			"Lancia Stratos",
			"Renault 5 Turbo",
			"Datsun 240Z",
			"Ford Sierra Cosworth RS500",
		},
		class.GroupBRWD: {
			"Lancia 037 Evo 2",
			"Opel Manta 400",
			"BMW M1 Procar",
			"Porsche 911 SC RS",
		},
		class.GroupB4WD: {
			"Audi Sport quattro S1 E2",
			"Peugeot 205 T16 Evo 2",
			"Lancia Delta S4",
			"Ford RS200",
			"MG Metro 6R4",
		},
		class.R2: {
			"Ford Fiesta R2",
			"Opel Adam R2",
			"Peugeot 208 R2",
		},
		class.F2: {
			"Peugeot 306 Maxi",
			"SEAT Ibiza Kitcar",
			"Volkswagen Golf Kitcar",
		},
		class.GroupA: {
			"Mitsubishi Lancer Evolution VI",
			"Subaru Impreza 1995",
			"Lancia Delta HF Integrale",
			"Ford Escort RS Cosworth",
			"Subaru Legacy RS",
		},
		class.NR4: {
			"Subaru WRX STI NR4",
			"Mitsubishi Lancer Evolution X",
		},
		class.WRC: {
			"Ford Focus RS Rally 2001",
			"Subaru Impreza (2001)",
			"Citroën C4 Rally",
			"Škoda Fabia Rally 2005",
			"Ford Focus RS Rally 2007",
			"Subaru Impreza",
			"Peugeot 206 Rally",
			"Subaru Impreza S4 Rally",
		},
		class.R5: {
			"Ford Fiesta R5",
			"Ford Fiesta R5 MKII",
			"Peugeot 208 R5 T16",
			"Mitsubishi Space Star R5",
			"ŠKODA Fabia R5",
			"Citroën C3 R5",
			"Volkswagen Polo GTI R5",
		},
		class.RGT: {
			"Chevrolet Camaro GT4-R",
			"Porsche 911 RGT Rally Spec",
			"Aston Martin V8 Vantage GT4",
			"Ford Mustang GT4",
			"BMW M2 Competition",
		},
	},
	game.WRC: {
		class.WRC_WRC: {
			"Ford Puma Rally1 HYBRID",
			"Hyundai i20 N Rally1 HYBRID",
			"Toyota GR Yaris Rally1 HYBRID",
		},
		class.WRC2: {
			"Citroën C3 Rally2",
			"Ford Fiesta Rally2",
			"Hyundai i20 N Rally2",
			"Škoda Fabia Rally2 Evo",
			"Škoda Fabia RS Rally2",
			"Volkswagen Polo GTI R5",
		},
		class.JuniorWRC: {
			"Ford Fiesta Rally3",
		},
		class.WRC2017to2021: {
			"Ford Fiesta WRC",
			"Volkswagen Polo 2017",
		},
		class.WRC1997to2011: {
			"Citroën C4 WRC",
			"Citroën Xsara WRC",
			"Ford Focus RS Rally 2001",
			"Ford Focus RS Rally 2008",
			"MINI Countryman Rally Edition",
			"Mitsubishi Lancer Evolution VI",
			"Peugeot 206 Rally",
			"Seat C\u00f3rdoba WRC",
			"Škoda Fabia WRC",
			"SUBARU Impreza 1998",
			"SUBARU Impreza 2001",
			"SUBARU Impreza 2008",
		},
		class.Rally2: {
			"Ford Fiesta R5 MK7 Evo 2",
			"Peugeot 208 T16 R5",
		},
		class.Rally4: {
			"Ford Fiesta MK8 Rally4",
			"Opel Adam R2",
			"Peugeot 208 Rally4",
			"Renault Twingo II",
		},
		class.NR4_WRC: {
			"McRae R4",
			"Mitsubishi Lancer Evolution X",
			"SUBARU WRX STI NR4",
		},
		class.S2000: {
			"Fiat Grande Punto Abarth S2000",
			"Opel Corsa S2000",
			"Peugeot 207 S2000",
		},
		class.S1600: {
			"Citroën C2 Super 1600",
			"Citroën Saxo Super 1600",
			"Ford Puma S1600",
			"Renault Clio S1600",
		},
		class.F2_WRC: {
			"Ford Escort Mk 6 Maxi",
			"Peugeot 306 Maxi",
			"Renault Maxi Mégane",
			"Seat Ibiza Kit Car",
			"Vauxhall Astra Rally Car",
			"Volkswagen Golf IV Kit Car",
		},
		class.GroupA_WRC: {
			"Ford Escort RS Cosworth",
			"Lancia Delta HF Integrale",
			"Mitsubishi Galant VR4",
			"SUBARU Impreza 1995",
			"SUBARU Legacy RS",
		},
		class.GroupB4WD_WRC: {
			"Audi Sport quattro S1 (E2)",
			"Ford RS200",
			"Lancia Delta S4",
			"MG Metro 6R4",
			"Peugeot 205 T16 Evo 2",
		},
		class.GroupBRWD_WRC: {
			"BMW M1 Procar Rally",
			"Lancia 037 Evo 2",
			"Opel Manta 400",
			"Porsche 911 SC RS",
		},
		class.H3RWD_WRC: {
			"BMW M3 Evo Rally",
			"Ford Escort MK2 McRae Motorsport",
			"Ford Sierra Cosworth RS500",
			"Lancia Stratos",
			"Opel Ascona 400",
			"Renault 5 Turbo",
		},
		class.H2RWD_WRC: {
			"Alpine Renault A110 1600 S",
			"Fiat 131 Abarth Rally",
			"Ford Escort MK2",
			"Hillman Avenger",
			"Opel Kadett C GT/E",
			"Talbot Sunbeam Lotus",
		},
		class.H2FWD_WRC: {
			"Peugeot 205 GTI",
			"Peugeot 309 GTI",
			"Volkswagen Golf GTI",
		},
		class.H1_WRC: {
			"Lancia Fulvia HF",
			"MINI Cooper S",
			"Vauxhall Nova Sport",
		},
	},
}

func InClass(c class.Model, g game.Model) []Model {
	names := namesByGameClass[g][c]
	models := make([]Model, 0, len(names))
	for _, name := range names {
		models = append(models, New(name, c))
	}
	return models
}
