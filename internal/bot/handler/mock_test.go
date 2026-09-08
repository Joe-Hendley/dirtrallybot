package handler_test

import (
	"context"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
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

var _ port.Store = &storeMock{}

type storeMock struct {
	mock.Mock
}

// DeleteChallenge implements port.Store.
func (sm *storeMock) DeleteChallenge(_ context.Context, challengeID string) error {
	args := sm.Called(challengeID)
	return args.Error(0)
}

// GetChallenge implements port.Store.
func (sm *storeMock) GetChallenge(_ context.Context, challengeID string) (challenge.Model, error) {
	args := sm.Called(challengeID)
	return args.Get(0).(challenge.Model), args.Error(1)
}

// PutChallenge implements port.Store.
func (sm *storeMock) PutChallenge(_ context.Context, challengeID string, challenge challenge.Model) error {
	args := sm.Called(challengeID)
	return args.Error(0)
}

// RegisterCompletion implements port.Store.
func (sm *storeMock) RegisterCompletion(_ context.Context, challengeID string, completion challenge.Completion) error {
	args := sm.Called(challengeID, completion)
	return args.Error(0)
}

// RegisterVote implements port.Store.
func (sm *storeMock) RegisterVote(_ context.Context, challengeID string, vote challenge.Vote) error {
	args := sm.Called(challengeID, vote)
	return args.Error(0)
}

// Popularity implements port.Store.
func (sm *storeMock) Popularity(_ context.Context) (popularity.Snapshot, error) {
	args := sm.Called()
	return args.Get(0).(popularity.Snapshot), args.Error(1)
}
