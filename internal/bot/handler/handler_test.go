package handler_test

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/completion"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	challengeModel "github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
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

			expectedUserID      = "someUserID"
			expectedChallengeID = "someMessageID"
			expectedResponseID  = strings.Join([]string{completion.SubmitCompletionPrefix, expectedChallengeID, expectedUserID}, "_")
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
					ID: expectedChallengeID,
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

	// When we recieve an invalid modal submit interaction,
	// with a customID corresponding to "submit completion",
	// we respond to inform the user the time is invalid
	// and do nothing else
	t.Run("TestDisplayCompletionEntryModal_InvalidTimestamp", func(t *testing.T) {
		// Arrange
		var (
			expectedOptions []discordgo.RequestOption

			userID             = "someUserID"
			challengeID        = "someMessageID"
			submissionCustomID = strings.Join([]string{completion.SubmitCompletionPrefix, challengeID, userID}, "_")

			expectedCustomID = completion.InvalidSubmissionID
		)
		interaction := discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionModalSubmit,
				Data: discordgo.ModalSubmitInteractionData{
					CustomID: submissionCustomID,
					Components: []discordgo.MessageComponent{
						&discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.TextInput{
									CustomID: completion.CompletionTextInputID,
									Value:    "invalid",
								},
							},
						},
					},
				},
			},
		}

		store := new(storeMock)
		session := new(sessionMock)
		session.On("InteractionRespond", interaction.Interaction, mock.AnythingOfType("*discordgo.InteractionResponse"), expectedOptions).Return(nil)

		// Act
		handler.ModalSubmit(store, session, &interaction)

		// Assert
		store.AssertExpectations(t)
		session.AssertExpectations(t)

		if assert.NotNil(t, session.Calls[0].Arguments[1]) {
			response := session.Calls[0].Arguments[1].(*discordgo.InteractionResponse)

			assert.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, response.Type)
			if assert.NotNil(t, response.Data) {
				assert.Equal(t, expectedCustomID, response.Data.CustomID)
			}
		}
	})

	// When we recieve a valid modal submit interaction,
	// with a customID corresponding to "submit completion",
	// we acknowledge the submission
	// and store the completion against the corresponding challenge
	// and update the challenge message if the top three has changed
	t.Run("TestDisplayCompletionEntryModal_ValidTimestamp", func(t *testing.T) {
		// Arrange
		var (
			expectedOptions []discordgo.RequestOption

			userID             = "someUserID"
			username           = "username"
			challengeID        = "someMessageID"
			guildID            = "someGuildID"
			submissionCustomID = strings.Join([]string{completion.SubmitCompletionPrefix, challengeID, userID}, "_")

			validTimestamp  = "1:23.450"
			parsedTimestamp = timestamp.Build(1, 23, 450)

			expectedResponseCustomID = completion.ValidSubmissionID
			expectedCompletion       = challengeModel.NewCompletion(userID, parsedTimestamp)

			storedChallenge = challengeModel.NewChallenge(stage.Model{}, weather.DRY, car.Model{}, []challengeModel.Completion{expectedCompletion})
		)
		interaction := discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionModalSubmit,
				Data: discordgo.ModalSubmitInteractionData{
					CustomID: submissionCustomID,
					Components: []discordgo.MessageComponent{
						&discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.TextInput{
									CustomID: completion.CompletionTextInputID,
									Value:    validTimestamp,
								},
							},
						},
					},
				},
				GuildID: guildID,
			},
		}

		store := new(storeMock)
		store.On("RegisterCompletion", challengeID, expectedCompletion).Return(nil)
		store.On("GetChallenge", challengeID).Return(storedChallenge, nil)

		session := new(sessionMock)
		session.On("InteractionRespond", interaction.Interaction, mock.AnythingOfType("*discordgo.InteractionResponse"), expectedOptions).Return(nil)
		session.On("GuildMember", guildID, userID, expectedOptions).Return(&discordgo.Member{User: &discordgo.User{GlobalName: username}}, nil)
		session.On("ChannelMessageEditComplex", mock.AnythingOfType("*discordgo.MessageEdit"), expectedOptions).Return(&discordgo.Message{}, nil)

		// Act
		handler.ModalSubmit(store, session, &interaction)

		// Assert
		store.AssertExpectations(t)
		session.AssertExpectations(t)

		if assert.NotNil(t, session.Calls[0].Arguments[1]) {
			response := session.Calls[0].Arguments[1].(*discordgo.InteractionResponse)

			assert.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, response.Type)
			if assert.NotNil(t, response.Data) {
				assert.Equal(t, expectedResponseCustomID, response.Data.CustomID)
			}
		}

		if assert.NotNil(t, session.Calls[2].Arguments[0]) {
			editedMessage := session.Calls[2].Arguments[0].(*discordgo.MessageEdit)
			assert.Contains(t, *editedMessage.Content, timestamp.Format(parsedTimestamp))
			assert.Contains(t, *editedMessage.Content, username)
		}
	})
}
