package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
)

func Start(update Update) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Hello! It's demo bot powered by <a href=\"https://github.com/L11R/go-chatwars-api\">go-chatwars-bot</a>.\nClick /auth to see how it works.",
	)
	msg.ParseMode = "HTML"

	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func Auth(update Update) (*tgbotapi.Message, error) {
	// Send create auth code request
	cwres, err := client.CreateAuthCodeSync(update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	// Debug log
	sugar.Debug(cwres)

	// Send user message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Code has sent, please forward message with code!")
	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func Code(update Update) (*tgbotapi.Message, error) {
	// Regexp to find six digit code
	re := regexp.MustCompile(`\d{6}`)

	// Send code to get token
	cwres, err := client.GrantTokenSync(update.Message.From.ID, re.FindString(update.Message.Text))
	if err != nil {
		return nil, err
	}

	// Debug log
	sugar.Debug(cwres)

	// Save user with token into DB
	if err = db.Save(&User{
		ID:    update.Message.From.ID,
		Token: cwres.Payload.Token,
	}); err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Thanks! Authenticated successfully, now try to retrieve you /profile!",
	)
	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func Profile(update Update) (*tgbotapi.Message, error) {
	// Get saved user from DB
	var user User
	if err := db.One("ID", update.Message.From.ID, &user); err != nil {
		return nil, err
	}

	// Request profile
	cwres, err := client.RequestProfileSync(user.Token, user.ID)
	if err != nil {
		return nil, err
	}

	// Debug log
	sugar.Debug(cwres)

	// Make profile text
	text := fmt.Sprintf(
		`<b>%s%s</b>
Class: %s
🏅Level: %d
⚔️Atk: %d 🛡Def: %d
🔥Exp: %d
🔋Stamina: %d`,
		cwres.Payload.Profile.Castle,
		cwres.Payload.Profile.UserName,
		cwres.Payload.Profile.Class,
		cwres.Payload.Profile.Level,
		cwres.Payload.Profile.Attack,
		cwres.Payload.Profile.Defense,
		cwres.Payload.Profile.Experience,
		cwres.Payload.Profile.Stamina,
	)

	// Check, maybe it's not Blacksmith
	if cwres.Payload.Profile.Mana != 0 {
		text += fmt.Sprintf("\n💧Mana: %d", cwres.Payload.Profile.Mana)
	}

	text += fmt.Sprintf("\n💰%d 👝%d", cwres.Payload.Profile.Gold, cwres.Payload.Profile.Pouches)

	// ... or not in the Guild?
	if cwres.Payload.Profile.Guild != "" {
		text += fmt.Sprintf("\n👥%s", cwres.Payload.Profile.Guild)
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		text,
	)
	msg.ParseMode = "HTML"

	res, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
