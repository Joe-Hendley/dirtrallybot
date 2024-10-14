package handler_test

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/completion"
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

// When we recieve an application interaction,
// with a customID corresponding to "new challenge",
// we respond with the correct message.
func TestNewChallengeRouting(t *testing.T) {
	testData := []routingTestCase{
		{title: "TestHandleNewDR2Challenge", incomingID: challenge.NewDR2ChallengeID, outgoingID: challenge.InitialDR2ChallengeResponseID},
		{title: "TestHandleNewWRCChallenge", incomingID: challenge.NewWRCChallengeID, outgoingID: challenge.InitialWRCChallengeResponseID},
	}

	for _, testCase := range testData {
		t.Run(testCase.title, func(t *testing.T) {
			// Arrange
			interaction := createApplicationCommandInteraction(testCase.incomingID)
			var expectedOptions []discordgo.RequestOption

			session := new(sessionMock)
			session.On("InteractionRespond", interaction.Interaction, mock.AnythingOfType("*discordgo.InteractionResponse"), expectedOptions).Return(nil)

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
func TestSubmitCompletion(t *testing.T) {
	// When we recieve a message component interaction,
	// with a customID corresponding to "display completion entry modal",
	// we respond with the completion submission modal
	t.Run("TestDisplayCompletionEntryModal", func(t *testing.T) {
		// Arrange
		var (
			expectedOptions []discordgo.RequestOption

			expectedUserID     = "someUserID"
			expectedMessageID  = "someMessageID"
			expectedResponseID = strings.Join([]string{completion.SubmitCompletionPrefix, expectedMessageID, expectedUserID}, "-")
		)
		interaction := discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionMessageComponent,
				Data: discordgo.MessageComponentInteractionData{
					CustomID: challenge.CompletedID,
				},
				Member: &discordgo.Member{
					User: &discordgo.User{
						ID: expectedUserID,
					},
				},
				Message: &discordgo.Message{
					ID: expectedMessageID,
				},
			},
		}

		store := new(storeMock)
		session := new(sessionMock)
		session.On("InteractionRespond", interaction.Interaction, mock.AnythingOfType("*discordgo.InteractionResponse"), expectedOptions).Return(nil)

		// Act
		handler.InteractionMessageComponent(store, session, &interaction)

		// Assert
		store.AssertExpectations(t)
		session.AssertExpectations(t)

		if assert.NotNil(t, session.Calls[0].Arguments[1]) {
			response := session.Calls[0].Arguments[1].(*discordgo.InteractionResponse)

			assert.Equal(t, discordgo.InteractionResponseModal, response.Type)
			if assert.NotNil(t, response.Data) {
				assert.Equal(t, expectedResponseID, response.Data.CustomID)
			}
		}
	})
	// When we recieve a modal submit interaction,
	// with a customID corresponding to "submit completion",
	// we store the completion against the corresponding challenge
	// and update the challenge message if the top three has changed
}
