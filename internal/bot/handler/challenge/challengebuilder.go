package challenge

import (
	"fmt"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/render"
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

const (
	RandomID = "random"
	DryID    = "dry"
	WetID    = "wet"

	baseMessage = "%s Challenge Builder v0.0.1"

	gameIndex      = 1
	componentIndex = 2

	locationID               = "location"
	distanceID               = "distance"
	stageID                  = "stage"
	weatherID                = "weather"
	SubmitLocationAndStageID = "submit1"

	drivetrainID = "drivetrain"
	classID      = "class"
	carID        = "car"
	SubmitCarID  = "submit2"
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
				buildDistanceMenu(config),
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
					CustomID: strings.Join([]string{ChallengeID, config.Game.ID(), SubmitLocationAndStageID}, idFieldDelimiter),
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
					CustomID: strings.Join([]string{ChallengeID, config.Game.ID(), SubmitCarID}, idFieldDelimiter),
				},
			},
		},
	}
}

func randomOption(category string) discordgo.SelectMenuOption {
	return discordgo.SelectMenuOption{
		Label: "Random " + category, Value: RandomID, Emoji: &discordgo.ComponentEmoji{Name: render.RandomEmoji},
	}
}

// menuEntry is one selectable value in a builder menu, before it is turned into
// a discordgo option.
type menuEntry struct {
	label       string
	value       string
	emoji       string
	description string
}

// buildSelectMenu assembles a string select menu with a leading "Random
// <category>" option. When no entry matches selectedValue the random option
// becomes the default, and a menu offering nothing but that option is disabled.
func buildSelectMenu(config challenge.Config, component, category, selectedValue string, entries []menuEntry) discordgo.SelectMenu {
	options := []discordgo.SelectMenuOption{randomOption(category)}

	matched := false
	for _, entry := range entries {
		option := discordgo.SelectMenuOption{
			Label:   entry.label,
			Value:   entry.value,
			Default: entry.value == selectedValue,
		}
		if entry.emoji != "" {
			option.Emoji = &discordgo.ComponentEmoji{Name: entry.emoji}
		}
		if entry.description != "" {
			option.Description = entry.description
		}
		if option.Default {
			matched = true
		}
		options = append(options, option)
	}

	if !matched {
		options[0].Default = true
	}

	return discordgo.SelectMenu{
		Placeholder: category,
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    buildComponentID(config, component),
		Options:     options,
		Disabled:    len(options) == 1,
	}
}

func buildComponentID(config challenge.Config, component string) string {
	return strings.Join([]string{ChallengeID, config.Game.ID(), component}, idFieldDelimiter)
}

func buildLocationsMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Location != nil {
		selected = strings.ToLower(config.Location.String())
	}

	entries := make([]menuEntry, 0, len(location.List(config.Game)))
	for _, loc := range location.List(config.Game) {
		entries = append(entries, menuEntry{
			label:       loc.String(),
			value:       strings.ToLower(loc.String()),
			emoji:       render.Flag(loc),
			description: loc.DetailedString(),
		})
	}

	return buildSelectMenu(config, locationID, "Location", selected, entries)
}

func buildDistanceMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Distance != nil {
		selected = strings.ToLower(config.Distance.String())
	}

	distances := []stage.Distance{stage.Short, stage.Long}
	if config.Game == game.WRC {
		distances = append(distances, stage.ReallyLong)
	}

	entries := make([]menuEntry, 0, len(distances))
	for _, distance := range distances {
		entries = append(entries, menuEntry{
			label: distance.String(),
			value: strings.ToLower(distance.String()),
			emoji: render.DistanceEmoji(distance),
		})
	}

	return buildSelectMenu(config, distanceID, "Distance", selected, entries)
}

func buildStageMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Stage != nil {
		selected = strings.ToLower(config.Stage.Name())
	}

	var entries []menuEntry
	if config.Location != nil {
		for _, s := range stage.AtLocation(*config.Location) {
			if config.Distance != nil && *config.Distance != s.Distance() {
				continue
			}
			entries = append(entries, menuEntry{
				label:       s.Name(),
				value:       strings.ToLower(s.Name()),
				emoji:       render.DistanceEmoji(s.Distance()),
				description: s.String(),
			})
		}
	}

	return buildSelectMenu(config, stageID, "Stage", selected, entries)
}

func buildWeatherMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Weather != nil {
		selected = strings.ToLower(config.Weather.String())
	}

	weathers := []weather.Model{weather.DRY, weather.WET}
	if config.Location != nil {
		weathers = config.Location.Weather()
	}

	// A location with a single possible weather leaves nothing to choose.
	if config.Location != nil && len(weathers) == 1 {
		only := weathers[0]
		return discordgo.SelectMenu{
			Placeholder: "Weather",
			MenuType:    discordgo.StringSelectMenu,
			CustomID:    buildComponentID(config, weatherID),
			Options: []discordgo.SelectMenuOption{{
				Label:   only.String(),
				Value:   strings.ToLower(only.String()),
				Emoji:   &discordgo.ComponentEmoji{Name: render.WeatherEmoji(only)},
				Default: true,
			}},
			Disabled: true,
		}
	}

	entries := make([]menuEntry, 0, len(weathers))
	for _, w := range weathers {
		entries = append(entries, menuEntry{
			label: w.String(),
			value: strings.ToLower(w.String()),
			emoji: render.WeatherEmoji(w),
		})
	}

	return buildSelectMenu(config, weatherID, "Weather", selected, entries)
}

func buildDriveTrainMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Drivetrain != nil {
		selected = strings.ToLower(config.Drivetrain.String())
	}

	entries := make([]menuEntry, 0, len(drivetrain.List(config.Game)))
	for _, dt := range drivetrain.List(config.Game) {
		entries = append(entries, menuEntry{
			label: dt.String(),
			value: strings.ToLower(dt.String()),
			emoji: render.DrivetrainEmoji(dt),
		})
	}

	return buildSelectMenu(config, drivetrainID, "Drivetrain", selected, entries)
}

func buildClassMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Class != nil {
		selected = strings.ToLower(config.Class.String())
	}

	classes := class.List(config.Game)
	if config.Drivetrain != nil {
		classes = class.WithDrivetrain(*config.Drivetrain, config.Game)
	}

	entries := make([]menuEntry, 0, len(classes))
	for _, c := range classes {
		entries = append(entries, menuEntry{
			label: c.String(),
			value: strings.ToLower(c.String()),
		})
	}

	return buildSelectMenu(config, classID, "Class", selected, entries)
}

func buildCarMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Car != nil {
		selected = strings.ToLower(config.Car.Name())
	}

	var entries []menuEntry
	if config.Class != nil {
		for _, c := range car.InClass(*config.Class, config.Game) {
			entries = append(entries, menuEntry{
				label: c.Name(),
				value: strings.ToLower(c.Name()),
			})
		}
	}

	return buildSelectMenu(config, carID, "Car", selected, entries)
}

func configFromInteraction(sessions SessionStore, interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	fields := strings.Split(customID, idFieldDelimiter)
	if len(fields) != 3 {
		return challenge.Config{}, fmt.Errorf("unexpected customID %s", customID)
	}

	whichGame := game.FromID(fields[gameIndex])
	if whichGame == game.NotSet {
		return challenge.Config{}, fmt.Errorf("invalid game from customID %s", customID)
	}

	builderID := interaction.Message.ID

	config, ok := sessions.Get(builderID)
	if !ok {
		config = challenge.Config{Game: whichGame}
	}

	switch changed := fields[componentIndex]; changed {
	case SubmitLocationAndStageID, SubmitCarID:
		// submit buttons carry no value and only advance the flow
	default:
		values := interaction.MessageComponentData().Values
		if len(values) == 0 {
			return challenge.Config{}, fmt.Errorf("no value for component %s", changed)
		}
		config = applyComponent(config, changed, values[0])
	}

	sessions.Put(builderID, config)

	return config, nil
}

func applyComponent(config challenge.Config, component, value string) challenge.Config {
	switch component {
	case locationID:
		return applyLocation(config, value)
	case distanceID:
		return applyDistance(config, value)
	case stageID:
		return applyStage(config, value)
	case weatherID:
		return applyWeather(config, value)
	case drivetrainID:
		return applyDrivetrain(config, value)
	case classID:
		return applyClass(config, value)
	case carID:
		return applyCar(config, value)
	}
	return config
}

func applyLocation(config challenge.Config, value string) challenge.Config {
	// Location scopes both stage and weather, so any change invalidates them.
	config.Stage = nil
	config.Weather = nil

	if value == RandomID {
		config.Location = nil
		return config
	}

	for _, loc := range location.List(config.Game) {
		if value == strings.ToLower(loc.String()) {
			config.Location = &loc
			return config
		}
	}

	config.Location = nil

	return config
}

func applyDistance(config challenge.Config, value string) challenge.Config {
	switch value {
	case strings.ToLower(stage.Short.String()):
		distance := stage.Short
		config.Distance = &distance
	case strings.ToLower(stage.Long.String()):
		distance := stage.Long
		config.Distance = &distance
	case strings.ToLower(stage.ReallyLong.String()):
		distance := stage.ReallyLong
		config.Distance = &distance
	default:
		config.Distance = nil
	}

	// A fixed distance filters the stage list, so drop a stage that no longer fits.
	if config.Stage != nil && config.Distance != nil && config.Stage.Distance() != *config.Distance {
		config.Stage = nil
	}

	return config
}

func applyStage(config challenge.Config, value string) challenge.Config {
	if value == RandomID || config.Location == nil {
		config.Stage = nil
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
	if value == RandomID {
		config.Weather = nil
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
	// Drivetrain scopes class, which scopes car.
	config.Class = nil
	config.Car = nil

	if value == RandomID {
		config.Drivetrain = nil
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
	// Class scopes car.
	config.Car = nil

	if value == RandomID {
		config.Class = nil
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
		config.Car = nil
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
