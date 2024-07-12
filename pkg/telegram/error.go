package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidURL     = errors.New("url is invalid")
	errUnauthorized   = errors.New("user is not authorized")
	errUnableToSave   = errors.New("unable to save")
	errUnknownCommand = errors.New("command is unknown")
)

func (b *Bot) handleError(message *tgbotapi.Message, err error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Default)
	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
	case errUnknownCommand:
		msg.Text = b.messages.UnknownCommand
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
	}
	_, _ = b.bot.Send(msg)
}
