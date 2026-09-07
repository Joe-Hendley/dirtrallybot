package bot

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        challenge.NewDR2ChallengeID,
		Description: "Generate a new Dirt Rally 2 challenge with custom settings",
	},
	{
		Name:        challenge.NewWRCChallengeID,
		Description: "Generate a new WRC challenge with custom settings",
	},
}

func createCommands(cfg config.Config, session *discordgo.Session) error {
	for _, cmd := range commands {
		registered, err := session.ApplicationCommandCreate(cfg.App, "", cmd)
		if err != nil {
			return fmt.Errorf("creating slash command %s: %w", cmd.Name, err)
		}

		slog.Debug("registered command", "cmd", registered.Name)
	}

	return nil
}

func cleanupGuildCommands(session *discordgo.Session) error {
	var errs []error

	for _, guild := range session.State.Guilds {
		registered, err := session.ApplicationCommands(session.State.User.ID, guild.ID)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetching commands for guild %s: %w", guild.ID, err))
			continue
		}

		for _, cmd := range registered {
			if err := session.ApplicationCommandDelete(session.State.User.ID, guild.ID, cmd.ID); err != nil {
				errs = append(errs, fmt.Errorf("deleting guild command %s: %w", cmd.Name, err))
			}
		}
	}

	return errors.Join(errs...)
}

func cleanupGlobalCommands(session *discordgo.Session) error {
	registered, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("fetching global commands: %w", err)
	}

	var errs []error
	for _, cmd := range registered {
		if err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID); err != nil {
			errs = append(errs, fmt.Errorf("deleting global command %s: %w", cmd.Name, err))
		}
	}

	return errors.Join(errs...)
}
