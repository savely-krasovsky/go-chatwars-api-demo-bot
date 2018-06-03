package main

import (
	"encoding/json"
	"github.com/L11R/go-chatwars-api"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

// Logger and Bot global variables
var (
	sugar  *zap.SugaredLogger
	bot    *tgbotapi.BotAPI
	db     *storm.DB
	client *cwapi.Client
)

// Command handler for logging purposes
func (u Update) Handle(c func(update Update) (*tgbotapi.Message, error)) {
	res, err := c(u)
	if err != nil {
		sugar.Error(err)
	}

	if res != nil {
		b, err := json.MarshalIndent(res, "", "\t")
		if err != nil {
			sugar.Warn(err)
		}

		sugar.Debug(string(b))
	}
}

func main() {
	if err := Init(); err != nil {
		sugar.Fatalf("initialization failed: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		sugar.Fatal(err)
	}

	for u := range updates {
		// Extend update
		update := Update{
			&u,
		}

		// Skip images, stickers, etc
		if update.Message == nil {
			continue
		}

		// Commands handler
		switch update.Message.Command() {
		case "start":
			go update.Handle(Start)
		case "auth":
			go update.Handle(Auth)
		case "profile":
			go update.Handle(Profile)
		default:
			// Check login
			if update.Message.ForwardFrom.UserName == "chtwrsbot" {
				go update.Handle(Code)
			}
		}
	}
}
