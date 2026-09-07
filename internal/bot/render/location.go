package render

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// Flag is the flag (or stand-in) emoji for a location.
func Flag(m location.Model) string {
	switch m {
	// DR2
	case location.ARG:
		return "🇦🇷"
	case location.AUS:
		return "🇦🇺"
	case location.FIN:
		return "🇫🇮"
	case location.DEU:
		return "🇩🇪"
	case location.GRC:
		return "🇬🇷"
	case location.MCO:
		return "🇲🇨"
	case location.NZL:
		return "🇳🇿"
	case location.POL:
		return "🇵🇱"
	case location.SCO:
		return "🏴󠁧󠁢󠁳󠁣󠁴󠁿"
	case location.ESP:
		return "🇪🇸"
	case location.SWE:
		return "🇸🇪"
	case location.USA:
		return "🇺🇸"
	case location.WAL:
		return "🏴󠁧󠁢󠁷󠁬󠁳󠁿"

	// WRC
	case location.MCO_WRC:
		return "🇲🇨"
	case location.SWE_WRC:
		return "🇸🇪"
	case location.MEX:
		return "🇲🇽"
	case location.HRV:
		return "🇭🇷"
	case location.PRT:
		return "🇵🇹"
	case location.ITA:
		return "🇮🇹"
	case location.KEN:
		return "🇰🇪"
	case location.EST:
		return "🇪🇪"
	case location.FIN_WRC:
		return "🇫🇮"
	case location.GRC_WRC:
		return "🇬🇷"
	case location.CHL:
		return "🇨🇱"
	case location.CER:
		return "🇪🇺"
	case location.JPN:
		return "🇯🇵"
	case location.MED:
		return "🫒"
	case location.PAC:
		return "🗿"
	case location.UUU:
		return "🦘"
	case location.SCA:
		return "🌲"
	case location.IBE:
		return "🐮"
	}
	return "invalid location"
}

// WeatherDescription is the human-readable weather line shown for a challenge's
// location and weather combination.
func WeatherDescription(l location.Model, w weather.Model) string {
	return weatherDescriptions[l][w]
}

const (
	DAYCLEARDRY        = "☀️ **Daytime / Clear / Dry Surface**"
	DAYCLOUDYWET       = "⛅💧 **Daytime / Cloudy / Wet Surface**"
	DAYCLOUDYSNOW      = "⛅❄️ **Daytime / Cloudy / Snow**"
	DAYHEAVYRAINWET    = "⛅🌧️ **Daytime / Heavy Rain / Wet Surface**"
	DUSKCLOUDYWET      = "🌆☁️💧 **Dusk / Cloudy / Wet Surface**"
	DUSKHEAVYRAINWET   = "🌆🌧️ **Dusk / Heavy Rain / Wet Surface**"
	SUNSETCLOUDYWET    = "🌇☁️💧 **Sunset / Cloudy / Wet Surface**"
	SUNSETHEAVYRAINWET = "🌇🌧️ **Sunset / Heavy Rain / Wet Surface**"
)

var weatherDescriptions = map[location.Model]map[weather.Model]string{
	// DR2
	location.ARG: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DUSKHEAVYRAINWET,
	},
	location.AUS: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.FIN: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DUSKCLOUDYWET,
	},
	location.DEU: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYHEAVYRAINWET,
	},
	location.GRC: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: SUNSETHEAVYRAINWET,
	},
	location.MCO: {
		weather.DRY: DAYCLEARDRY,
	},
	location.NZL: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.POL: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: SUNSETCLOUDYWET,
	},
	location.SCO: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.ESP: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.SWE: {
		weather.SNOW: DAYCLOUDYSNOW,
	},
	location.USA: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.WAL: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: SUNSETCLOUDYWET,
	},

	// WRC
	location.MCO_WRC: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.SWE_WRC: {
		weather.SNOW: DAYCLOUDYSNOW,
	},
	location.MEX: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.HRV: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.PRT: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.ITA: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.KEN: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.EST: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.FIN_WRC: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.GRC_WRC: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.CHL: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.CER: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.JPN: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.MED: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.PAC: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.UUU: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
	location.SCA: {
		weather.SNOW: DAYCLOUDYSNOW,
	},
	location.IBE: {
		weather.DRY: DAYCLEARDRY,
		weather.WET: DAYCLOUDYWET,
	},
}
