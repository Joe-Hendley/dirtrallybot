package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/debug"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/web"
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
	web      *web.Server
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

	if cfg.WebAddr != "" {
		bot.web = web.New(cfg.WebAddr, store)
	}

	challenge.SetGenerator(generatorFor(cfg.Randomiser))

	session.AddHandler(bot.HandleReady)
	session.AddHandler(bot.HandleMessageCreate)
	session.AddHandler(bot.HandleInteractionCreate)

	if err := createCommands(cfg, session); err != nil {
		return nil, fmt.Errorf("registering commands: %w", err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	return bot, nil
}

// generatorFor maps the configured randomiser choice to its generator, defaulting
// to the popularity-biased one.
func generatorFor(kind config.RandomiserType) challenge.Generator {
	switch kind {
	case config.RandomiserDeterministic:
		return challenge.DeterministicGenerator
	case config.RandomiserRandom:
		return challenge.RandomGenerator
	default:
		return challenge.BiasedGenerator
	}
}

// ServeWeb runs the local challenge viewer until Shutdown is called. It returns
// immediately when the viewer is disabled.
func (bot *Bot) ServeWeb() error {
	if bot.web == nil {
		return nil
	}
	return bot.web.Start()
}

func (bot *Bot) Shutdown() error {
	errs := []error{
		cleanupGuildCommands(bot.session),
		cleanupGlobalCommands(bot.session),
	}

	if bot.web != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		errs = append(errs, bot.web.Shutdown(ctx))
	}

	return errors.Join(errs...)
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
