package stage

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type distance bool

const (
	Short distance = true
	Long  distance = false
)

type Model struct {
	name     string
	location location.Model
	weather  weather.Model
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

func WeightedMap() map[Model]int {
	stages := make(map[Model]int)

	for _, l := range location.List() {
		for _, s := range AtLocation(l) {
			stages[s] = 0
		}
	}

	return stages
}

func AtLocation(l location.Model) []Model {
	switch l {
	case location.ARG:
		return []Model{
			New("Las Juntas", location.ARG, Long),
			New("Valle de los puentes", location.ARG, Long),
			New("Camino de acantilados y rocas", location.ARG, Short),
			New("San Isidro", location.ARG, Short),
			New("Miraflores", location.ARG, Short),
			New("El Rodeo", location.ARG, Short),

			New("Camino a La Puerta", location.ARG, Long),
			New("Valle de los puentes a la inversa", location.ARG, Long),
			New("Camino de acantilados y rocas inverso", location.ARG, Short),
			New("Camino a Coneta", location.ARG, Short),
			New("Huillaprima", location.ARG, Short),
			New("La Merced", location.ARG, Short),
		}
	case location.AUS:
		return []Model{
			New("Mount Kaye Pass", location.AUS, Long),
			New("Chandlers Creek", location.AUS, Long),
			New("Bondi Forest", location.AUS, Short),
			New("Rockton Plains", location.AUS, Short),
			New("Yambulla Mountain Ascent", location.AUS, Short),
			New("Noorinbee Ridge Ascent", location.AUS, Short),

			New("Mount Kaye Pass Reverse", location.AUS, Long),
			New("Chandlers Creek Reverse", location.AUS, Long),
			New("Taylor Farm Sprint", location.AUS, Short),
			New("Rockton Plains Reverse", location.AUS, Short),
			New("Yambulla Mountain Descent", location.AUS, Short),
			New("Noorinbee Ridge Descent", location.AUS, Short),
		}
	case location.FIN:
		return []Model{New("Kakaristo", location.FIN, Long),
			New("Kontinjärvi", location.FIN, Long),
			New("Kotajärvi", location.FIN, Short),
			New("Iso Oksjärvi", location.FIN, Short),
			New("Kailajärvi", location.FIN, Short),
			New("Naarajärvi", location.FIN, Short),

			New("Pitkäjärvi", location.FIN, Long),
			New("Hämelahti", location.FIN, Long),
			New("Oksala", location.FIN, Short),
			New("Järvenkylä", location.FIN, Short),
			New("Jyrkysjärvi", location.FIN, Short),
			New("Paskuri", location.FIN, Short)}
	case location.DEU:
		return []Model{New("Oberstein", location.DEU, Long),
			New("Hammerstein", location.DEU, Long),
			New("Kreuzungsring", location.DEU, Short),
			New("Verbundsring", location.DEU, Short),
			New("Innerer Feld-Sprint", location.DEU, Short),
			New("Waldaufstieg", location.DEU, Short),

			New("Frauenberg", location.DEU, Long),
			New("Ruschberg", location.DEU, Long),
			New("Kreuzungsring reverse", location.DEU, Short),
			New("Verbundsring Reverse", location.DEU, Short),
			New("Innerer Feld-Sprint (umgekehrt)", location.DEU, Short),
			New("Waldabstieg", location.DEU, Short)}
	case location.GRC:
		return []Model{New("Anodou Farmakas", location.GRC, Long),
			New("Pomona Érixi", location.GRC, Long),
			New("Koryfi Dafni", location.GRC, Short),
			New("Perasma Platani", location.GRC, Short),
			New("Ourea Spevsi", location.GRC, Short),
			New("Abies Koiláda", location.GRC, Short),

			New("Kathodo Leontiou", location.GRC, Long),
			New("Fourkéta Kourva", location.GRC, Long),
			New("Ampelonas Ormi", location.GRC, Short),
			New("Tsiristra Théa", location.GRC, Short),
			New("Pedines Epidaxi", location.GRC, Short),
			New("Ypsona tou Dasos", location.GRC, Short)}
	case location.MCO:
		return []Model{New("Vallée descendante", location.MCO, Long),
			New("Pra d’Alart", location.MCO, Long),
			New("Col de Turini - Départ en descente", location.MCO, Short),
			New("Gordolon - Courte montée", location.MCO, Short),
			New("Col de Turini sprint en montée", location.MCO, Short),
			New("Route de Turini Descente", location.MCO, Short),

			New("Route de Turini", location.MCO, Long),
			New("Col de Turini Départ", location.MCO, Long),
			New("Route de Turini Montée", location.MCO, Short),
			New("Col de Turini - Descente", location.MCO, Short),
			New("Col de Turini - Sprint en descente", location.MCO, Short),
			New("Approche du Col de Turini - Montée", location.MCO, Short),
		}
	case location.NZL:
		return []Model{
			New("Waimarama Point Forward", location.NZL, Long),
			New("Te Awanga Forward", location.NZL, Long),
			New("Waimarama Sprint Forward", location.NZL, Short),
			New("Elsthorpe Sprint Forward", location.NZL, Short),
			New("Ocean Beach Sprint Forward", location.NZL, Short),
			New("Te Awanga Sprint Forward", location.NZL, Short),

			New("Waimarama Point Reverse", location.NZL, Long),
			New("Ocean Beach", location.NZL, Long),
			New("Waimarama Sprint Reverse", location.NZL, Short),
			New("Elsthorpe Sprint Reverse", location.NZL, Short),
			New("Ocean Beach Sprint Reverse", location.NZL, Short),
			New("Te Awanga Sprint Reverse", location.NZL, Short),
		}
	case location.POL:
		return []Model{
			New("Zaróbka", location.POL, Long),
			New("Zienki", location.POL, Long),
			New("Marynka", location.POL, Short),
			New("Kopina", location.POL, Short),
			New("Lejno", location.POL, Short),
			New("Czarny Las", location.POL, Short),

			New("Zagórze", location.POL, Long),
			New("Jezioro Rotcze", location.POL, Long),
			New("Borysik", location.POL, Short),
			New("Józefin", location.POL, Short),
			New("Jagodno", location.POL, Short),
			New("Jezioro Lukie", location.POL, Short),
		}
	case location.SCO:
		return []Model{New("Newhouse Bridge", location.SCO, Long),
			New("South Morningside", location.SCO, Long),
			New("Annbank Station", location.SCO, Short),
			New("Rosebank Farm", location.SCO, Short),
			New("Old Butterstone Muir", location.SCO, Short),
			New("Glencastle Farm", location.SCO, Short),

			New("Newhouse Bridge Reverse", location.SCO, Long),
			New("South Morningside Reverse", location.SCO, Long),
			New("Annbank Station Reverse", location.SCO, Short),
			New("Rosebank Farm Reverse", location.SCO, Short),
			New("Old Butterstone Muir Reverse", location.SCO, Short),
			New("Glencastle Farm Reverse", location.SCO, Short),
		}
	case location.ESP:
		return []Model{
			New("Comienzo De Bellriu", location.ESP, Long),
			New("Centenera", location.ESP, Long),
			New("Ascenso por valle el Gualet", location.ESP, Short),
			New("Viñedos dentro del valle Parra", location.ESP, Short),
			New("Viñedos Dardenyà", location.ESP, Short),
			New("Descenso por carretera", location.ESP, Short),

			New("Final de Bellriu", location.ESP, Long),
			New("Camino a Centenera", location.ESP, Long),
			New("Salida desde Montverd", location.ESP, Short),
			New("Ascenso bosque Montverd", location.ESP, Short),
			New("Viñedos Dardenyà inversa", location.ESP, Short),
			New("Subida por carretera", location.ESP, Short),
		}
	case location.SWE:
		return []Model{
			New("Hamra", location.SWE, Long),
			New("Ransbysäter", location.SWE, Long),
			New("Elgsjön", location.SWE, Short),
			New("Stor-jangen Sprint", location.SWE, Short),
			New("Älgsjön Sprint", location.SWE, Short),
			New("Östra Hinnsjön", location.SWE, Short),

			New("Lysvik", location.SWE, Long),
			New("Norraskoga", location.SWE, Long),
			New("Älgsjön", location.SWE, Short),
			New("Stor-jangen Sprint Reverse", location.SWE, Short),
			New("Skogsrallyt", location.SWE, Short),
			New("Björklangen", location.SWE, Short),
		}
	case location.USA:
		return []Model{
			New("Beaver Creek Trail Forward", location.USA, Long),
			New("North Fork Pass", location.USA, Long),
			New("Hancock Creek Burst", location.USA, Short),
			New("Fuller Mountain Ascent", location.USA, Short),
			New("Tolt Valley Sprint Forward", location.USA, Short),
			New("Hancock Hill Sprint Forward", location.USA, Short),

			New("Beaver Creek Trail Reverse", location.USA, Long),
			New("North Fork Pass Reverse", location.USA, Long),
			New("Fury Lake Depart", location.USA, Short),
			New("Fuller Mountain Descent", location.USA, Short),
			New("Tolt Valley Sprint Reverse", location.USA, Short),
			New("Hancock Hill Sprint Reverse", location.USA, Short),
		}
	case location.WAL:
		return []Model{
			New("River Severn Valley", location.WAL, Long),
			New("Sweet Lamb", location.WAL, Long),
			New("Fferm Wynt", location.WAL, Short),
			New("Dyffryn Afon", location.WAL, Short),
			New("Bidno Moorland", location.WAL, Short),
			New("Pant Mawr", location.WAL, Short),

			New("Bronfelen", location.WAL, Long),
			New("Geufron Forest", location.WAL, Long),
			New("Fferm Wynt Reverse", location.WAL, Short),
			New("Dyffryn Afon Reverse", location.WAL, Short),
			New("Bidno Moorland Reverse", location.WAL, Short),
			New("Pant Mawr Reverse", location.WAL, Short)}
	}
	return []Model{}
}
