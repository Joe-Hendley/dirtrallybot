package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/debug"
	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/buildersession"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	ctx      context.Context
	cfg      config.Config
	session  *discordgo.Session
	store    port.Store
	sessions *buildersession.Store
}

// New wires up the bot and registers its slash commands. ctx bounds the lifetime
// of work started from interaction handlers; cancelling it unwinds in-flight
// store operations at shutdown.
func New(ctx context.Context, cfg config.Config, store port.Store, session *discordgo.Session) (*Bot, error) {
	bot := &Bot{
		ctx:      ctx,
		cfg:      cfg,
		session:  session,
		store:    store,
		sessions: buildersession.New(),
	}

	session.AddHandler(bot.HandleReady)
	session.AddHandler(bot.HandleMessageCreate)
	session.AddHandler(bot.HandleInteractionCreate)

	if err := createCommands(cfg, session); err != nil {
		return nil, fmt.Errorf("registering commands: %w", err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	return bot, nil
}

func (bot *Bot) Shutdown() error {
	return errors.Join(
		cleanupGuildCommands(bot.session),
		cleanupGlobalCommands(bot.session),
	)
}

func (bot *Bot) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is ready")
}

func (bot *Bot) HandleMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	slog.Debug("message", "id", message.ID, "author_id", message.Author.ID, "author_name", message.Author.GlobalName)

	lowercase := strings.ToLower(message.Content)

	switch lowercase {
	case "!cars":
		if message.GuildID != bot.cfg.TestServerID {
			return
		}
		debug.HandleCars(session, message)

	case "!stages":
		if message.GuildID != bot.cfg.TestServerID {
			return
		}
		debug.HandleStages(session, message)

	default:
		slog.Debug("message ignored", "id", message.ID)
	}
}

func (bot *Bot) HandleInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		handler.ApplicationCommand(session, interaction)
	case discordgo.InteractionMessageComponent:
		handler.InteractionMessageComponent(bot.ctx, bot.sessions, bot.store, session, interaction)
	case discordgo.InteractionModalSubmit:
		handler.ModalSubmit(bot.ctx, bot.store, session, interaction)
	}
}
