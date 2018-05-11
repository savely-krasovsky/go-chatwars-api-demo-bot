package main

import (
	"github.com/L11R/go-chatwars-api"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init configuration manager, logger, bot, database
func Init() error {
	// Init and read config file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Configuration defaults
	// Log level: INFO (-1 for DEBUG)
	viper.SetDefault("log.level", 0)
	// Log type: "production" or "development"
	viper.SetDefault("log.type", "production")

	// Init logger
	var loggerConfig zap.Config
	if viper.GetString("log.type") == "production" {
		loggerConfig = zap.NewProductionConfig()
	}
	if viper.GetString("log.type") == "development" {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.Level.SetLevel(zapcore.Level(viper.GetInt("log.level")))

	logger, _ := loggerConfig.Build()
	defer logger.Sync()

	sugar = logger.Sugar()

	// Init database
	db, err = storm.Open(viper.GetString("db.name"))
	if err != nil {
		return err
	}

	// Check token
	if !viper.IsSet("token") {
		return err
	}

	// Init Telegram Bot API
	bot, err = tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		return err
	}

	sugar.Infof("authorized on @%s", bot.Self.UserName)

	// Init Chat Wars API client
	client, err = cwapi.NewClient(viper.GetString("cw.user"), viper.GetString("cw.password"))
	if err != nil {
		return err
	}

	return nil
}
