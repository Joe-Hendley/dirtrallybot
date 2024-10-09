package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type Session interface {
	InteractionResponder
	ChannelMessageSender
	ChannelMessageEditor
	GuildMemberGetter
}

type InteractionResponder interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
}

type ChannelMessageSender interface {
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
}

type ChannelMessageEditor interface {
	ChannelMessageEditComplex(m *discordgo.MessageEdit, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
}

type GuildMemberGetter interface {
	GuildMember(guildID, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
}

func GetGuildMemberDisplayName(session GuildMemberGetter, guildID, userID string) string {
	u, err := session.GuildMember(guildID, userID)
	if err != nil {
		slog.Error("getting display name", "guildID", guildID, "userID", userID, "err", err)
		return userID
	}

	return u.DisplayName()
}
