package handler_test

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

var _ discord.Session = &sessionMock{}

type sessionMock struct {
	mock.Mock
}

// ChannelMessageEditComplex implements discord.Session.
func (sm *sessionMock) ChannelMessageEditComplex(m *discordgo.MessageEdit, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	args := sm.Called(m, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// ChannelMessageSendComplex implements discord.Session.
func (sm *sessionMock) ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	args := sm.Called(channelID, data, options)
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// GuildMember implements discord.Session.
func (sm *sessionMock) GuildMember(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error) {
	args := sm.Called(guildID, userID, options)
	return args.Get(0).(*discordgo.Member), args.Error(1)
}

// InteractionRespond implements discord.Session.
func (sm *sessionMock) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
	args := sm.Called(interaction, resp, options)
	return args.Error(0)
}

var _ model.Store = &storeMock{}

type storeMock struct {
	mock.Mock
}

// DeleteChallenge implements model.Store.
func (sm *storeMock) DeleteChallenge(challengeID string) error {
	args := sm.Called(challengeID)
	return args.Error(0)
}

// GetChallenge implements model.Store.
func (sm *storeMock) GetChallenge(challengeID string) (challenge.Model, error) {
	args := sm.Called(challengeID)
	return args.Get(0).(challenge.Model), args.Error(1)
}

// PutChallenge implements model.Store.
func (sm *storeMock) PutChallenge(challengeID string, challenge challenge.Model) error {
	args := sm.Called(challengeID)
	return args.Error(0)
}

// RegisterCompletion implements model.Store.
func (sm *storeMock) RegisterCompletion(challengeID string, completion challenge.Completion) error {
	args := sm.Called(challengeID, completion)
	return args.Error(0)
}
