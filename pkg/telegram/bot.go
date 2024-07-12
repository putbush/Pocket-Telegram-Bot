package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"pocketer_bot/pkg/config"
	"pocketer_bot/pkg/storage"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenReposiroty storage.TokenRepository
	redirectURL     string
	messages        config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, redirectURL string, tr storage.TokenRepository, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: client, redirectURL: redirectURL, tokenReposiroty: tr, messages: messages}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	b.handleUpdates()

	return nil
}

func (b *Bot) initUpdateChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	return b.bot.GetUpdatesChan(u)

}

func (b *Bot) handleUpdates() {
	updates := b.initUpdateChannel()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message, err)
		}
	}
}
