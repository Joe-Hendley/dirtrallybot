package builder

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
)

func Stages() map[stage.Model]int {
	stages := make(map[stage.Model]int)

	for _, l := range location.List() {
		for _, s := range StagesAtLocation(l) {
			stages[s] = 0
		}
	}

	return stages
}

func StagesAtLocation(l location.Model) []stage.Model {
	switch l {
	case location.ARG:
		return []stage.Model{
			stage.New("Las Juntas", location.ARG, stage.Long),
			stage.New("Valle de los puentes", location.ARG, stage.Long),
			stage.New("Camino de acantilados y rocas", location.ARG, stage.Short),
			stage.New("San Isidro", location.ARG, stage.Short),
			stage.New("Miraflores", location.ARG, stage.Short),
			stage.New("El Rodeo", location.ARG, stage.Short),

			stage.New("Camino a La Puerta", location.ARG, stage.Long),
			stage.New("Valle de los puentes a la inversa", location.ARG, stage.Long),
			stage.New("Camino de acantilados y rocas inverso", location.ARG, stage.Short),
			stage.New("Camino a Coneta", location.ARG, stage.Short),
			stage.New("Huillaprima", location.ARG, stage.Short),
			stage.New("La Merced", location.ARG, stage.Short),
		}
	case location.AUS:
		return []stage.Model{
			stage.New("Mount Kaye Pass", location.AUS, stage.Long),
			stage.New("Chandlers Creek", location.AUS, stage.Long),
			stage.New("Bondi Forest", location.AUS, stage.Short),
			stage.New("Rockton Plains", location.AUS, stage.Short),
			stage.New("Yambulla Mountain Ascent", location.AUS, stage.Short),
			stage.New("Noorinbee Ridge Ascent", location.AUS, stage.Short),

			stage.New("Mount Kaye Pass Reverse", location.AUS, stage.Long),
			stage.New("Chandlers Creek Reverse", location.AUS, stage.Long),
			stage.New("Taylor Farm Sprint", location.AUS, stage.Short),
			stage.New("Rockton Plains Reverse", location.AUS, stage.Short),
			stage.New("Yambulla Mountain Descent", location.AUS, stage.Short),
			stage.New("Noorinbee Ridge Descent", location.AUS, stage.Short),
		}
	case location.FIN:
		return []stage.Model{stage.New("Kakaristo", location.FIN, stage.Long),
			stage.New("Kontinjärvi", location.FIN, stage.Long),
			stage.New("Kotajärvi", location.FIN, stage.Short),
			stage.New("Iso Oksjärvi", location.FIN, stage.Short),
			stage.New("Kailajärvi", location.FIN, stage.Short),
			stage.New("Naarajärvi", location.FIN, stage.Short),

			stage.New("Pitkäjärvi", location.FIN, stage.Long),
			stage.New("Hämelahti", location.FIN, stage.Long),
			stage.New("Oksala", location.FIN, stage.Short),
			stage.New("Järvenkylä", location.FIN, stage.Short),
			stage.New("Jyrkysjärvi", location.FIN, stage.Short),
			stage.New("Paskuri", location.FIN, stage.Short)}
	case location.DEU:
		return []stage.Model{stage.New("Oberstein", location.DEU, stage.Long),
			stage.New("Hammerstein", location.DEU, stage.Long),
			stage.New("Kreuzungsring", location.DEU, stage.Short),
			stage.New("Verbundsring", location.DEU, stage.Short),
			stage.New("Innerer Feld-Sprint", location.DEU, stage.Short),
			stage.New("Waldaufstieg", location.DEU, stage.Short),

			stage.New("Frauenberg", location.DEU, stage.Long),
			stage.New("Ruschberg", location.DEU, stage.Long),
			stage.New("Kreuzungsring reverse", location.DEU, stage.Short),
			stage.New("Verbundsring Reverse", location.DEU, stage.Short),
			stage.New("Innerer Feld-Sprint (umgekehrt)", location.DEU, stage.Short),
			stage.New("Waldabstieg", location.DEU, stage.Short)}
	case location.GRC:
		return []stage.Model{stage.New("Anodou Farmakas", location.GRC, stage.Long),
			stage.New("Pomona Érixi", location.GRC, stage.Long),
			stage.New("Koryfi Dafni", location.GRC, stage.Short),
			stage.New("Perasma Platani", location.GRC, stage.Short),
			stage.New("Ourea Spevsi", location.GRC, stage.Short),
			stage.New("Abies Koiláda", location.GRC, stage.Short),

			stage.New("Kathodo Leontiou", location.GRC, stage.Long),
			stage.New("Fourkéta Kourva", location.GRC, stage.Long),
			stage.New("Ampelonas Ormi", location.GRC, stage.Short),
			stage.New("Tsiristra Théa", location.GRC, stage.Short),
			stage.New("Pedines Epidaxi", location.GRC, stage.Short),
			stage.New("Ypsona tou Dasos", location.GRC, stage.Short)}
	case location.MCO:
		return []stage.Model{stage.New("Vallée descendante", location.MCO, stage.Long),
			stage.New("Pra d’Alart", location.MCO, stage.Long),
			stage.New("Col de Turini - Départ en descente", location.MCO, stage.Short),
			stage.New("Gordolon - Courte montée", location.MCO, stage.Short),
			stage.New("Col de Turini sprint en montée", location.MCO, stage.Short),
			stage.New("Route de Turini Descente", location.MCO, stage.Short),

			stage.New("Route de Turini", location.MCO, stage.Long),
			stage.New("Col de Turini Départ", location.MCO, stage.Long),
			stage.New("Route de Turini Montée", location.MCO, stage.Short),
			stage.New("Col de Turini - Descente", location.MCO, stage.Short),
			stage.New("Col de Turini - Sprint en descente", location.MCO, stage.Short),
			stage.New("Approche du Col de Turini - Montée", location.MCO, stage.Short),
		}
	case location.NZL:
		return []stage.Model{
			stage.New("Waimarama Point Forward", location.NZL, stage.Long),
			stage.New("Te Awanga Forward", location.NZL, stage.Long),
			stage.New("Waimarama Sprint Forward", location.NZL, stage.Short),
			stage.New("Elsthorpe Sprint Forward", location.NZL, stage.Short),
			stage.New("Ocean Beach Sprint Forward", location.NZL, stage.Short),
			stage.New("Te Awanga Sprint Forward", location.NZL, stage.Short),

			stage.New("Waimarama Point Reverse", location.NZL, stage.Long),
			stage.New("Ocean Beach", location.NZL, stage.Long),
			stage.New("Waimarama Sprint Reverse", location.NZL, stage.Short),
			stage.New("Elsthorpe Sprint Reverse", location.NZL, stage.Short),
			stage.New("Ocean Beach Sprint Reverse", location.NZL, stage.Short),
			stage.New("Te Awanga Sprint Reverse", location.NZL, stage.Short),
		}
	case location.POL:
		return []stage.Model{
			stage.New("Zaróbka", location.POL, stage.Long),
			stage.New("Zienki", location.POL, stage.Long),
			stage.New("Marynka", location.POL, stage.Short),
			stage.New("Kopina", location.POL, stage.Short),
			stage.New("Lejno", location.POL, stage.Short),
			stage.New("Czarny Las", location.POL, stage.Short),

			stage.New("Zagórze", location.POL, stage.Long),
			stage.New("Jezioro Rotcze", location.POL, stage.Long),
			stage.New("Borysik", location.POL, stage.Short),
			stage.New("Józefin", location.POL, stage.Short),
			stage.New("Jagodno", location.POL, stage.Short),
			stage.New("Jezioro Lukie", location.POL, stage.Short),
		}
	case location.SCO:
		return []stage.Model{stage.New("Newhouse Bridge", location.SCO, stage.Long),
			stage.New("South Morningside", location.SCO, stage.Long),
			stage.New("Annbank Station", location.SCO, stage.Short),
			stage.New("Rosebank Farm", location.SCO, stage.Short),
			stage.New("Old Butterstone Muir", location.SCO, stage.Short),
			stage.New("Glencastle Farm", location.SCO, stage.Short),

			stage.New("Newhouse Bridge Reverse", location.SCO, stage.Long),
			stage.New("South Morningside Reverse", location.SCO, stage.Long),
			stage.New("Annbank Station Reverse", location.SCO, stage.Short),
			stage.New("Rosebank Farm Reverse", location.SCO, stage.Short),
			stage.New("Old Butterstone Muir Reverse", location.SCO, stage.Short),
			stage.New("Glencastle Farm Reverse", location.SCO, stage.Short),
		}
	case location.ESP:
		return []stage.Model{
			stage.New("Comienzo De Bellriu", location.ESP, stage.Long),
			stage.New("Centenera", location.ESP, stage.Long),
			stage.New("Ascenso por valle el Gualet", location.ESP, stage.Short),
			stage.New("Viñedos dentro del valle Parra", location.ESP, stage.Short),
			stage.New("Viñedos Dardenyà", location.ESP, stage.Short),
			stage.New("Descenso por carretera", location.ESP, stage.Short),

			stage.New("Final de Bellriu", location.ESP, stage.Long),
			stage.New("Camino a Centenera", location.ESP, stage.Long),
			stage.New("Salida desde Montverd", location.ESP, stage.Short),
			stage.New("Ascenso bosque Montverd", location.ESP, stage.Short),
			stage.New("Viñedos Dardenyà inversa", location.ESP, stage.Short),
			stage.New("Subida por carretera", location.ESP, stage.Short),
		}
	case location.SWE:
		return []stage.Model{
			stage.New("Hamra", location.SWE, stage.Long),
			stage.New("Ransbysäter", location.SWE, stage.Long),
			stage.New("Elgsjön", location.SWE, stage.Short),
			stage.New("Stor-jangen Sprint", location.SWE, stage.Short),
			stage.New("Älgsjön Sprint", location.SWE, stage.Short),
			stage.New("Östra Hinnsjön", location.SWE, stage.Short),

			stage.New("Lysvik", location.SWE, stage.Long),
			stage.New("Norraskoga", location.SWE, stage.Long),
			stage.New("Älgsjön", location.SWE, stage.Short),
			stage.New("Stor-jangen Sprint Reverse", location.SWE, stage.Short),
			stage.New("Skogsrallyt", location.SWE, stage.Short),
			stage.New("Björklangen", location.SWE, stage.Short),
		}
	case location.USA:
		return []stage.Model{
			stage.New("Beaver Creek Trail Forward", location.USA, stage.Long),
			stage.New("North Fork Pass", location.USA, stage.Long),
			stage.New("Hancock Creek Burst", location.USA, stage.Short),
			stage.New("Fuller Mountain Ascent", location.USA, stage.Short),
			stage.New("Tolt Valley Sprint Forward", location.USA, stage.Short),
			stage.New("Hancock Hill Sprint Forward", location.USA, stage.Short),

			stage.New("Beaver Creek Trail Reverse", location.USA, stage.Long),
			stage.New("North Fork Pass Reverse", location.USA, stage.Long),
			stage.New("Fury Lake Depart", location.USA, stage.Short),
			stage.New("Fuller Mountain Descent", location.USA, stage.Short),
			stage.New("Tolt Valley Sprint Reverse", location.USA, stage.Short),
			stage.New("Hancock Hill Sprint Reverse", location.USA, stage.Short),
		}
	case location.WAL:
		return []stage.Model{
			stage.New("River Severn Valley", location.WAL, stage.Long),
			stage.New("Sweet Lamb", location.WAL, stage.Long),
			stage.New("Fferm Wynt", location.WAL, stage.Short),
			stage.New("Dyffryn Afon", location.WAL, stage.Short),
			stage.New("Bidno Moorland", location.WAL, stage.Short),
			stage.New("Pant Mawr", location.WAL, stage.Short),

			stage.New("Bronfelen", location.WAL, stage.Long),
			stage.New("Geufron Forest", location.WAL, stage.Long),
			stage.New("Fferm Wynt Reverse", location.WAL, stage.Short),
			stage.New("Dyffryn Afon Reverse", location.WAL, stage.Short),
			stage.New("Bidno Moorland Reverse", location.WAL, stage.Short),
			stage.New("Pant Mawr Reverse", location.WAL, stage.Short)}
	}
	return []stage.Model{}
}
