package stage

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
)

type Distance int

const (
	Short      Distance = iota // 4 sectors
	Long                       // 8 sectors
	ReallyLong                 // 16 sectors?
	Unknown
)

func (d Distance) String() string {
	switch d {
	case Short:
		return "4 Sector"
	case Long:
		return "8 Sector"
	case ReallyLong:
		return "16 Sector"
	case Unknown:
		return "❓"
	}
	return "invalid distance"
}

func (d Distance) Emoji() string {
	switch d {
	case Short:
		return "4️⃣"
	case Long:
		return "8️⃣"
	case ReallyLong:
		return "♾️"
	case Unknown:
		return "❓"
	}
	return "invalid distance"
}

type Model struct {
	name     string
	location location.Model
	distance Distance
}

func New(name string, location location.Model, distance Distance) Model {
	return Model{
		name:     name,
		location: location,
		distance: distance,
	}
}

func (m Model) Name() string {
	return m.name
}

func (m Model) Distance() Distance {
	return m.distance
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

func AtLocationWithDistance(l location.Model, d Distance) []Model {
	atLocation := AtLocation(l)
	withDistance := []Model{}
	for _, stage := range atLocation {
		if stage.distance == d {
			withDistance = append(withDistance, stage)
		}
	}

	return withDistance
}

type stageData struct {
	name     string
	distance Distance
}

var byLocation = map[location.Model][]stageData{
	// DR2
	location.ARG: {
		{"Las Juntas", Long},
		{"Valle de los puentes", Long},
		{"Camino de acantilados y rocas", Short},
		{"San Isidro", Short},
		{"Miraflores", Short},
		{"El Rodeo", Short},
		{"Camino a La Puerta", Long},
		{"Valle de los puentes a la inversa", Long},
		{"Camino de acantilados y rocas inverso", Short},
		{"Camino a Coneta", Short},
		{"Huillaprima", Short},
		{"La Merced", Short},
	},
	location.AUS: {
		{"Mount Kaye Pass", Long},
		{"Chandlers Creek", Long},
		{"Bondi Forest", Short},
		{"Rockton Plains", Short},
		{"Yambulla Mountain Ascent", Short},
		{"Noorinbee Ridge Ascent", Short},
		{"Mount Kaye Pass Reverse", Long},
		{"Chandlers Creek Reverse", Long},
		{"Taylor Farm Sprint", Short},
		{"Rockton Plains Reverse", Short},
		{"Yambulla Mountain Descent", Short},
		{"Noorinbee Ridge Descent", Short},
	},
	location.FIN: {
		{"Kakaristo", Long},
		{"Kontinjärvi", Long},
		{"Kotajärvi", Short},
		{"Iso Oksjärvi", Short},
		{"Kailajärvi", Short},
		{"Naarajärvi", Short},
		{"Pitkäjärvi", Long},
		{"Hämelahti", Long},
		{"Oksala", Short},
		{"Järvenkylä", Short},
		{"Jyrkysjärvi", Short},
		{"Paskuri", Short},
	},
	location.DEU: {
		{"Oberstein", Long},
		{"Hammerstein", Long},
		{"Kreuzungsring", Short},
		{"Verbundsring", Short},
		{"Innerer Feld-Sprint", Short},
		{"Waldaufstieg", Short},
		{"Frauenberg", Long},
		{"Ruschberg", Long},
		{"Kreuzungsring reverse", Short},
		{"Verbundsring Reverse", Short},
		{"Innerer Feld-Sprint (umgekehrt)", Short},
		{"Waldabstieg", Short},
	},
	location.GRC: {
		{"Anodou Farmakas", Long},
		{"Pomona Érixi", Long},
		{"Koryfi Dafni", Short},
		{"Perasma Platani", Short},
		{"Ourea Spevsi", Short},
		{"Abies Koiláda", Short},
		{"Kathodo Leontiou", Long},
		{"Fourkéta Kourva", Long},
		{"Ampelonas Ormi", Short},
		{"Tsiristra Théa", Short},
		{"Pedines Epidaxi", Short},
		{"Ypsona tou Dasos", Short},
	},
	location.MCO: {
		{"Vallée descendante", Long},
		{"Pra d’Alart", Long},
		{"Col de Turini - Départ en descente", Short},
		{"Gordolon - Courte montée", Short},
		{"Col de Turini sprint en montée", Short},
		{"Route de Turini Descente", Short},
		{"Route de Turini", Long},
		{"Col de Turini Départ", Long},
		{"Route de Turini Montée", Short},
		{"Col de Turini - Descente", Short},
		{"Col de Turini - Sprint en descente", Short},
		{"Approche du Col de Turini - Montée", Short},
	},
	location.NZL: {
		{"Waimarama Point Forward", Long},
		{"Te Awanga Forward", Long},
		{"Waimarama Sprint Forward", Short},
		{"Elsthorpe Sprint Forward", Short},
		{"Ocean Beach Sprint Forward", Short},
		{"Te Awanga Sprint Forward", Short},
		{"Waimarama Point Reverse", Long},
		{"Ocean Beach", Long},
		{"Waimarama Sprint Reverse", Short},
		{"Elsthorpe Sprint Reverse", Short},
		{"Ocean Beach Sprint Reverse", Short},
		{"Te Awanga Sprint Reverse", Short},
	},
	location.POL: {
		{"Zaróbka", Long},
		{"Zienki", Long},
		{"Marynka", Short},
		{"Kopina", Short},
		{"Lejno", Short},
		{"Czarny Las", Short},
		{"Zagórze", Long},
		{"Jezioro Rotcze", Long},
		{"Borysik", Short},
		{"Józefin", Short},
		{"Jagodno", Short},
		{"Jezioro Lukie", Short},
	},
	location.SCO: {
		{"Newhouse Bridge", Long},
		{"South Morningside", Long},
		{"Annbank Station", Short},
		{"Rosebank Farm", Short},
		{"Old Butterstone Muir", Short},
		{"Glencastle Farm", Short},
		{"Newhouse Bridge Reverse", Long},
		{"South Morningside Reverse", Long},
		{"Annbank Station Reverse", Short},
		{"Rosebank Farm Reverse", Short},
		{"Old Butterstone Muir Reverse", Short},
		{"Glencastle Farm Reverse", Short},
	},
	location.ESP: {
		{"Comienzo De Bellriu", Long},
		{"Centenera", Long},
		{"Ascenso por valle el Gualet", Short},
		{"Viñedos dentro del valle Parra", Short},
		{"Viñedos Dardenyà", Short},
		{"Descenso por carretera", Short},
		{"Final de Bellriu", Long},
		{"Camino a Centenera", Long},
		{"Salida desde Montverd", Short},
		{"Ascenso bosque Montverd", Short},
		{"Viñedos Dardenyà inversa", Short},
		{"Subida por carretera", Short},
	},
	location.SWE: {
		{"Hamra", Long},
		{"Ransbysäter", Long},
		{"Elgsjön", Short},
		{"Stor-jangen Sprint", Short},
		{"Älgsjön Sprint", Short},
		{"Östra Hinnsjön", Short},
		{"Lysvik", Long},
		{"Norraskoga", Long},
		{"Älgsjön", Short},
		{"Stor-jangen Sprint Reverse", Short},
		{"Skogsrallyt", Short},
		{"Björklangen", Short},
	},
	location.USA: {
		{"Beaver Creek Trail Forward", Long},
		{"North Fork Pass", Long},
		{"Hancock Creek Burst", Short},
		{"Fuller Mountain Ascent", Short},
		{"Tolt Valley Sprint Forward", Short},
		{"Hancock Hill Sprint Forward", Short},
		{"Beaver Creek Trail Reverse", Long},
		{"North Fork Pass Reverse", Long},
		{"Fury Lake Depart", Short},
		{"Fuller Mountain Descent", Short},
		{"Tolt Valley Sprint Reverse", Short},
		{"Hancock Hill Sprint Reverse", Short},
	},
	location.WAL: {
		{"River Severn Valley", Long},
		{"Sweet Lamb", Long},
		{"Fferm Wynt", Short},
		{"Dyffryn Afon", Short},
		{"Bidno Moorland", Short},
		{"Pant Mawr", Short},
		{"Bronfelen", Long},
		{"Geufron Forest", Long},
		{"Fferm Wynt Reverse", Short},
		{"Dyffryn Afon Reverse", Short},
		{"Bidno Moorland Reverse", Short},
		{"Pant Mawr Reverse", Short},
	},
	// WRC
	location.MCO_WRC: {
		{"Ancelle", Short},
		{"Baisse de Patronel", Short},
		{"La Bâtie-Neuve - Saint-Léger-les-Mélèzes", Long},
		{"La Bollène-Vésubie - Col de Turini", Short},
		{"La Bollène-Vésubie - Peïra Cava", ReallyLong},
		{"La Maïris", Short},
		{"Les Borels", Short},
		{"Moissière", Short},
		{"Peïra Cava - La Bollène-Vésubie", ReallyLong},
		{"Pra d'Alart", Short},
		{"Ravin de Coste Belle", Short},
		{"Saint-Léger-les-Mélèzes - La Bâtie-Neuve", Long},
	},
	location.SWE_WRC: {
		{"Älgsjön", Short},
		{"Åslia", Short},
		{"Åsnes", ReallyLong},
		{"Ekshärad", Short},
		{"Hof-Finnskog", ReallyLong},
		{"Knapptjernet", Short},
		{"Lauksjøen", Long},
		{"Lövstaholm", Short},
		{"Spikbrenna", Long},
		{"Stora Jangen", Short},
		{"Sunne", Short},
		{"Vargasen", Short},
	},
	location.MEX: {
		{"Alfaro", Short},
		{"Derramadero", Short},
		{"El Brinco", Short},
		{"El Chocolate", ReallyLong},
		{"El Mosquito", Short},
		{"Guanajuatito", Short},
		{"Ibarrilla", Long},
		{"Las Minas", Long},
		{"Mesa Cuata", Short},
		{"Ortega", Short},
		{"Otates", ReallyLong},
		{"San Diego", Short},
	},
	location.HRV: {
		{"Bliznec", ReallyLong},
		{"Grdanjci", Long},
		{"Hartje", Short},
		{"Kostanjevac", Short},
		{"Krašić", Short},
		{"Kumrovec", Long},
		{"Mali Lipovec", Short},
		{"Petruš Vrh", Short},
		{"Stojdraga", Short},
		{"Trakošćan", ReallyLong},
		{"Vrbno", Short},
		{"Zagorska Sela", Short},
	},
	location.PRT: {
		{"Anjos", Short},
		{"Baião", ReallyLong},
		{"Caminha", ReallyLong},
		{"Carrazedo", Short},
		{"Celeiro", Short},
		{"Ervideiro", Short},
		{"Fridão", Long},
		{"Marão", Long},
		{"Ponte de Lima", Short},
		{"Touca", Short},
		{"Viana do Castelo", Short},
		{"Vila Boa", Short},
	},
	location.ITA: {
		{"Alà del Sardi", Long},
		{"Bassacutena", Short},
		{"Bortigiadas", Short},
		{"Li Pinnenti", Short},
		{"Littichedda", Short},
		{"Malti", Short},
		{"Mamone", Long},
		{"Monte Acuto", Short},
		{"Monte Muvri", Short},
		{"Monte Olia", ReallyLong},
		{"Rena Majore", ReallyLong},
		{"Sa Mela", Short},
	},
	location.KEN: {
		{"Kanyawa", Long},
		{"Kanyawa - Nakura", Long},
		{"Kingono", Short},
		{"Malewa", Short},
		{"Marula", Short},
		{"Mbaruk", ReallyLong},
		{"Moi North", Short},
		{"Nakuru", Short},
		{"Soysambu", ReallyLong},
		{"Sugunoi", Short},
		{"Tarambete", Short},
		{"Wileli", Short},
	},
	location.EST: {
		{"Elva", Long},
		{"Koigu", Short},
		{"Kooraste", Short},
		{"Külaaseme", Short},
		{"Metsalaane", Long},
		{"Nüpli", Short},
		{"Otepää", ReallyLong},
		{"Rebaste", ReallyLong},
		{"Truuta", Short},
		{"Vahessaare", Short},
		{"Vellavere", Short},
		{"Vissi", Short},
	},
	location.FIN_WRC: {
		{"Hatanpää", Long},
		{"Honkanen", Short},
		{"Lahdenkylä", Short},
		{"Leustu", Short},
		{"Maahi", Short},
		{"Päijälä", ReallyLong},
		{"Painaa", Short},
		{"Peltola", Short},
		{"Ruokolahti", ReallyLong},
		{"Saakoski", Short},
		{"Vehmas", Long},
		{"Venkajärvi", Short},
	},
	location.GRC_WRC: {
		{"Amfissa", Short},
		{"Bauxites", Short},
		{"Delphi", Short},
		{"Drosochori", Short},
		{"Eptalofos", Short},
		{"Gravia", ReallyLong},
		{"Karoutes", Long},
		{"Lilea", Short},
		{"Mariolata", Long},
		{"Parnassós", Short},
		{"Prosilio", ReallyLong},
		{"Viniani", Short},
	},
	location.CHL: {
		{"Arauco", ReallyLong},
		{"Bio Bío", ReallyLong},
		{"Coronel", Long},
		{"Lota", Long},
		{"Santa Juana", Long},
		{"Los Angeles", Short},
		{"San Rosendo", Short},
		{"Laja", Short},
		{"Yumbel", Short},
		{"Florida", Short},
		{"Hualqui", Short},
		{"Reputo", Short},
	},
	location.CER: {
		{"Brusné", Short},
		{"Chvalčov", Long},
		{"Libosváry", Short},
		{"Lukoveček", ReallyLong},
		{"Osíčko", Short},
		{"Příkazy", Short},
		{"Provodovice", Long},
		{"Raztoka", Short},
		{"Rouské", ReallyLong},
		{"Rusava", Short},
		{"Vítová", Short},
		{"Žabárna", Short},
	},
	location.JPN: {
		{"Habucho", Short},
		{"Habu Dam", Short},
		{"Higashino", Short},
		{"Hokono Lake", Short},
		{"Kudarisawa", ReallyLong},
		{"Lake Mikawa", ReallyLong},
		{"Nakatsugawa", Short},
		{"Nenoue Highlands", Short},
		{"Nenoue Plateau", Long},
		{"Okuwacho", Short},
		{"Oninotaira", Short},
		{"Tegano", Long},
	},
	location.MED: {
		{"Albarello", Long},
		{"Asco", ReallyLong},
		{"Cabanella", Short},
		{"Capannace", Long},
		{"Maririe", Short},
		{"Moltifao", Short},
		{"Monte Alloradu", Short},
		{"Monte Cinto", Short},
		{"Poggiola", Short},
		{"Ponte", ReallyLong},
		{"Ravin de Finelio", Short},
		{"Serra Di Cuzzioli", Short},
	},
	location.PAC: {
		{"Abai", Short},
		{"Batukangkung", Short},
		{"Bidaralam", Short},
		{"Gunung Tujuh", Short},
		{"Kebun Raya Solok", Short},
		{"Loeboekmalaka", Short},
		{"Moearaikoer", Short},
		{"Sangir Balai Janggo", Long},
		{"South Solok", Long},
		{"Sungai Kunit", Short},
		{"Talanghilirair", ReallyLong},
		{"Talao", ReallyLong},
	},
	location.UUU: {
		{"Brynderwyn", Short},
		{"Doctors Hill", Long},
		{"Makarau", Short},
		{"Mangapai", Short},
		{"Mareretu", Short},
		{"Noakes Hill", Short},
		{"Oakleigh", Long},
		{"Orewa", Short},
		{"Tahekeroa", ReallyLong},
		{"Tahekeroa - Orewa", Short},
		{"Taipuha", Short},
		{"Waiwera", ReallyLong},
	},
	location.SCA: {
		{"Bergsøytjønn", Short},
		{"Dagtrolltjønn", Short},
		{"Fordol", Short},
		{"Fyresdal", Short},
		{"Fyresvatn", Long},
		{"Hengeltjønn", ReallyLong},
		{"Holtjønn", ReallyLong},
		{"Kottjønn", Short},
		{"Ljosdalstjønn", Short},
		{"Russvatn", Long},
		{"Tovsli", Short},
		{"Tovslioytjorn", Short},
	},
	location.IBE: {
		{"Aiguamúrcia", Short},
		{"Alforja", Long},
		{"Botarell", Short},
		{"Campdasens", Short},
		{"L'Argentera", Short},
		{"Les Irles", Long},
		{"Les Voltes", Short},
		{"Montagut", Short},
		{"Montclar", Short},
		{"Pontils", Short},
		{"Santes Creus", ReallyLong},
		{"Valldossera", ReallyLong},
	},
}

func AtLocation(l location.Model) []Model {
	data := byLocation[l]
	stages := make([]Model, 0, len(data))
	for _, d := range data {
		stages = append(stages, New(d.name, l, d.distance))
	}
	return stages
}
