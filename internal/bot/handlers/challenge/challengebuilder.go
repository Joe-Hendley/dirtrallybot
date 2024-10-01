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

	LocationSelectID = "challenge-location"
	StageSelectID    = "challenge-stage"
	WeatherSelectID  = "challenge-weather"

	DrivetrainSelectID = "challenge-drivetrain"
	ClassSelectID      = "challenge-class"
	CarSelectID        = "challenge-car"

	SubmitID = "challenge-submit"
)

func buildChallengeMessageComponents(config challenge.Config) []discordgo.MessageComponent {
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
				},
			},
		},
	}
}

var RandomOption = discordgo.SelectMenuOption{
	Label: "Random", Value: RandomID, Emoji: &discordgo.ComponentEmoji{Name: "🎲"},
}

func buildLocationsMenu(config challenge.Config) discordgo.SelectMenu {
	var selected string
	if config.Location != nil {
		selected = strings.ToLower(config.Location.String())
	}

	options := []discordgo.SelectMenuOption{
		RandomOption,
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
		RandomOption,
	}

	hasDefault := false

	if config.Location != nil {
		for _, stage := range stage.AtLocation(*config.Location) {
			stageID := strings.ToLower(stage.String())
			options = append(options, discordgo.SelectMenuOption{
				Label:       stage.String(),
				Value:       stageID,
				Emoji:       &discordgo.ComponentEmoji{},
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
			RandomOption,
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
			RandomOption,
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
		RandomOption,
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
		RandomOption,
	}

	hasDefault := false

	if config.Drivetrain != nil {
		for _, class := range class.WithDrivetrain(*config.Drivetrain) {
			classID := strings.ToLower(config.Class.String())
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
		RandomOption,
	}

	hasDefault := false

	if config.Class != nil {
		for _, car := range car.InClass(*config.Class) {
			carID := strings.ToLower(config.Class.String())
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

func buildConfigFromInteraction(interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	customID := interaction.MessageComponentData().CustomID
	var newValue string
	if customID != SubmitID {
		newValue = interaction.MessageComponentData().Values[0]
	}
	config := challenge.Config{}

	// each component is contained in a separate action row
	// should be some number of select menus then a button
	for _, componentRow := range interaction.Message.Components {
		actionsRow, ok := componentRow.(discordgo.ActionsRow)
		if !ok {
			return challenge.Config{}, errors.New("expected action row")
		}

		switch component := actionsRow.Components[0].(type) {
		case discordgo.Button:
			continue
		case discordgo.SelectMenu:
			currentValue := ""
			for _, option := range component.Options {
				if option.Default {
					currentValue = option.Value
					break
				}
			}
			config = applyValueToConfig(config, component.CustomID, currentValue)
		default:
			return challenge.Config{}, errors.New("unexpected component type")
		}
	}

	if newValue != "" {

	}

	return config, nil
}

func applyValueToConfig(config challenge.Config, fieldID, value string) challenge.Config {
	return config
}
