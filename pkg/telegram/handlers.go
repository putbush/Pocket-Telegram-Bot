package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	start = "start"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavesSuccessfully)
	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{URL: message.Text, AccessToken: accessToken}); err != nil {
		return errUnableToSave
	}
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case start:
		err := b.handleStartCommand(message)
		return err
	default:
		err := b.handleUnknownCommand()
		return err
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	if _, err := b.getAccessToken(message.Chat.ID); err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuth)
	_, err := b.bot.Send(msg)
	return err

}

func (b *Bot) handleUnknownCommand() error {
	return errUnknownCommand
}
