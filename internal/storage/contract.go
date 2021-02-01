package storage

import (
	"errors"

	"dev/tgbot/cmd/service/internal/models"
)

type Storage interface {
	Delete(chatID int64) error
	Get(chatID int64) (models.Chat, error)
	Save(chat models.Chat) error
}

var (
	ErrChatNotFound = errors.New("chat not found")
)