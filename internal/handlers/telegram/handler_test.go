package telegram

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/kyokomi/emoji"
	"github.com/ritlug/teleirc/internal"
	"github.com/stretchr/testify/assert"
)

/*
TestPartFullOn tests the ability of the partHandler to send messages
when ShowLeaveMessage is set to true
*/
func TestPartFullOn(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		FirstName: "test",
		Username:  "testUser",
	}
	correct := "test (@testUser) has left the Telegram Group!"

	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:           "<",
			Suffix:           ">",
			ShowLeaveMessage: true,
			ShowZWSP:         false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	partHandler(clientObj, testUser)
}

/*
TestPartFullZwsp tests the full capacity of the Part handler with zero-width spaces
*/
func TestPartFullZwsp(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		FirstName: "test",
		Username:  "testUser",
	}
	correct := "test (@t" + "\u200b" + "estUser) has left the Telegram Group!"

	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:           "<",
			Suffix:           ">",
			ShowLeaveMessage: true,
			ShowZWSP:         true,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	partHandler(clientObj, testUser)
}

/*
TestPartFullOff tests the ability of the partHandler to not send messages
when ShowLeaveMessage is set to false
*/
func TestPartFullOff(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		FirstName: "test",
		Username:  "testUser",
	}
	correct := ""
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:           "<",
			Suffix:           ">",
			ShowLeaveMessage: false,
			ShowZWSP:         false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	partHandler(clientObj, testUser)
}

/*
TestPartNoUsername tests the ability of the partHandler to send correctly
formatted messages when a TG user has no Username
*/
func TestPartNoUsername(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		FirstName: "test",
	}
	correct := testUser.FirstName + " has left the Telegram Group!"
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:           "<",
			Suffix:           ">",
			ShowLeaveMessage: true,
			ShowZWSP:         false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	partHandler(clientObj, testUser)
}

/*
TestJoinFullOn tests the ability of the joinHandler to send messages
when ShowJoinMessage is set to true
*/
func TestJoinFullOn(t *testing.T) {
	testListUser := &[]models.User{
		{
			ID:        1,
			FirstName: "test",
			Username:  "testUser",
		},
	}
	correct := "test (@testUser) has joined the Telegram Group!"
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:          "<",
			Suffix:          ">",
			ShowJoinMessage: true,
			ShowZWSP:        false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	joinHandler(clientObj, testListUser)
}

/*
TestJoinFullZwsp tests the full capacity of the Join handler with zero-width spaces
*/
func TestJoinFullZwsp(t *testing.T) {
	testListUser := &[]models.User{
		{
			ID:        1,
			FirstName: "test",
			Username:  "testUser",
		},
	}
	correct := "test (@t" + "\u200b" + "estUser) has joined the Telegram Group!"
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:          "<",
			Suffix:          ">",
			ShowJoinMessage: true,
			ShowZWSP:        true,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	joinHandler(clientObj, testListUser)
}

/*
TestJoinFullOff tests the ability of the joinHandler to not send messages
when ShowJoinMessage is set to false
*/
func TestJoinFullOff(t *testing.T) {
	testListUser := &[]models.User{
		{
			ID:        1,
			FirstName: "test",
			Username:  "testUser",
		},
	}
	correct := ""
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:          "<",
			Suffix:          ">",
			ShowJoinMessage: false,
			ShowZWSP:        false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	joinHandler(clientObj, testListUser)
}

/*
TestJoinNoUsername tests the ability of the joinHandler to send correctly
formatted messages when a TG user has no Username
*/
func TestJoinNoUsername(t *testing.T) {
	testListUser := &[]models.User{
		{
			ID:        1,
			FirstName: "test",
		},
	}
	correct := "test has joined the Telegram Group!"
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			Prefix:          "<",
			Suffix:          ">",
			ShowJoinMessage: true,
			ShowZWSP:        false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	joinHandler(clientObj, testListUser)
}

/*
TestDocumentPlain checks the behavior of the document handlers if only required
fields are available.
*/
func TestDocumentPlain(t *testing.T) {
	correct := "test shared a file"
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
			},
			Document: &models.Document{
				FileID: "https://teleirc.com/file.txt",
			},
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestDocumentBasic checks the behavior of the document handlers when
the update just has required informations in addition to the caption.
*/
func TestDocumentBasic(t *testing.T) {
	correct := "test shared a file on Telegram with caption: 'Random Caption'."
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
			},
			Document: &models.Document{
				FileID: "https://teleirc.com/file.txt",
			},
			Caption: "Random Caption",
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestDocumentMime checks the behavior of the document handlers when
the document contains the mimetype information.
*/
func TestDocumentMime(t *testing.T) {
	correct := "test shared a file (test/txt) on Telegram with caption: 'Random Caption'."
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
			},
			Document: &models.Document{
				FileID:   "https://teleirc.com/file.txt",
				MimeType: "test/txt",
			},
			Caption: "Random Caption",
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestDocumentUsername checks the behavior of the document handlers when
both firstname and Username exist. It also incorporates the availability of a mimetype.
*/
func TestDocumentUsername(t *testing.T) {
	correct := "user shared a file (test/txt) on Telegram with caption: 'Random Caption'."
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
				Username:  "user",
			},
			Document: &models.Document{
				FileID:   "https://teleirc.com/file.txt",
				MimeType: "test/txt",
			},
			Caption: "Random Caption",
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestDocumentNoCaption checks the behavior of the document handlers when neither
a caption nor a Username is attached to the document. It also test a case where
both filename and mimetype exist.
*/
func TestDocumentNoCaption(t *testing.T) {
	correct := "test shared a file (test/txt) on Telegram with title: 'test.txt'."
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
			},
			Document: &models.Document{
				FileID:   "https://teleirc.com/file.txt",
				MimeType: "test/txt",
				FileName: "test.txt",
			},
			Caption: "",
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestDocumentFull checks the behavior of the document handlers when
both caption and filename exist. It also incorporates the availability of both
firstname and Username
*/
func TestDocumentFull(t *testing.T) {
	correct := "u" + "\u200b" +
		"ser shared a file (test/txt) on Telegram with caption: 'Random Caption'."
	updateObj := &models.Update{
		Message: &models.Message{
			From: &models.User{
				FirstName: "test",
				Username:  "user",
			},
			Document: &models.Document{
				FileID:   "https://teleirc.com/file.txt",
				MimeType: "test/txt",
				FileName: "test.txt",
			},
			Caption: "Random Caption",
		},
	}
	clientObj := &Client{
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: true,
		},
	}
	documentHandler(clientObj, updateObj.Message)
}

/*
TestPhotoFull tests a complete Photo object
*/
// func TestPhotoFull(t *testing.T) {
// 	t.Skip()
// 	correct := "u" + "\u200b" +
// 		"ser shared a photo on Telegram with caption: 'Random Caption'"
// 	updateObj := models.Update{
// 		Message: &models.Message{
// 			From: &models.User{
// 				FirstName: "test",
// 				Username:  "user",
// 			},
// 			Photo: &[]tgbotapi.PhotoSize{
// 				{
// 					FileID:   "https://teleirc.com/file.png",
// 					Width:    1,
// 					Height:   1,
// 					FileSize: 1,
// 				},
// 			},
// 			Caption: "Random Caption",
// 		},
// 	}
// 	clientObj := &Client{
// 		sendToIrc: func(s string) {
// 			assert.Equal(t, correct, s)
// 		},
// 		IRCSettings: &internal.IRCSettings{
// 			ShowZWSP: true,
// 		},
// 	}
// 	photoHandler(clientObj, updateObj)
// }

/*
TestPhotoNoUsername tests a Photo object with no Username present. Should default
to user's FirstName
*/
// func TestPhotoNoUsername(t *testing.T) {
// 	t.Skip()
// 	correct := "test shared a photo on Telegram with caption: 'Random Caption'"
// 	updateObj := models.Update{
// 		Message: &models.Message{
// 			From: &models.User{
// 				FirstName: "test",
// 			},
// 			Photo: &[]tgbotapi.PhotoSize{
// 				{
// 					FileID:   "https://teleirc.com/file.png",
// 					Width:    1,
// 					Height:   1,
// 					FileSize: 1,
// 				},
// 			},
// 			Caption: "Random Caption",
// 		},
// 	}
// 	clientObj := &Client{
// 		sendToIrc: func(s string) {
// 			assert.Equal(t, correct, s)
// 		},
// 		IRCSettings: &internal.IRCSettings{
// 			ShowZWSP: false,
// 		},
// 	}
// 	photoHandler(clientObj, updateObj)
// }

/*
TestPhotoNoCaption tests messages are correctly formatted when a photo
is uploaded without a caption
*/
// func TestPhotoNoCaption(t *testing.T) {
// 	t.Skip()
// 	correct := "user shared a photo on Telegram with caption: ''"
// 	updateObj := models.Update{
// 		Message: &models.Message{
// 			From: &models.User{
// 				FirstName: "test",
// 				Username:  "user",
// 			},
// 			Photo: &[]tgbotapi.PhotoSize{
// 				{
// 					FileID:   "https://teleirc.com/file.png",
// 					Width:    1,
// 					Height:   1,
// 					FileSize: 1,
// 				},
// 			},
// 			Caption: "",
// 		},
// 	}
// 	clientObj := &Client{
// 		sendToIrc: func(s string) {
// 			assert.Equal(t, correct, s)
// 		},
// 		IRCSettings: &internal.IRCSettings{
// 			ShowZWSP: false,
// 		},
// 	}
// 	photoHandler(clientObj, updateObj)
// }

func TestStickerSmileWithUsername(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	correct := fmt.Sprintf("<%s> 😄", testUser.Username)
	updateObj := models.Update{
		Message: &models.Message{
			From: testUser,
			Sticker: &models.Sticker{
				Emoji: strings.Trim(emoji.Sprint(":smile:"), " "),
			},
		},
	}

	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}

	stickerHandler(clientObj, updateObj)

}

func TestStickerSmileZWSP(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	correct := "<t" + "\u200b" + "est> 😄"
	updateObj := models.Update{
		Message: &models.Message{
			From: testUser,
			Sticker: &models.Sticker{
				Emoji: strings.Trim(emoji.Sprint(":smile:"), " "),
			},
		},
	}

	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: true,
		},
	}

	stickerHandler(clientObj, updateObj)

}

func TestStickerSmileWithoutUsername(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "",
		FirstName: "testing",
		LastName:  "123",
	}
	correct := fmt.Sprintf("<%s> 😄", testUser.FirstName)
	updateObj := models.Update{
		Message: &models.Message{
			From: testUser,
			Sticker: &models.Sticker{
				Emoji: strings.Trim(emoji.Sprint(":smile:"), " "),
			},
		},
	}

	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
	}

	stickerHandler(clientObj, updateObj)
}

func TestMessageRandomWithUsername(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	testChat := models.Chat{
		ID: 100,
	}
	correct := fmt.Sprintf("<%s> Random Text", testUser.Username)

	updateObj := &models.Update{
		Message: &models.Message{
			From: testUser,
			Text: "Random Text",
			Chat: testChat,
		},
	}

	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
			ChatID: 100,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}

	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)

}

func TestMessageRandomWithoutUsername(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "",
		FirstName: "testing",
		LastName:  "123",
	}
	testChat := models.Chat{
		ID: 100,
	}

	correct := fmt.Sprintf("<%s> Random Text", testUser.FirstName)

	updateObj := &models.Update{
		Message: &models.Message{
			From: testUser,
			Text: "Random Text",
			Chat: testChat,
		},
	}
	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
			ChatID: 100,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: false,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}
	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)
}

func TestMessageRandomWithNoForward(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "",
		FirstName: "testing",
		LastName:  "123",
	}
	testChat := models.Chat{
		ID: 100,
	}

	updateObj := &models.Update{
		Message: &models.Message{
			From: testUser,
			Text: "[off] Random Text",
			Chat: testChat,
		},
	}
	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
			ChatID: 100,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP:        false,
			NoForwardPrefix: "[off]",
		},
		sendToIrc: func(s string) {
			assert.True(t, false)
		},
	}
	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)

}

func TestMessageZwsp(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	testChat := models.Chat{
		ID: 100,
	}
	correct := fmt.Sprintf("<%s> Random Text", "t"+"\u200b"+"est")

	updateObj := &models.Update{
		Message: &models.Message{
			From: testUser,
			Text: "Random Text",
			Chat: testChat,
		},
	}
	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
			ChatID: 100,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: true,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}

	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)
}

func TestMessageReply(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	replyUser := &models.User{
		ID:        2,
		Username:  "replyUser",
		FirstName: "Reply",
		LastName:  "User",
	}
	testChat := models.Chat{
		ID: 100,
	}

	tests := []struct {
		name     string
		updateFn func() *models.Update
		expected string
	}{
		{
			name: "ascii",
			updateFn: func() *models.Update {
				return &models.Update{
					Message: &models.Message{
						From: replyUser,
						Text: "Response Text",
						Chat: testChat,
						ReplyToMessage: &models.Message{
							From: testUser,
							Text: "Initial Text",
							Chat: testChat,
						},
					},
				}
			},
			expected: "<replyUser> [Re test: Initial Text] Response Text",
		},
		{
			name: "cyrillic-short",
			updateFn: func() *models.Update {
				return &models.Update{
					Message: &models.Message{
						From: replyUser,
						Text: "Response Text",
						Chat: testChat,
						ReplyToMessage: &models.Message{
							From: testUser,
							Text: "Тест",
							Chat: testChat,
						},
					},
				}
			},
			expected: "<replyUser> [Re test: Тест] Response Text",
		},
		{
			name: "cyrillic-long",
			updateFn: func() *models.Update {
				return &models.Update{
					Message: &models.Message{
						From: replyUser,
						Text: "Response Text",
						Chat: testChat,
						ReplyToMessage: &models.Message{
							From: testUser,
							Text: "Уикипедия е свободна енциклопедия",
							Chat: testChat,
						},
					},
				}
			},
			expected: "<replyUser> [Re test: Уикипедия е сво…] Response Text",
		},
		{
			name: "japanese-long",
			updateFn: func() *models.Update {
				return &models.Update{
					Message: &models.Message{
						From: replyUser,
						Text: "Response Text",
						Chat: testChat,
						ReplyToMessage: &models.Message{
							From: testUser,
							Text: "1234567テストテストテスト",
							Chat: testChat,
						},
					},
				}
			},
			expected: "<replyUser> [Re test: 1234567テストテストテス…] Response Text",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clientObj := &Client{
				Settings: &internal.TelegramSettings{
					Prefix:      "<",
					Suffix:      ">",
					ReplyPrefix: "[",
					ReplySuffix: "]",
					ReplyLength: 15,
					ChatID:      100,
				},
				IRCSettings: &internal.IRCSettings{
					ShowZWSP: false,
				},
				sendToIrc: func(actual string) {
					assert.Equal(t, test.expected, actual)
				},
			}
			messageHandler(clientObj)(clientObj.ctx, clientObj.API, test.updateFn())
		})
	}
}

func TestMessageReplyZwsp(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	replyUser := &models.User{
		ID:        2,
		Username:  "replyUser",
		FirstName: "Reply",
		LastName:  "User",
	}
	testChat := models.Chat{
		ID: 100,
	}
	initMessage := &models.Message{
		From: testUser,
		Text: "Initial Text",
		Chat: testChat,
	}
	correct := fmt.Sprintf("<%s> [Re %s: Initial Text] Response Text", "r"+"\u200b"+"eplyUser", "t"+"\u200b"+"est")

	updateObj := &models.Update{
		Message: &models.Message{
			From:           replyUser,
			Text:           "Response Text",
			Chat:           testChat,
			ReplyToMessage: initMessage,
		},
	}
	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix:      "<",
			Suffix:      ">",
			ReplyPrefix: "[",
			ReplySuffix: "]",
			ReplyLength: 15,
			ChatID:      100,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: true,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}

	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)
}

func TestMessageFromWrongTelegramChat(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}
	testChat := models.Chat{
		ID: 100,
	}

	updateObj := &models.Update{
		Message: &models.Message{
			From: testUser,
			Text: "Random Text",
			Chat: testChat,
		},
	}
	clientObj := &Client{
		Settings: &internal.TelegramSettings{
			Prefix: "<",
			Suffix: ">",
			ChatID: 101,
		},
		IRCSettings: &internal.IRCSettings{
			ShowZWSP: true,
		},
		sendToIrc: func(s string) {
			assert.False(t, true, "sendToIrc should not be called if the telegram chat ID has a mismatch")
		},
	}

	messageHandler(clientObj)(clientObj.ctx, clientObj.API, updateObj)
}

func TestLocationHandlerWithLocationEnabled(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}

	location := &models.Location{
		Latitude:  43.0845274,
		Longitude: -77.6781174,
	}

	correct := "test shared their location: (43.0845274, -77.6781174)."

	messageObj := models.Message{
		From:     testUser,
		Location: location,
	}
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			ShowZWSP:            false,
			ShowLocationMessage: true,
		},
		sendToIrc: func(s string) {
			assert.Equal(t, correct, s)
		},
	}

	locationHandler(clientObj, &messageObj)
}

func TestLocationHandlerWithLocationDisabled(t *testing.T) {
	testUser := &models.User{
		ID:        1,
		Username:  "test",
		FirstName: "testing",
		LastName:  "123",
	}

	location := &models.Location{
		Latitude:  43.0845274,
		Longitude: -77.6781174,
	}

	messageObj := models.Message{
		From:     testUser,
		Location: location,
	}
	clientObj := &Client{
		IRCSettings: &internal.IRCSettings{
			ShowLocationMessage: false,
		},
		sendToIrc: func(s string) {
			assert.Fail(t, "Setting disabled, this should not have been called")
		},
	}

	locationHandler(clientObj, &messageObj)
}
