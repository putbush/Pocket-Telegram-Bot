package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"pocketer_bot/pkg/storage"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenReposiroty(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (t *TokenRepository) Save(chatID int64, requestKey string, bucket storage.Bucket) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put(intToBytes(chatID), []byte(requestKey))
	})
}

func (t *TokenRepository) Get(chatID int64, bucket storage.Bucket) (string, error) {
	var token string

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, err
}

func intToBytes(chatID int64) []byte {
	return []byte(strconv.FormatInt(chatID, 10))
}
