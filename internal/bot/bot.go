package bot

import (
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/debug"
	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	cfg     config.Config
	session *discordgo.Session
	store   model.Store
}

func New(cfg config.Config, store model.Store, session *discordgo.Session) (*bot, error) {

	bot := &bot{
		session: session,
		store:   store,
	}

	session.AddHandler(bot.HandleReady)
	session.AddHandler(bot.HandleMessageCreate)
	session.AddHandler(bot.HandleInteractionCreate)

	CreateCommands(cfg, session)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	return bot, nil
}

func (bot *bot) Shutdown() {
	CleanupCommands(bot.cfg, bot.session)
}

func (bot *bot) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is ready")
}

func (bot *bot) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	slog.Debug("message", "id", m.ID, "author_id", m.Author.ID, "author_name", m.Author.GlobalName)

	lowercase := strings.ToLower(m.Content)

	switch lowercase {
	case "!cars":
		debug.HandleCars(s, m)

	case "!stages":
		debug.HandleStages(s, m)

	case "!newstage":
		challenge.HandleNewChallenge(bot.store, s, m)

	default:
		slog.Debug("message ignored", "id", m.ID)
	}
}

func (bot *bot) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		handlers.InteractionMessageComponent(bot.store, s, i)
	case discordgo.InteractionModalSubmit:
		handlers.ModalSubmit(bot.store, s, i)
	}
}
