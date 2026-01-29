package telegram

import (
	"github.com/go-telegram/bot/models"
)

/*
GetUsername takes showZWSP condition and user then returns username with or without ​.
*/
func GetUsername(showZWSP bool, u *models.User) string {
	if u.Username == "" {
		return u.FirstName
	}
	if showZWSP {
		return ZwspUsername(u)
	}
	return u.Username
}

/*
GetFullUsername takes showZWSP condition and user then returns full username with or without ​.
*/
func GetFullUsername(showZWSP bool, u *models.User) string {
	if u.Username == "" {
		return u.FirstName
	}
	if showZWSP {
		return GetFullUserZwsp(u)
	}
	return u.FirstName + " (@" + u.Username + ")"
}

/*
GetFullUserZwsp returns both the Telegram user's first name and username, if available.
Adds ZWSP to username to prevent username pinging across platform.
*/
func GetFullUserZwsp(u *models.User) string {
	// Add ZWSP to prevent pinging across platforms
	// See https://github.com/42wim/matterbridge/issues/175
	userNameAsRunes := []rune(u.Username)
	return u.FirstName + " (@" + string(userNameAsRunes[:1]) + "\u200b" + string(userNameAsRunes[1:]) + ")"
}

/*
ZwspUsername adds a zero-width space after the first character of a Telegram user's
username.
*/
func ZwspUsername(u *models.User) string {
	// Add ZWSP to prevent pinging across platforms
	// See https://github.com/42wim/matterbridge/issues/175
	userNameAsRunes := []rune(u.Username)
	return string(userNameAsRunes[:1]) + "\u200b" + string(userNameAsRunes[1:])
}

/*
uploadImage uploads a Photo object from Telegram to the Imgur API and
returns a string with the Imgur URL.
*/
// TODO: Fix this if needed
// func uploadImage(tg *Client, u models.Update) string {
// 	photo := (u.Message.Photo)[len(u.Message.Photo)-1]

// 	// Get Telegram file URL
// 	tgLink, err := tg.api.FileDownloadLink(photo.FileID)
// 	if err != nil {
// 		tg.logger.LogError("Could not get Telegram Photo URL:", err)
// 	}

// 	return getImgurLink(tg, tgLink)
// }
