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
			New("Ancelle", l, Unknown),
			New("Baisse de Patronel", l, Unknown),
			New("La Bâtie-Neuve - Saint-Léger-les-Mélèzes", l, Unknown),
			New("La Bollène-Vésubie - Col de Turini", l, Unknown),
			New("La Bollène-Vésubie - Peïra Cava", l, Unknown),
			New("La Maïris", l, Unknown),
			New("Les Borels", l, Unknown),
			New("Moissière", l, Unknown),
			New("Peïra Cava - La Bollène-Vésubie", l, Unknown),
			New("Pra d'Alart", l, Unknown),
			New("Ravin de Coste Belle", l, Unknown),
			New("Saint-Léger-les-Mélèzes - La Bâtie-Neuve", l, Unknown),
		}
	case location.SWE_WRC:
		return []Model{
			New("Älgsjön", l, Unknown),
			New("Åslia", l, Unknown),
			New("Åsnes", l, Unknown),
			New("Ekshärad", l, Unknown),
			New("Hof-Finnskog", l, Unknown),
			New("Knapptjernet", l, Unknown),
			New("Lauksjøen", l, Unknown),
			New("Lövstaholm", l, Unknown),
			New("Spikbrenna", l, Unknown),
			New("Stora Jangen", l, Unknown),
			New("Sunne", l, Unknown),
			New("Vargasen", l, Unknown),
		}
	case location.MEX:
		return []Model{
			New("Alfaro", l, Unknown),
			New("Derramadero", l, Unknown),
			New("El Brinco", l, Unknown),
			New("El Chocolate", l, Unknown),
			New("El Mosquito", l, Unknown),
			New("Guanajuatito", l, Unknown),
			New("Ibarrilla", l, Unknown),
			New("Las Minas", l, Unknown),
			New("Mesa Cuata", l, Unknown),
			New("Ortega", l, Unknown),
			New("Otates", l, Unknown),
			New("San Diego", l, Unknown),
		}
	case location.HRV:
		return []Model{
			New("Bliznec", l, Unknown),
			New("Grdanjci", l, Unknown),
			New("Hartje", l, Unknown),
			New("Kostanjevac", l, Unknown),
			New("Krašić", l, Unknown),
			New("Kumrovec", l, Unknown),
			New("Mali Lipovec", l, Unknown),
			New("Petruš Vrh", l, Unknown),
			New("Stojdraga", l, Unknown),
			New("Trakošćan", l, Unknown),
			New("Vrbno", l, Unknown),
			New("Zagorska Sela", l, Unknown),
		}
	case location.PRT:
		return []Model{
			New("Anjos", l, Unknown),
			New("Baião", l, Unknown),
			New("Caminha", l, Unknown),
			New("Carrazedo", l, Unknown),
			New("Celeiro", l, Unknown),
			New("Ervideiro", l, Unknown),
			New("Fridão", l, Unknown),
			New("Marão", l, Unknown),
			New("Ponte de Lima", l, Unknown),
			New("Touca", l, Unknown),
			New("Viana do Castelo", l, Unknown),
			New("Vila Boa", l, Unknown),
		}
	case location.ITA:
		return []Model{
			New("Alà del Sardi", l, Unknown),
			New("Bassacutena", l, Unknown),
			New("Bortigiadas", l, Unknown),
			New("Li Pinnenti", l, Unknown),
			New("Littichedda", l, Unknown),
			New("Malti", l, Unknown),
			New("Mamone", l, Unknown),
			New("Monte Acuto", l, Unknown),
			New("Monte Muvri", l, Unknown),
			New("Monte Olia", l, Unknown),
			New("Rena Majore", l, Unknown),
			New("Sa Mela", l, Unknown),
		}
	case location.KEN:
		return []Model{
			New("Kanyawa", l, Unknown),
			New("Kanyawa - Nakura", l, Unknown),
			New("Kingono", l, Unknown),
			New("Malewa", l, Unknown),
			New("Marula", l, Unknown),
			New("Mbaruk", l, Unknown),
			New("Moi North", l, Unknown),
			New("Nakuru", l, Unknown),
			New("Soysambu", l, Unknown),
			New("Sugunoi", l, Unknown),
			New("Tarambete", l, Unknown),
			New("Wileli", l, Unknown),
		}
	case location.EST:
		return []Model{
			New("Elva", l, Unknown),
			New("Koigu", l, Unknown),
			New("Kooraste", l, Unknown),
			New("Külaaseme", l, Unknown),
			New("Metsalaane", l, Unknown),
			New("Nüpli", l, Unknown),
			New("Otepää", l, Unknown),
			New("Rebaste", l, Unknown),
			New("Truuta", l, Unknown),
			New("Vahessaare", l, Unknown),
			New("Vellavere", l, Unknown),
			New("Vissi", l, Unknown),
		}
	case location.FIN_WRC:
		return []Model{
			New("Hatanpää", l, Unknown),
			New("Honkanen", l, Unknown),
			New("Lahdenkylä", l, Unknown),
			New("Leustu", l, Unknown),
			New("Maahi", l, Unknown),
			New("Päijälä", l, Unknown),
			New("Painaa", l, Unknown),
			New("Peltola", l, Unknown),
			New("Ruokolahti", l, Unknown),
			New("Saakoski", l, Unknown),
			New("Vehmas", l, Unknown),
			New("Venkajärvi", l, Unknown),
		}
	case location.GRC_WRC:
		return []Model{
			New("Amfissa", l, Unknown),
			New("Bauxites", l, Unknown),
			New("Delphi", l, Unknown),
			New("Drosochori", l, Unknown),
			New("Eptalofos", l, Unknown),
			New("Gravia", l, Unknown),
			New("Karoutes", l, Unknown),
			New("Lilea", l, Unknown),
			New("Mariolata", l, Unknown),
			New("Parnassós", l, Unknown),
			New("Prosilio", l, Unknown),
			New("Prosilio", l, Unknown),
		}
	case location.CHL:
		return []Model{
			New("Bio Bío", l, Unknown),
			New("Chivilingo", l, Unknown),
			New("El Poñen", l, Unknown),
			New("Hualqui", l, Unknown),
			New("Laja", l, Unknown),
			New("Las Patagues", l, Unknown),
			New("María Las Cruces", l, Unknown),
			New("Pulpería", l, Unknown),
			New("Rere", l, Unknown),
			New("Río Claro", l, Unknown),
			New("Río Lía", l, Unknown),
			New("Yumbel", l, Unknown),
		}
	case location.CER:
		return []Model{
			New("Brusné", l, Unknown),
			New("Chvalčov", l, Unknown),
			New("Libosváry", l, Unknown),
			New("Lukoveček", l, Unknown),
			New("Osíčko", l, Unknown),
			New("Příkazy", l, Unknown),
			New("Provodovice", l, Unknown),
			New("Raztoka", l, Unknown),
			New("Rouské", l, Unknown),
			New("Rusava", l, Unknown),
			New("Vítová", l, Unknown),
			New("Žabárna", l, Unknown),
		}
	case location.JPN:
		return []Model{
			New("Habucho", l, Unknown),
			New("Habu Dam", l, Unknown),
			New("Higashino", l, Unknown),
			New("Hokono Lake", l, Unknown),
			New("Kudarisawa", l, Unknown),
			New("Lake Mikawa", l, Unknown),
			New("Nakatsugawa", l, Unknown),
			New("Nenoue Highlands", l, Unknown),
			New("Nenoue Plateau", l, Unknown),
			New("Okuwacho", l, Unknown),
			New("Oninotaira", l, Unknown),
			New("Tegano", l, Unknown),
		}
	case location.MED:
		return []Model{
			New("Albarello", l, Unknown),
			New("Asco", l, Unknown),
			New("Cabanella", l, Unknown),
			New("Capannace", l, Unknown),
			New("Maririe", l, Unknown),
			New("Moltifao", l, Unknown),
			New("Monte Alloradu", l, Unknown),
			New("Monte Cinto", l, Unknown),
			New("Poggiola", l, Unknown),
			New("Ponte", l, Unknown),
			New("Ravin de Finelio", l, Unknown),
			New("Serra Di Cuzzioli", l, Unknown),
		}
	case location.PAC:
		return []Model{
			New("Abai", l, Unknown),
			New("Batukangkung", l, Unknown),
			New("Bidaralam", l, Unknown),
			New("Gunung Tujuh", l, Unknown),
			New("Kebun Raya Solok", l, Unknown),
			New("Loeboekmalaka", l, Unknown),
			New("Moearaikoer", l, Unknown),
			New("Sangir Balai Janggo", l, Unknown),
			New("South Solok", l, Unknown),
			New("Sungai Kunit", l, Unknown),
			New("Talanghilirair", l, Unknown),
			New("Talao", l, Unknown),
		}
	case location.UUU:
		return []Model{
			New("Brynderwyn", l, Unknown),
			New("Doctors Hill", l, Unknown),
			New("Makarau", l, Unknown),
			New("Mangapai", l, Unknown),
			New("Mareretu", l, Unknown),
			New("Noakes Hill", l, Unknown),
			New("Oakleigh", l, Unknown),
			New("Orewa", l, Unknown),
			New("Tahekeroa", l, Unknown),
			New("Tahekeroa - Orewa", l, Unknown),
			New("Taipuha", l, Unknown),
			New("Waiwera", l, Unknown),
		}
	case location.SCA:
		return []Model{
			New("Bergsøytjønn", l, Unknown),
			New("Dagtrolltjønn", l, Unknown),
			New("Fordol", l, Unknown),
			New("Fyresdal", l, Unknown),
			New("Fyresvatn", l, Unknown),
			New("Hengeltjønn", l, Unknown),
			New("Holtjønn", l, Unknown),
			New("Kottjønn", l, Unknown),
			New("Ljosdalstjønn", l, Unknown),
			New("Russvatn", l, Unknown),
			New("Tovsli", l, Unknown),
			New("Tovslioytjorn", l, Unknown),
		}
	case location.IBE:
		return []Model{
			New("Aiguamúrcia", l, Unknown),
			New("Alforja", l, Unknown),
			New("Botarell", l, Unknown),
			New("Campdasens", l, Unknown),
			New("L'Argentera", l, Unknown),
			New("Les Irles", l, Unknown),
			New("Les Voltes", l, Unknown),
			New("Montagut", l, Unknown),
			New("Montclar", l, Unknown),
			New("Pontils", l, Unknown),
			New("Santes Creus", l, Unknown),
			New("Valldossera", l, Unknown),
		}
	}
	return []Model{}
}
