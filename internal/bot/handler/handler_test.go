package handler_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createApplicationCommandInteraction(name string) discordgo.InteractionCreate {
	return discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: name,
			},
		},
	}
}

type routingTestCase struct {
	title, incomingID, outgoingID string
}

func TestNewChallengeRouting(t *testing.T) {
	testData := []routingTestCase{
		{title: "TestHandleNewDR2Challenge", incomingID: challenge.NewDR2ChallengeID, outgoingID: challenge.InitialDR2ChallengeResponseID},
		{title: "TestHandleNewWRCChallenge", incomingID: challenge.NewWRCChallengeID, outgoingID: challenge.InitialWRCChallengeResponseID},
	}

	for _, testCase := range testData {
		t.Run(testCase.title, func(t *testing.T) {
			// Arrange
			interaction := createApplicationCommandInteraction(testCase.incomingID)

			session := new(sessionMock)
			session.On("InteractionRespond", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			// Act
			handler.ApplicationCommand(session, &interaction)

			// Assert
			session.AssertExpectations(t)

			if assert.NotNil(t, session.Calls[0].Arguments[1]) {
				response := session.Calls[0].Arguments[1].(*discordgo.InteractionResponse)

				assert.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, response.Type)
				if assert.NotNil(t, response.Data) {
					assert.Equal(t, testCase.outgoingID, response.Data.CustomID)
				}
			}
		})
	}
}

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
