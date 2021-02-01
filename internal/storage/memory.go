package storage

import (
	"sync"

	"dev/tgbot/cmd/service/internal/models"
)

type memory struct {
	mux  sync.Mutex
	data map[int64]models.Chat
}

func NewMemoryStorage() Storage {
	return &memory{
		data: make(map[int64]models.Chat),
	}
}

func (m *memory) Delete(chatID int64) error {
	m.mux.Lock()
	delete(m.data, chatID)
	m.mux.Unlock()

	return nil
}

func (m *memory) Get(chatID int64) (models.Chat, error) {
	m.mux.Lock()
	chat, ok := m.data[chatID]
	m.mux.Unlock()

	if ok {
		return chat, nil
	}
	return models.Chat{}, ErrChatNotFound
}

func (m *memory) Save(chat models.Chat) error {
	m.mux.Lock()
	m.data[chat.ID] = chat
	m.mux.Unlock()
	return nil
}
