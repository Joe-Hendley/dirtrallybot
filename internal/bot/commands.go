package bot

import (
	"log/slog"
	"os"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        challenge.NewstageDR2DefaultID,
			Description: "Generate a new Dirt Rally 2 challenge with default settings",
		},
		{
			Name:        challenge.NewstageDR2CustomID,
			Description: "Generate a new Dirt Rally 2 challenge with custom settings",
			Options:     options,
		},
		{
			Name:        challenge.NewstageWRCDefaultID,
			Description: "Generate a new WRC challenge with default settings",
		},
		{
			Name:        challenge.NewstageWRCCustomID,
			Description: "Generate a new WRC challenge with custom settings",
		},
	}
	cmdIDs map[string]string

	options = []*discordgo.ApplicationCommandOption{
		{
			Name:        "stage",
			Description: "Stage",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "🎲 Random",
					Value: "random",
				},
				{
					Name:  "🌍 Specific Location",
					Value: "location",
				},
				{
					Name:  "🏞️ Specific Stage",
					Value: "stage",
				},
			},
		},
		{
			Name:        "weather",
			Description: "Weather",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "🎲 Random",
					Value: "random",
				},
				{
					Name:  "☀️ Dry",
					Value: "random",
				},
				{
					Name:  "🌧️ Wet",
					Value: "stage",
				},
			},
		},
		{
			Name:        "car",
			Description: "Car",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "🎲 Fully Random",
					Value: "random",
				},
				{
					Name:  "⚙️ Specific Drivetrain",
					Value: "drivetrain",
				},
				{
					Name:  "📋 Specific Class",
					Value: "class",
				},
				{
					Name:  "🏎️ Specific Car",
					Value: "class",
				},
			},
		},
	}
)

func CreateCommands(config config.Config, session *discordgo.Session) {

	cmdIDs = make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := session.ApplicationCommandCreate(config.App, "", cmd)
		if err != nil {
			slog.Error("creating slash command", "cmd", cmd.Name, "err", err)
			os.Exit(1)
		}

		slog.Debug("registered command", "cmd", rcmd.Name)

		cmdIDs[rcmd.ID] = rcmd.Name
	}
}

func CleanupCommands(config config.Config, session *discordgo.Session) {
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		slog.Error("fetching registered slash commands", "err", err)
		os.Exit(1)
	}

	for _, cmd := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
		if err != nil {
			slog.Error("deleting slash command", "cmd", cmd.Name, "err", err)
			os.Exit(1)
		}
	}
}
