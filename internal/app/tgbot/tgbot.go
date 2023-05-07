package tgbot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vk_tarantool_test_task/internal/app/config"
	"vk_tarantool_test_task/internal/infrastructure"
)

const setCommand = "/set"
const setCommandFormat = "/set service login password"
const getCommand = "/get"
const getCommandFormat = "/get service"
const delCommand = "/del"
const delCommandFormat = "/del service"

type TgBot struct {
	bot *tgbotapi.BotAPI
	db  *infrastructure.Database
}

func New(cfg *config.Config, db *infrastructure.Database) (*TgBot, error) {
	tgBot := &TgBot{}
	var err error

	tgBot.bot, err = tgbotapi.NewBotAPI(cfg.Bot.TelegramApiToken)
	if err != nil {
		return nil, err
	}

	// TODO: change debug mode
	tgBot.bot.Debug = true
	log.Printf("Authorized on account %s", tgBot.bot.Self.UserName)

	tgBot.db = db

	return tgBot, nil
}

func (b *TgBot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var err error
			if update.Message.Text[:len(setCommand)] == setCommand {
				err = b.setCommandHandle(ctx, update.Message)
			} else if update.Message.Text[:len(getCommand)] == getCommand {
				err = b.getCommandHandle(ctx, update.Message)
			} else if update.Message.Text[:len(delCommand)] == delCommand {
				err = b.delCommandHandle(ctx, update.Message)
			}
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

const setCommandSuccessMessage = "Credentials successfully added"

func (b *TgBot) setCommandHandle(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	parts := strings.Split(message.Text, " ")
	if len(parts) != len(strings.Split(setCommandFormat, " ")) {
		return fmt.Errorf("incorrect command format")
	}
	service, login, password := parts[1], parts[2], parts[3]

	err := b.db.InsertCredentialsHandler(ctx, userID, service, login, password)
	if err != nil {
		return err
	}

	err = b.replyToUser(message, setCommandSuccessMessage)
	if err != nil {
		return err
	}

	return nil
}

const getCommandCantFindMessage = "Can't find credentials for service %s"
const getCommandSuccessMessage = "service: %s\nlogin: %s\npassword: %s"

func (b *TgBot) getCommandHandle(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	parts := strings.Split(message.Text, " ")
	if len(parts) != len(strings.Split(getCommandFormat, " ")) {
		return fmt.Errorf("incorrect command format")
	}
	service := parts[1]

	creds, err := b.db.GetCredentialsHandler(ctx, userID, service)
	if err != nil {
		return err
	}

	if creds == nil {
		err = b.replyToUser(message, fmt.Sprintf(getCommandCantFindMessage, service))
	} else {
		for _, cred := range creds {
			err = b.replyToUser(message, fmt.Sprintf(getCommandSuccessMessage, cred.Service, cred.Login, cred.Password))
		}
	}
	if err != nil {
		return err
	}

	return nil
}

const delCommandSuccessMessage = "Credentials successfully deleted"

func (b *TgBot) delCommandHandle(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	parts := strings.Split(message.Text, " ")
	if len(parts) != len(strings.Split(delCommandFormat, " ")) {
		return fmt.Errorf("incorrect command format")
	}
	service := parts[1]

	err := b.db.DeleteCredentialsHandler(ctx, userID, service)
	if err != nil {
		return err
	}

	err = b.replyToUser(message, delCommandSuccessMessage)
	if err != nil {
		return err
	}

	return nil
}

func (b *TgBot) replyToUser(message *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}