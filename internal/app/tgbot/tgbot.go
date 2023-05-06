package tgbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	"vk_tarantool_test_task/internal/app/config"
)

type TgBot struct {
	bot *tgbotapi.BotAPI
}

func New(cfg *config.Config) (*TgBot, error) {
	tgBot := &TgBot{}
	var err error

	tgBot.bot, err = tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		return nil, err
	}

	// TODO: change debug mode
	tgBot.bot.Debug = true
	log.Printf("Authorized on account %s", tgBot.bot.Self.UserName)

	return tgBot, nil
}

func (b *TgBot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			b.bot.Send(msg)
		}
	}

	return nil
}
