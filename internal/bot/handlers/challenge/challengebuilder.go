package challenge

import (
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/bwmarrin/discordgo"
)

func buildChallengeConfig(session *discordgo.Session, interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Custom DR2 hasn't been implemented yet! Sorry!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge", "err", err)
		return challenge.Config{}, err
	}

	return challenge.Config{}, nil
}

// so, we get the interaction in with the following options:
// location - fully random, specific location, specific stage
// weather 	- fully random, specific weather
// car		- fully random, specific drivetrain, specific class, specific car

// we then need to populate the select options with the list of options for these

const (
	RandomID     = "random"
	LocationID   = "location"
	StageID      = "stage"
	DryID        = "dry"
	WetID        = "wet"
	DrivetrainID = "drivetrain"
	ClassID      = "class"
	CarID        = "car"
)

func buildChallengeMessageComponents(config challenge.Config) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{},
		},
	}
}

var RandomOption = discordgo.SelectMenuOption{
	Label: "Random", Value: RandomID, Emoji: &discordgo.ComponentEmoji{Name: "🎲"},
}

func buildLocationsMenu() discordgo.SelectMenu {
	options := []discordgo.SelectMenuOption{
		RandomOption,
	}

	for _, loc := range location.List() {
		options = append(options, discordgo.SelectMenuOption{
			Label: loc.String(), Value: strings.ToLower(loc.String()), Emoji: &discordgo.ComponentEmoji{Name: loc.Flag()}, Description: loc.LongString(),
		})
	}

	return discordgo.SelectMenu{
		Placeholder: "Location",
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    LocationID,
		Options:     options,
	}
}

func buildStageMenu(challenge.Config) discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{
		RandomOption,
	}

	return
}
