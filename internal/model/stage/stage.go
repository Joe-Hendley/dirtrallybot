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

func AtLocation(l location.Model) []Model {
	switch l {
	// DR2
	case location.ARG:
		return []Model{
			New("Las Juntas", l, Long),
			New("Valle de los puentes", l, Long),
			New("Camino de acantilados y rocas", l, Short),
			New("San Isidro", l, Short),
			New("Miraflores", l, Short),
			New("El Rodeo", l, Short),

			New("Camino a La Puerta", l, Long),
			New("Valle de los puentes a la inversa", l, Long),
			New("Camino de acantilados y rocas inverso", l, Short),
			New("Camino a Coneta", l, Short),
			New("Huillaprima", l, Short),
			New("La Merced", l, Short),
		}
	case location.AUS:
		return []Model{
			New("Mount Kaye Pass", l, Long),
			New("Chandlers Creek", l, Long),
			New("Bondi Forest", l, Short),
			New("Rockton Plains", l, Short),
			New("Yambulla Mountain Ascent", l, Short),
			New("Noorinbee Ridge Ascent", l, Short),

			New("Mount Kaye Pass Reverse", l, Long),
			New("Chandlers Creek Reverse", l, Long),
			New("Taylor Farm Sprint", l, Short),
			New("Rockton Plains Reverse", l, Short),
			New("Yambulla Mountain Descent", l, Short),
			New("Noorinbee Ridge Descent", l, Short),
		}
	case location.FIN:
		return []Model{New("Kakaristo", l, Long),
			New("Kontinjärvi", l, Long),
			New("Kotajärvi", l, Short),
			New("Iso Oksjärvi", l, Short),
			New("Kailajärvi", l, Short),
			New("Naarajärvi", l, Short),

			New("Pitkäjärvi", l, Long),
			New("Hämelahti", l, Long),
			New("Oksala", l, Short),
			New("Järvenkylä", l, Short),
			New("Jyrkysjärvi", l, Short),
			New("Paskuri", l, Short)}
	case location.DEU:
		return []Model{New("Oberstein", l, Long),
			New("Hammerstein", l, Long),
			New("Kreuzungsring", l, Short),
			New("Verbundsring", l, Short),
			New("Innerer Feld-Sprint", l, Short),
			New("Waldaufstieg", l, Short),

			New("Frauenberg", l, Long),
			New("Ruschberg", l, Long),
			New("Kreuzungsring reverse", l, Short),
			New("Verbundsring Reverse", l, Short),
			New("Innerer Feld-Sprint (umgekehrt)", l, Short),
			New("Waldabstieg", l, Short)}
	case location.GRC:
		return []Model{New("Anodou Farmakas", l, Long),
			New("Pomona Érixi", l, Long),
			New("Koryfi Dafni", l, Short),
			New("Perasma Platani", l, Short),
			New("Ourea Spevsi", l, Short),
			New("Abies Koiláda", l, Short),

			New("Kathodo Leontiou", l, Long),
			New("Fourkéta Kourva", l, Long),
			New("Ampelonas Ormi", l, Short),
			New("Tsiristra Théa", l, Short),
			New("Pedines Epidaxi", l, Short),
			New("Ypsona tou Dasos", l, Short)}
	case location.MCO:
		return []Model{New("Vallée descendante", l, Long),
			New("Pra d’Alart", l, Long),
			New("Col de Turini - Départ en descente", l, Short),
			New("Gordolon - Courte montée", l, Short),
			New("Col de Turini sprint en montée", l, Short),
			New("Route de Turini Descente", l, Short),

			New("Route de Turini", l, Long),
			New("Col de Turini Départ", l, Long),
			New("Route de Turini Montée", l, Short),
			New("Col de Turini - Descente", l, Short),
			New("Col de Turini - Sprint en descente", l, Short),
			New("Approche du Col de Turini - Montée", l, Short),
		}
	case location.NZL:
		return []Model{
			New("Waimarama Point Forward", l, Long),
			New("Te Awanga Forward", l, Long),
			New("Waimarama Sprint Forward", l, Short),
			New("Elsthorpe Sprint Forward", l, Short),
			New("Ocean Beach Sprint Forward", l, Short),
			New("Te Awanga Sprint Forward", l, Short),

			New("Waimarama Point Reverse", l, Long),
			New("Ocean Beach", l, Long),
			New("Waimarama Sprint Reverse", l, Short),
			New("Elsthorpe Sprint Reverse", l, Short),
			New("Ocean Beach Sprint Reverse", l, Short),
			New("Te Awanga Sprint Reverse", l, Short),
		}
	case location.POL:
		return []Model{
			New("Zaróbka", l, Long),
			New("Zienki", l, Long),
			New("Marynka", l, Short),
			New("Kopina", l, Short),
			New("Lejno", l, Short),
			New("Czarny Las", l, Short),

			New("Zagórze", l, Long),
			New("Jezioro Rotcze", l, Long),
			New("Borysik", l, Short),
			New("Józefin", l, Short),
			New("Jagodno", l, Short),
			New("Jezioro Lukie", l, Short),
		}
	case location.SCO:
		return []Model{New("Newhouse Bridge", l, Long),
			New("South Morningside", l, Long),
			New("Annbank Station", l, Short),
			New("Rosebank Farm", l, Short),
			New("Old Butterstone Muir", l, Short),
			New("Glencastle Farm", l, Short),

			New("Newhouse Bridge Reverse", l, Long),
			New("South Morningside Reverse", l, Long),
			New("Annbank Station Reverse", l, Short),
			New("Rosebank Farm Reverse", l, Short),
			New("Old Butterstone Muir Reverse", l, Short),
			New("Glencastle Farm Reverse", l, Short),
		}
	case location.ESP:
		return []Model{
			New("Comienzo De Bellriu", l, Long),
			New("Centenera", l, Long),
			New("Ascenso por valle el Gualet", l, Short),
			New("Viñedos dentro del valle Parra", l, Short),
			New("Viñedos Dardenyà", l, Short),
			New("Descenso por carretera", l, Short),

			New("Final de Bellriu", l, Long),
			New("Camino a Centenera", l, Long),
			New("Salida desde Montverd", l, Short),
			New("Ascenso bosque Montverd", l, Short),
			New("Viñedos Dardenyà inversa", l, Short),
			New("Subida por carretera", l, Short),
		}
	case location.SWE:
		return []Model{
			New("Hamra", l, Long),
			New("Ransbysäter", l, Long),
			New("Elgsjön", l, Short),
			New("Stor-jangen Sprint", l, Short),
			New("Älgsjön Sprint", l, Short),
			New("Östra Hinnsjön", l, Short),

			New("Lysvik", l, Long),
			New("Norraskoga", l, Long),
			New("Älgsjön", l, Short),
			New("Stor-jangen Sprint Reverse", l, Short),
			New("Skogsrallyt", l, Short),
			New("Björklangen", l, Short),
		}
	case location.USA:
		return []Model{
			New("Beaver Creek Trail Forward", l, Long),
			New("North Fork Pass", l, Long),
			New("Hancock Creek Burst", l, Short),
			New("Fuller Mountain Ascent", l, Short),
			New("Tolt Valley Sprint Forward", l, Short),
			New("Hancock Hill Sprint Forward", l, Short),

			New("Beaver Creek Trail Reverse", l, Long),
			New("North Fork Pass Reverse", l, Long),
			New("Fury Lake Depart", l, Short),
			New("Fuller Mountain Descent", l, Short),
			New("Tolt Valley Sprint Reverse", l, Short),
			New("Hancock Hill Sprint Reverse", l, Short),
		}
	case location.WAL:
		return []Model{
			New("River Severn Valley", l, Long),
			New("Sweet Lamb", l, Long),
			New("Fferm Wynt", l, Short),
			New("Dyffryn Afon", l, Short),
			New("Bidno Moorland", l, Short),
			New("Pant Mawr", l, Short),

			New("Bronfelen", l, Long),
			New("Geufron Forest", l, Long),
			New("Fferm Wynt Reverse", l, Short),
			New("Dyffryn Afon Reverse", l, Short),
			New("Bidno Moorland Reverse", l, Short),
			New("Pant Mawr Reverse", l, Short),
		}

		// WRC
	case location.MCO_WRC:
		return []Model{
			New("Ancelle", l, Short),
			New("Baisse de Patronel", l, Short),
			New("La Bâtie-Neuve - Saint-Léger-les-Mélèzes", l, Long),
			New("La Bollène-Vésubie - Col de Turini", l, Short),
			New("La Bollène-Vésubie - Peïra Cava", l, ReallyLong),
			New("La Maïris", l, Short),
			New("Les Borels", l, Short),
			New("Moissière", l, Short),
			New("Peïra Cava - La Bollène-Vésubie", l, ReallyLong),
			New("Pra d'Alart", l, Short),
			New("Ravin de Coste Belle", l, Short),
			New("Saint-Léger-les-Mélèzes - La Bâtie-Neuve", l, Long),
		}
	case location.SWE_WRC:
		return []Model{
			New("Älgsjön", l, Short),
			New("Åslia", l, Short),
			New("Åsnes", l, ReallyLong),
			New("Ekshärad", l, Short),
			New("Hof-Finnskog", l, ReallyLong),
			New("Knapptjernet", l, Short),
			New("Lauksjøen", l, Long),
			New("Lövstaholm", l, Short),
			New("Spikbrenna", l, Long),
			New("Stora Jangen", l, Short),
			New("Sunne", l, Short),
			New("Vargasen", l, Short),
		}
	case location.MEX:
		return []Model{
			New("Alfaro", l, Short),
			New("Derramadero", l, Short),
			New("El Brinco", l, Short),
			New("El Chocolate", l, ReallyLong),
			New("El Mosquito", l, Short),
			New("Guanajuatito", l, Short),
			New("Ibarrilla", l, Long),
			New("Las Minas", l, Long),
			New("Mesa Cuata", l, Short),
			New("Ortega", l, Short),
			New("Otates", l, ReallyLong),
			New("San Diego", l, Short),
		}
	case location.HRV:
		return []Model{
			New("Bliznec", l, ReallyLong),
			New("Grdanjci", l, Long),
			New("Hartje", l, Short),
			New("Kostanjevac", l, Short),
			New("Krašić", l, Short),
			New("Kumrovec", l, Long),
			New("Mali Lipovec", l, Short),
			New("Petruš Vrh", l, Short),
			New("Stojdraga", l, Short),
			New("Trakošćan", l, ReallyLong),
			New("Vrbno", l, Short),
			New("Zagorska Sela", l, Short),
		}
	case location.PRT:
		return []Model{
			New("Anjos", l, Short),
			New("Baião", l, ReallyLong),
			New("Caminha", l, ReallyLong),
			New("Carrazedo", l, Short),
			New("Celeiro", l, Short),
			New("Ervideiro", l, Short),
			New("Fridão", l, Long),
			New("Marão", l, Long),
			New("Ponte de Lima", l, Short),
			New("Touca", l, Short),
			New("Viana do Castelo", l, Short),
			New("Vila Boa", l, Short),
		}
	case location.ITA:
		return []Model{
			New("Alà del Sardi", l, Long),
			New("Bassacutena", l, Short),
			New("Bortigiadas", l, Short),
			New("Li Pinnenti", l, Short),
			New("Littichedda", l, Short),
			New("Malti", l, Short),
			New("Mamone", l, Long),
			New("Monte Acuto", l, Short),
			New("Monte Muvri", l, Short),
			New("Monte Olia", l, ReallyLong),
			New("Rena Majore", l, ReallyLong),
			New("Sa Mela", l, Short),
		}
	case location.KEN:
		return []Model{
			New("Kanyawa", l, Long),
			New("Kanyawa - Nakura", l, Long),
			New("Kingono", l, Short),
			New("Malewa", l, Short),
			New("Marula", l, Short),
			New("Mbaruk", l, ReallyLong),
			New("Moi North", l, Short),
			New("Nakuru", l, Short),
			New("Soysambu", l, ReallyLong),
			New("Sugunoi", l, Short),
			New("Tarambete", l, Short),
			New("Wileli", l, Short),
		}
	case location.EST:
		return []Model{
			New("Elva", l, Long),
			New("Koigu", l, Short),
			New("Kooraste", l, Short),
			New("Külaaseme", l, Short),
			New("Metsalaane", l, Long),
			New("Nüpli", l, Short),
			New("Otepää", l, ReallyLong),
			New("Rebaste", l, ReallyLong),
			New("Truuta", l, Short),
			New("Vahessaare", l, Short),
			New("Vellavere", l, Short),
			New("Vissi", l, Short),
		}
	case location.FIN_WRC:
		return []Model{
			New("Hatanpää", l, Long),
			New("Honkanen", l, Short),
			New("Lahdenkylä", l, Short),
			New("Leustu", l, Short),
			New("Maahi", l, Short),
			New("Päijälä", l, ReallyLong),
			New("Painaa", l, Short),
			New("Peltola", l, Short),
			New("Ruokolahti", l, ReallyLong),
			New("Saakoski", l, Short),
			New("Vehmas", l, Long),
			New("Venkajärvi", l, Short),
		}
	case location.GRC_WRC:
		return []Model{
			New("Amfissa", l, Short),
			New("Bauxites", l, Short),
			New("Delphi", l, Short),
			New("Drosochori", l, Short),
			New("Eptalofos", l, Short),
			New("Gravia", l, ReallyLong),
			New("Karoutes", l, Long),
			New("Lilea", l, Short),
			New("Mariolata", l, Long),
			New("Parnassós", l, Short),
			New("Prosilio", l, ReallyLong),
			New("Viniani", l, Short),
		}
	case location.CHL:
		return []Model{
			New("Arauco", l, ReallyLong),
			New("Bio Bío", l, ReallyLong),
			New("Coronel", l, Long),
			New("Lota", l, Long),
			New("Santa Juana", l, Long),
			New("Los Angeles", l, Short),
			New("San Rosendo", l, Short),
			New("Laja", l, Short),
			New("Yumbel", l, Short),
			New("Florida", l, Short),
			New("Hualqui", l, Short),
			New("Reputo", l, Short),
		}
	case location.CER:
		return []Model{
			New("Brusné", l, Short),
			New("Chvalčov", l, Long),
			New("Libosváry", l, Short),
			New("Lukoveček", l, ReallyLong),
			New("Osíčko", l, Short),
			New("Příkazy", l, Short),
			New("Provodovice", l, Long),
			New("Raztoka", l, Short),
			New("Rouské", l, ReallyLong),
			New("Rusava", l, Short),
			New("Vítová", l, Short),
			New("Žabárna", l, Short),
		}
	case location.JPN:
		return []Model{
			New("Habucho", l, Short),
			New("Habu Dam", l, Short),
			New("Higashino", l, Short),
			New("Hokono Lake", l, Short),
			New("Kudarisawa", l, ReallyLong),
			New("Lake Mikawa", l, ReallyLong),
			New("Nakatsugawa", l, Short),
			New("Nenoue Highlands", l, Short),
			New("Nenoue Plateau", l, Long),
			New("Okuwacho", l, Short),
			New("Oninotaira", l, Short),
			New("Tegano", l, Long),
		}
	case location.MED:
		return []Model{
			New("Albarello", l, Long),
			New("Asco", l, ReallyLong),
			New("Cabanella", l, Short),
			New("Capannace", l, Long),
			New("Maririe", l, Short),
			New("Moltifao", l, Short),
			New("Monte Alloradu", l, Short),
			New("Monte Cinto", l, Short),
			New("Poggiola", l, Short),
			New("Ponte", l, ReallyLong),
			New("Ravin de Finelio", l, Short),
			New("Serra Di Cuzzioli", l, Short),
		}
	case location.PAC:
		return []Model{
			New("Abai", l, Short),
			New("Batukangkung", l, Short),
			New("Bidaralam", l, Short),
			New("Gunung Tujuh", l, Short),
			New("Kebun Raya Solok", l, Short),
			New("Loeboekmalaka", l, Short),
			New("Moearaikoer", l, Short),
			New("Sangir Balai Janggo", l, Long),
			New("South Solok", l, Long),
			New("Sungai Kunit", l, Short),
			New("Talanghilirair", l, ReallyLong),
			New("Talao", l, ReallyLong),
		}
	case location.UUU:
		return []Model{
			New("Brynderwyn", l, Short),
			New("Doctors Hill", l, Long),
			New("Makarau", l, Short),
			New("Mangapai", l, Short),
			New("Mareretu", l, Short),
			New("Noakes Hill", l, Short),
			New("Oakleigh", l, Long),
			New("Orewa", l, Short),
			New("Tahekeroa", l, ReallyLong),
			New("Tahekeroa - Orewa", l, Short),
			New("Taipuha", l, Short),
			New("Waiwera", l, ReallyLong),
		}
	case location.SCA:
		return []Model{
			New("Bergsøytjønn", l, Short),
			New("Dagtrolltjønn", l, Short),
			New("Fordol", l, Short),
			New("Fyresdal", l, Short),
			New("Fyresvatn", l, Long),
			New("Hengeltjønn", l, ReallyLong),
			New("Holtjønn", l, ReallyLong),
			New("Kottjønn", l, Short),
			New("Ljosdalstjønn", l, Short),
			New("Russvatn", l, Long),
			New("Tovsli", l, Short),
			New("Tovslioytjorn", l, Short),
		}
	case location.IBE:
		return []Model{
			New("Aiguamúrcia", l, Short),
			New("Alforja", l, Long),
			New("Botarell", l, Short),
			New("Campdasens", l, Short),
			New("L'Argentera", l, Short),
			New("Les Irles", l, Long),
			New("Les Voltes", l, Short),
			New("Montagut", l, Short),
			New("Montclar", l, Short),
			New("Pontils", l, Short),
			New("Santes Creus", l, ReallyLong),
			New("Valldossera", l, ReallyLong),
		}
	}
	return []Model{}
}
