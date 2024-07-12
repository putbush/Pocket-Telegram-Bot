package main

import (
	"github.com/boltdb/bolt"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"pocketer_bot/pkg/config"
	"pocketer_bot/pkg/server"
	"pocketer_bot/pkg/storage"
	"pocketer_bot/pkg/storage/boltdb"
	"pocketer_bot/pkg/telegram"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Println(1)
		log.Fatal(err)
	}

	client, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	tokenRepository := boltdb.NewTokenReposiroty(db)

	go func() {
		telegramBot := telegram.NewBot(bot, client, cfg.AuthServerURL, tokenRepository, cfg.Message)
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	authorizationServer := server.NewAuthorizationServer(client, tokenRepository, cfg.TelegramBotURL)
	if err = authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(storage.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
