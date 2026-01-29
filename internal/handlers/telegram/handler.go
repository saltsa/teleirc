package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

/*
Handler specifies a function that handles a Telegram update.
In this case, we take a Telegram client and update object,
where the specific Handler will "handle" the given event.
*/
type Handler = func(tg *Client, u models.Update)

/*
messageHandler handles the Message Telegram Object, which formats the
Telegram update into a simple string for IRC.
*/
func messageHandler(tg *Client) tgbotapi.HandlerFunc {
	return func(ctx context.Context, b *tgbotapi.Bot, u *models.Update) {
		date := time.Unix(int64(u.Message.Date), 0)
		if since := time.Since(date); since > time.Minute {
			tg.logger.LogWarning("received message was %s old, ignoring sending", since)
			return
		}
		username := GetUsername(tg.IRCSettings.ShowZWSP, u.Message.From)
		formatted := ""

		if tg.IRCSettings.NoForwardPrefix != "" && strings.HasPrefix(u.Message.Text, tg.IRCSettings.NoForwardPrefix) {
			return
		}

		// Don't forward messages to IRC that didn't come from the
		// chat we're bridging
		if u.Message.Chat.ID != tg.Settings.ChatID {
			return
		}

		// Telegram user replied to a message
		if u.Message.ReplyToMessage != nil {
			replyHandler(tg, u)
			return
		}

		formatted = fmt.Sprintf("%s%s%s %s",
			tg.Settings.Prefix,
			username,
			tg.Settings.Suffix,
			// Trim unexpected trailing whitespace
			strings.Trim(u.Message.Text, " "))

		tg.sendToIrc(formatted)
	}
}

/*
replyHandler handles when users reply to a Telegram message
*/
func replyHandler(tg *Client, u *models.Update) {
	replyText := strings.Trim(u.Message.ReplyToMessage.Text, " ")
	username := GetUsername(tg.IRCSettings.ShowZWSP, u.Message.From)
	replyUser := GetUsername(tg.IRCSettings.ShowZWSP, u.Message.ReplyToMessage.From)

	// Only show a portion of the reply text
	if replyTextAsRunes := []rune(replyText); len(replyTextAsRunes) > tg.Settings.ReplyLength {
		replyText = string(replyTextAsRunes[:tg.Settings.ReplyLength]) + "…"
	}

	var replyMsg string

	if u.Message.ReplyToMessage.IsTopicMessage {
		// If message was sent to Forum Topic
		if u.Message.ReplyToMessage.ForumTopicCreated != nil {
			// It was directly to topic, ie. not a reply
			replyMsg = fmt.Sprintf("%sTopic: %s%s",
				tg.Settings.ReplyPrefix,
				u.Message.ReplyToMessage.ForumTopicCreated.Name,
				tg.Settings.ReplySuffix)
		} else {
			// It was a reply in topic so we do not know the topic name
			replyMsg = fmt.Sprintf("%sTopic Re %s: %s%s",
				tg.Settings.ReplyPrefix,
				replyUser,
				replyText,
				tg.Settings.ReplySuffix)
		}
	} else {
		// Reply in generic channel
		replyMsg = fmt.Sprintf("%sRe %s: %s%s",
			tg.Settings.ReplyPrefix,
			replyUser,
			replyText,
			tg.Settings.ReplySuffix)
	}

	formatted := fmt.Sprintf("%s%s%s %s %s",
		tg.Settings.Prefix,
		username,
		tg.Settings.Suffix,
		replyMsg,
		u.Message.Text)

	tg.sendToIrc(formatted)
}

/*
joinHandler handles when users join the Telegram group
*/
func joinHandler(tg *Client, users *[]models.User) {
	if tg.IRCSettings.ShowJoinMessage {
		for _, user := range *users {
			user := user
			username := GetFullUsername(tg.IRCSettings.ShowZWSP, &user)
			formatted := username + " has joined the Telegram Group!"
			tg.sendToIrc(formatted)
		}
	}
}

/*
partHandler handles when users leave the Telegram group
*/
func partHandler(tg *Client, user *models.User) {
	if tg.IRCSettings.ShowLeaveMessage {
		username := GetFullUsername(tg.IRCSettings.ShowZWSP, user)
		formatted := username + " has left the Telegram Group!"

		tg.sendToIrc(formatted)
	}
}

/*
stickerHandler handles the Message.Sticker Telegram Object, which formats the
Telegram message into its base Emoji unicode character.
*/
func stickerHandler(tg *Client, u models.Update) {
	username := GetUsername(tg.IRCSettings.ShowZWSP, u.Message.From)
	formatted := fmt.Sprintf("%s%s%s %s",
		tg.Settings.Prefix,
		username,
		tg.Settings.Suffix,
		u.Message.Sticker.Emoji)
	tg.sendToIrc(formatted)
}

/*
photoHandler handles the Message.Photo Telegram object. Only acknowledges Photo
exists, and sends notification to IRC
*/
// TODO: Not working yet
// func photoHandler(tg *Client, u models.Update) {
// 	link := uploadImage(tg, u)
// 	username := GetUsername(tg.IRCSettings.ShowZWSP, u.Message.From)
// 	caption := u.Message.Caption
// 	if caption == "" {
// 		caption = "No caption provided."
// 	}

// 	// TeleIRC can fail to upload to Imgur
// 	if link == "" {
// 		tg.logger.LogError("Failed imgur photo upload for", username)
// 	} else {
// 		formatted := "'" + caption + "' uploaded by " + username + ": " + link
// 		tg.sendToIrc(formatted)
// 	}
// }

/*
documentHandler receives a document object from Telegram, and sends
a notification to IRC.
*/
func documentHandler(tg *Client, u *models.Message) {
	username := GetUsername(tg.IRCSettings.ShowZWSP, u.From)
	formatted := username + " shared a file"
	if u.Document.MimeType != "" {
		formatted += " (" + u.Document.MimeType + ")"
	}

	if u.Caption != "" {
		formatted += " on Telegram with caption: " + "'" + u.Caption + "'."
	} else if u.Document.FileName != "" {
		formatted += " on Telegram with title: " + "'" + u.Document.FileName + "'."
	}

	tg.sendToIrc(formatted)
}

/*
locationHandler receivers a location object from Telegram, and sends
a notification to IRC.
*/
func locationHandler(tg *Client, u *models.Message) {
	if !tg.IRCSettings.ShowLocationMessage {
		return
	}

	username := GetUsername(tg.IRCSettings.ShowZWSP, u.From)
	formatted := username + " shared their location: ("

	// f means do not use an exponent.
	// -1 means use the smallest number of digits needed so parseFloat will return f exactly.
	// 64 to represent a standard 64 bit floating point number.
	formatted += strconv.FormatFloat(u.Location.Latitude, 'f', -1, 64)
	formatted += ", "
	formatted += strconv.FormatFloat(u.Location.Longitude, 'f', -1, 64)
	formatted += ")."

	tg.sendToIrc(formatted)
}
