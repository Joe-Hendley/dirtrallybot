package challenge

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/bwmarrin/discordgo"
)

// TODO:
// rewrite the entire custom ID system

const (
	RandomID = "random"
	DryID    = "dry"
	WetID    = "wet"

	baseMessage = "%s Challenge Builder v0.0.1"

	gameIndex      = 1
	componentIndex = 2

	locationID               = "location"
	stageID                  = "stage"
	weatherID                = "weather"
	SubmitLocationAndStageID = "submit1"

	drivetrainID = "drivetrain"
	classID      = "class"
	carID        = "car"
	SubmitCarID  = "submit2"
)

func gameIDString(config challenge.Config) string {
	switch config.Game {
	case game.DR2:
		return "dr2"
	case game.WRC:
		return "wrc"
	}
	return ""
}

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
					CustomID: strings.Join([]string{ChallengeID, gameIDString(config), SubmitLocationAndStageID}, idFieldDelimiter),
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
					CustomID: strings.Join([]string{ChallengeID, gameIDString(config), SubmitCarID}, idFieldDelimiter),
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

	for _, loc := range location.List(config.Game) {
		locID := strings.ToLower(loc.String())
		if locID == selected {
			hasDefault = true
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       loc.String(),
			Value:       locID,
			Emoji:       &discordgo.ComponentEmoji{Name: loc.Flag()},
			Description: loc.DetailedString(),
			Default:     locID == selected,
		})
	}

	if !hasDefault {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: "Location",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), locationID}, idFieldDelimiter),
		Options:     options,
	}
}

func buildStageMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Stage != nil {
		selected = strings.ToLower(config.Stage.Name())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Stage"),
	}

	hasDefault := false

	if config.Location != nil {
		for _, stage := range stage.AtLocation(*config.Location) {
			stageID := strings.ToLower(stage.Name())
			if stageID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:       stage.Name(),
				Value:       stageID,
				Emoji:       &discordgo.ComponentEmoji{Name: stage.Distance().Emoji()},
				Description: stage.String(),
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
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), stageID}, idFieldDelimiter),
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
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), weatherID}, idFieldDelimiter),
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

	for _, drivetrain := range drivetrain.List(config.Game) {
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
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), drivetrainID}, idFieldDelimiter),
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
		for _, class := range class.WithDrivetrain(*config.Drivetrain, config.Game) {
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
		for _, class := range class.List(config.Game) {
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
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), classID}, idFieldDelimiter),
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildCarMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Car != nil {
		selected = strings.ToLower(config.Car.Name())
	}

	options := []discordgo.SelectMenuOption{
		randomOption("Car"),
	}

	hasDefault := false

	if config.Class != nil {
		for _, car := range car.InClass(*config.Class, config.Game) {
			carID := strings.ToLower(car.Name())
			if carID == selected {
				hasDefault = true
			}

			options = append(options, discordgo.SelectMenuOption{
				Label:   car.Name(),
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
		CustomID:    strings.Join([]string{ChallengeID, gameIDString(config), carID}, idFieldDelimiter),
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildStageConfigFromInteraction(interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	customIDFields := strings.Split(customID, idFieldDelimiter)
	if len(customIDFields) != 3 {
		return challenge.Config{}, fmt.Errorf("unexpected customID %s", customID)
	}

	gameID := customIDFields[gameIndex]
	changedComponentID := customIDFields[componentIndex]
	fmt.Println(customIDFields)

	var newValue string
	if changedComponentID != SubmitCarID && changedComponentID != SubmitLocationAndStageID {
		newValue = interaction.MessageComponentData().Values[0]
	}
	fmt.Println("newValue", newValue)

	whichGame := gameFromID(gameID)
	if whichGame == game.NotSet {
		return challenge.Config{}, fmt.Errorf("invalid game from customID %s", gameID)
	}
	fmt.Println("game", whichGame)

	config := challenge.Config{Game: whichGame}

	componentValues := map[string]string{}

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
					componentID := strings.Split(component.CustomID, idFieldDelimiter)[componentIndex]
					componentValues[componentID] = option.Value
					fmt.Println("found default", component.CustomID, option.Value)
					break
				}
			}
		default:
			return challenge.Config{}, errors.New("unexpected component type")
		}
	}

	fmt.Println(componentValues)

	config = applyLocation(config, componentValues[locationID])
	config = applyStage(config, componentValues[stageID])
	config = applyWeather(config, componentValues[weatherID])

	fmt.Println(config)

	fmt.Println("changedID", changedComponentID)
	switch changedComponentID {
	case locationID:
		config = applyLocation(config, newValue)
	case stageID:
		config = applyStage(config, newValue)
	case weatherID:
		config = applyWeather(config, newValue)
	}

	return config, nil
}

func buildCarConfigFromInteraction(interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	customIDFields := strings.Split(customID, idFieldDelimiter)
	if len(customIDFields) != 3 {
		return challenge.Config{}, fmt.Errorf("unexpected customID %s", customID)
	}

	gameID := customIDFields[gameIndex]
	changedComponentID := customIDFields[componentIndex]

	var newValue string
	if changedComponentID != SubmitCarID && changedComponentID != SubmitLocationAndStageID {
		newValue = interaction.MessageComponentData().Values[0]
	}

	whichGame := gameFromID(gameID)
	if whichGame == game.NotSet {
		return challenge.Config{}, fmt.Errorf("invalid game from customID %s", gameID)
	}

	config := challenge.Config{Game: whichGame}

	componentValues := map[string]string{}

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
					componentID := strings.Split(component.CustomID, idFieldDelimiter)[componentIndex]
					componentValues[componentID] = option.Value
					break
				}
			}
		default:
			return challenge.Config{}, errors.New("unexpected component type")
		}
	}

	lines := strings.Split(interaction.Message.Content, "\n")
	for _, line := range lines {
		emojiDelimited := strings.Split(line, challenge.EmojiDelimiter)
		if len(emojiDelimited) > 1 {
			switch {
			case strings.HasPrefix(strings.ToLower(emojiDelimited[0]), locationID):
				componentValues[locationID] = strings.ToLower(emojiDelimited[1])
			case strings.HasPrefix(strings.ToLower(emojiDelimited[0]), stageID):
				componentValues[stageID] = strings.ToLower(emojiDelimited[1])
			case strings.HasPrefix(strings.ToLower(emojiDelimited[0]), weatherID):
				componentValues[weatherID] = strings.ToLower(emojiDelimited[1])
			}
		}
	}

	config = applyLocation(config, componentValues[locationID])
	config = applyStage(config, componentValues[stageID])
	config = applyWeather(config, componentValues[weatherID])
	config = applyDrivetrain(config, componentValues[drivetrainID])
	config = applyClass(config, componentValues[classID])
	config = applyCar(config, componentValues[carID])

	switch changedComponentID {
	case drivetrainID:
		config = applyDrivetrain(config, newValue)
	case classID:
		config = applyClass(config, newValue)
	case carID:
		config = applyCar(config, newValue)
	}

	return config, nil
}

func gameFromID(gameID string) game.Model {
	switch gameID {
	case DR2ID:
		return game.DR2
	case WRCID:
		return game.WRC
	}
	return game.NotSet
}

func applyLocation(config challenge.Config, value string) challenge.Config {
	fmt.Println("applying loc")
	if value == RandomID {
		config.Stage = nil
		config.Weather = nil
		return config
	}

	for _, loc := range location.List(config.Game) {
		fmt.Println("loc:", loc)
		if value == strings.ToLower(loc.String()) {
			config.Location = &loc
			return config
		}
	}

	config.Location = nil

	return config
}

func applyStage(config challenge.Config, value string) challenge.Config {
	fmt.Println("applying stage", value)
	fmt.Println("config", config)
	if value == RandomID || config.Location == nil {
		fmt.Println("resetting config")
		return config
	}

	for _, stage := range stage.AtLocation(*config.Location) {
		if value == strings.ToLower(stage.Name()) {
			config.Stage = &stage
			return config
		}
	}

	config.Stage = nil

	return config
}

func applyWeather(config challenge.Config, value string) challenge.Config {
	fmt.Println("applying weather")
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

	for _, drivetrain := range drivetrain.List(config.Game) {
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

	for _, class := range class.List(config.Game) {
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

	for _, car := range car.InClass(*config.Class, config.Game) {
		if value == strings.ToLower(car.Name()) {
			config.Car = &car
			return config
		}
	}

	config.Car = nil

	return config
}
