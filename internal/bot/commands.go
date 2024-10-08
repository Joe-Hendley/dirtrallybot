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
			Name:        challenge.NewstageDR2CustomID,
			Description: "Generate a new Dirt Rally 2 challenge with custom settings",
		},
		{
			Name:        challenge.NewstageWRCCustomID,
			Description: "Generate a new WRC challenge with custom settings",
		},
	}
	cmdIDs map[string]string
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

func CleanupGuildCommands(config config.Config, session *discordgo.Session) {
	for _, guild := range session.State.Guilds {
		guildID := guild.ID
		registeredCommands, err := session.ApplicationCommands(session.State.User.ID, guildID)
		if err != nil {
			slog.Error("fetching registered slash commands", "err", err)
			os.Exit(1)
		}

		for _, cmd := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, guildID, cmd.ID)
			if err != nil {
				slog.Error("deleting slash command", "cmd", cmd.Name, "err", err)
			}
		}
	}
}

func CleanupGlobalCommands(config config.Config, session *discordgo.Session) {
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		slog.Error("fetching registered slash commands", "err", err)
		os.Exit(1)
	}

	for _, cmd := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
		if err != nil {
			slog.Error("deleting slash command", "cmd", cmd.Name, "err", err)
		}
	}
}
