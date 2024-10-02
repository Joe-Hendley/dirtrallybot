package challenge

import (
	"errors"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/bwmarrin/discordgo"
)

// so, we get the interaction in with the following options:
// location - fully random, specific location, specific stage
// weather 	- fully random, specific weather
// car		- fully random, specific drivetrain, specific class, specific car

// we then need to populate the select options with the list of options for these

const (
	RandomID = "random"
	DryID    = "dry"
	WetID    = "wet"

	DR2ChallengePrefix = "challenge-dr2-"

	locationFieldID          = "location"
	LocationSelectID         = DR2ChallengePrefix + locationFieldID
	stageFieldID             = "stage"
	StageSelectID            = DR2ChallengePrefix + stageFieldID
	weatherFieldID           = "weather"
	WeatherSelectID          = DR2ChallengePrefix + weatherFieldID
	SubmitLocationAndStageID = DR2ChallengePrefix + "submit-1"

	drivetrainFieldID  = "drivetrain"
	DrivetrainSelectID = DR2ChallengePrefix + drivetrainFieldID
	classFieldID       = "class"
	ClassSelectID      = DR2ChallengePrefix + classFieldID
	carFieldID         = "car"
	CarSelectID        = DR2ChallengePrefix + carFieldID
	SubmitCarID        = DR2ChallengePrefix + "submit-2"
)

func buildChallengeLocationMessageComponents(config challenge.Config) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildLocationsMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildStageMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildWeatherMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Submit Stage",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					CustomID: SubmitLocationAndStageID,
				},
			},
		},
	}
}

func buildChallengeCarMessageComponents(config challenge.Config) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildDriveTrainMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildClassMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buildCarMenu(config),
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Submit",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					CustomID: SubmitCarID,
				},
			},
		},
	}
}

func randomOption(category string) discordgo.SelectMenuOption {
	return discordgo.SelectMenuOption{
		Label: "Random " + category, Value: RandomID, Emoji: &discordgo.ComponentEmoji{Name: challenge.RandomEmoji},
	}
}

func buildLocationsMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Location != nil {
		selected = strings.ToLower(config.Location.String())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Location"),
	}

	hasDefault := false

	for _, loc := range location.List() {
		locID := strings.ToLower(loc.String())
		if locID == selected {
			hasDefault = true
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       loc.String(),
			Value:       locID,
			Emoji:       &discordgo.ComponentEmoji{Name: loc.Flag()},
			Description: loc.LongString(),
			Default:     locID == selected,
		})
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Location",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    LocationSelectID,
		Options:     options,
	}
}

func buildStageMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Stage != nil {
		selected = strings.ToLower(config.Stage.String())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Stage"),
	}

	hasDefault := false

	if config.Location != nil {
		for _, stage := range stage.AtLocation(*config.Location) {
			stageID := strings.ToLower(stage.String())
			if stageID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:       stage.String(),
				Value:       stageID,
				Description: stage.LongString(),
				Default:     stageID == selected,
			})
		}
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Stage",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    StageSelectID,
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildWeatherMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Weather != nil {
		selected = strings.ToLower(config.Weather.String())
	}

	var options []discordgo.SelectMenuOption
	switch {
	case config.Location == nil:
		dryID := strings.ToLower(weather.DRY.String())
		wetID := strings.ToLower(weather.WET.String())
		options = []discordgo.SelectMenuOption{
			randomOption("Weather"),
			{
				Label:   weather.DRY.String(),
				Value:   dryID,
				Emoji:   &discordgo.ComponentEmoji{Name: weather.DRY.Emoji()},
				Default: dryID == selected,
			},
			{
				Label:   weather.WET.String(),
				Value:   strings.ToLower(weather.WET.String()),
				Emoji:   &discordgo.ComponentEmoji{Name: weather.WET.Emoji()},
				Default: wetID == selected,
			},
		}

		if selected != dryID && selected != wetID {
			options[0].Default = true
		}

	case len(config.Location.Weather()) == 1:
		options = []discordgo.SelectMenuOption{
			{
				Label:   config.Location.Weather()[0].String(),
				Value:   strings.ToLower(config.Location.Weather()[0].String()),
				Emoji:   &discordgo.ComponentEmoji{Name: config.Location.Weather()[0].Emoji()},
				Default: true,
			},
		}

	default:
		options = []discordgo.SelectMenuOption{
			randomOption("Weather"),
		}

		hasDefault := false

		for _, weather := range config.Location.Weather() {
			weatherID := strings.ToLower(weather.String())
			if weatherID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:   weather.String(),
				Value:   weatherID,
				Emoji:   &discordgo.ComponentEmoji{Name: weather.Emoji()},
				Default: weatherID == selected,
			})
		}

		if !hasDefault {
			options[0].Default = true
		}
	}

	return discordgo.SelectMenu{
		Placeholder: "Weather",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    WeatherSelectID,
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildDriveTrainMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Drivetrain != nil {
		selected = strings.ToLower(config.Drivetrain.String())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Drivetrain"),
	}

	hasDefault := false

	for _, drivetrain := range drivetrain.List() {
		drivetrainID := strings.ToLower(drivetrain.String())
		if drivetrainID == selected {
			hasDefault = true
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:   drivetrain.String(),
			Value:   drivetrainID,
			Emoji:   &discordgo.ComponentEmoji{Name: drivetrain.Emoji()},
			Default: drivetrainID == selected,
		})
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Drivetrain",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    DrivetrainSelectID,
		Options:     options,
	}
}

func buildClassMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Class != nil {
		selected = strings.ToLower(config.Class.String())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Class"),
	}

	hasDefault := false

	if config.Drivetrain != nil {
		for _, class := range class.WithDrivetrain(*config.Drivetrain) {
			classID := strings.ToLower(class.String())
			if classID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:   class.String(),
				Value:   classID,
				Default: classID == selected,
			})
		}
	} else {
		for _, class := range class.List() {
			classID := strings.ToLower(class.String())
			if classID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:   class.String(),
				Value:   classID,
				Default: classID == selected,
			})
		}
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Class",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    ClassSelectID,
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildCarMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Car != nil {
		selected = strings.ToLower(config.Car.String())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Car"),
	}

	hasDefault := false

	if config.Class != nil {
		for _, car := range car.InClass(*config.Class) {
			carID := strings.ToLower(car.String())
			if carID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:   car.String(),
				Value:   carID,
				Default: carID == selected,
			})
		}
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Car",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    CarSelectID,
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildStageConfigFromInteraction(interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	var newValue string
	if customID != SubmitCarID && customID != SubmitLocationAndStageID {
		newValue = interaction.MessageComponentData().Values[0]
	}
	config := challenge.Config{}

	currentValues := map[string]string{}

	// each component is contained in a separate action row
	// should be some number of select menus then a button
	for _, componentRow := range interaction.Message.Components {
		actionsRow, ok := componentRow.(*discordgo.ActionsRow)
		if !ok {
			return challenge.Config{}, errors.New("expected action row")
		}

		switch component := actionsRow.Components[0].(type) {
		case *discordgo.Button:
			continue
		case *discordgo.SelectMenu:
			for _, option := range component.Options {
				if option.Default {
					currentValues[component.CustomID] = option.Value
					break
				}
			}
		default:
			return challenge.Config{}, errors.New("unexpected component type")
		}
	}

	config = applyLocation(config, currentValues[LocationSelectID])
	config = applyStage(config, currentValues[StageSelectID])
	config = applyWeather(config, currentValues[WeatherSelectID])

	switch customID {
	case LocationSelectID:
		config = applyLocation(config, newValue)
	case StageSelectID:
		config = applyStage(config, newValue)
	case WeatherSelectID:
		config = applyWeather(config, newValue)
	}

	return config, nil
}

func buildCarConfigFromInteraction(interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	var newValue string
	if customID != SubmitCarID && customID != SubmitLocationAndStageID {
		newValue = interaction.MessageComponentData().Values[0]
	}
	config := challenge.Config{}

	currentValues := map[string]string{}

	// each component is contained in a separate action row
	// should be some number of select menus then a button
	for _, componentRow := range interaction.Message.Components {
		actionsRow, ok := componentRow.(*discordgo.ActionsRow)
		if !ok {
			return challenge.Config{}, errors.New("expected action row")
		}

		switch component := actionsRow.Components[0].(type) {
		case *discordgo.Button:
			continue
		case *discordgo.SelectMenu:
			for _, option := range component.Options {
				if option.Default {
					currentValues[component.CustomID] = option.Value
					break
				}
			}
		default:
			return challenge.Config{}, errors.New("unexpected component type")
		}
	}

	config = applyDrivetrain(config, currentValues[DrivetrainSelectID])
	config = applyClass(config, currentValues[ClassSelectID])
	config = applyCar(config, currentValues[CarSelectID])

	switch customID {
	case DrivetrainSelectID:
		config = applyDrivetrain(config, newValue)
	case ClassSelectID:
		config = applyClass(config, newValue)
	case CarSelectID:
		config = applyCar(config, newValue)
	}

	return config, nil
}

func applyLocation(config challenge.Config, value string) challenge.Config {
	if value == RandomID {
		config.Stage = nil
		config.Weather = nil
		return config
	}

	for _, loc := range location.List() {
		if value == strings.ToLower(loc.String()) {
			config.Location = &loc
			return config
		}
	}

	config.Location = nil

	return config
}

func applyStage(config challenge.Config, value string) challenge.Config {
	if value == RandomID || config.Location == nil {
		return config
	}

	for _, stage := range stage.AtLocation(*config.Location) {
		if value == strings.ToLower(stage.String()) {
			config.Stage = &stage
			return config
		}
	}

	config.Stage = nil

	return config
}

func applyWeather(config challenge.Config, value string) challenge.Config {
	if value == RandomID {
		return config
	}

	var validWeathers []weather.Model
	if config.Location != nil {
		validWeathers = config.Location.Weather()
	} else {
		validWeathers = []weather.Model{weather.DRY, weather.WET}
	}

	if len(validWeathers) == 1 {
		config.Weather = &validWeathers[0]
		return config
	}

	for _, weather := range validWeathers {
		if value == strings.ToLower(weather.String()) {
			config.Weather = &weather
			return config
		}
	}

	config.Weather = nil

	return config
}

func applyDrivetrain(config challenge.Config, value string) challenge.Config {
	if value == RandomID {
		config.Class = nil
		return config
	}

	for _, drivetrain := range drivetrain.List() {
		if value == strings.ToLower(drivetrain.String()) {
			config.Drivetrain = &drivetrain
			return config
		}
	}

	config.Drivetrain = nil

	return config
}

func applyClass(config challenge.Config, value string) challenge.Config {
	if value == RandomID {
		config.Car = nil
		return config
	}

	for _, class := range class.List() {
		if value == strings.ToLower(class.String()) {
			config.Class = &class
			drivetrain := class.Drivetrain()
			config.Drivetrain = &drivetrain
			return config
		}
	}

	config.Class = nil

	return config
}

func applyCar(config challenge.Config, value string) challenge.Config {
	if value == RandomID || config.Class == nil {
		return config
	}

	for _, car := range car.InClass(*config.Class) {
		if value == strings.ToLower(car.String()) {
			config.Car = &car
			return config
		}
	}

	config.Car = nil

	return config
}

// don't do this, just use base36 encoding instead
const fieldDelim = 0x1f

func buildParsableConfigString(config challenge.Config) string {
	stringParts := []string{}
	if config.Location != nil {
		stringParts = append(stringParts, locationFieldID+":"+strings.ToLower(config.Location.String()))
	}

	if config.Stage != nil {
		stringParts = append(stringParts, stageFieldID+":"+strings.ToLower(config.Stage.String()))
	}

	if config.Weather != nil {
		stringParts = append(stringParts, weatherFieldID+":"+strings.ToLower(config.Weather.String()))
	}

	if config.Car != nil {
		stringParts = append(stringParts, carFieldID+":"+strings.ToLower(config.Car.String()))
	}

	if config.Class != nil {
		stringParts = append(stringParts, classFieldID+":"+strings.ToLower(config.Class.String()))
	}

	if config.Drivetrain != nil {
		stringParts = append(stringParts, drivetrainFieldID+":"+strings.ToLower(config.Drivetrain.String()))
	}

	return strings.Join(stringParts, ", ")
}

func parseConfigString(configString string) challenge.Config {

}
