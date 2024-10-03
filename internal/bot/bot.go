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
	CleanupGuildCommands(bot.cfg, bot.session)
	CleanupGlobalCommands(bot.cfg, bot.session)
}

func (bot *bot) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is ready")
}

func (bot *bot) HandleMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	slog.Debug("message", "id", message.ID, "author_id", message.Author.ID, "author_name", message.Author.GlobalName)

	lowercase := strings.ToLower(message.Content)

	switch lowercase {
	case "!cars":
		debug.HandleCars(session, message)

	case "!stages":
		debug.HandleStages(session, message)

	case "!newstage":
		challenge.HandleCreateDR2ChallengeDefault(bot.store, session, challenge.NewInvocationFromMessageCreate(*message))

	default:
		slog.Debug("message ignored", "id", message.ID)
	}
}

func (bot *bot) HandleInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		handlers.ApplicationCommand(bot.store, session, interaction)
	case discordgo.InteractionMessageComponent:
		handlers.InteractionMessageComponent(bot.store, session, interaction)
	case discordgo.InteractionModalSubmit:
		handlers.ModalSubmit(bot.store, session, interaction)
	}
}
