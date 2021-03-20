package storage

import (
	"sync"
)

type Storage struct {
	storage sync.Map
}

func New() *Storage {
	return &Storage{
		storage: sync.Map{},
	}
}

func (s *Storage) SaveMessage(ID string, message string) {
	messages := s.GetMessagesByRoom(ID)
	messages = append(messages, message)

	s.storage.Store(ID, messages)
}

func (s *Storage) GetMessagesByRoom(ID string) []string {
	messages, ok := s.storage.Load(ID)
	if !ok {
		return []string{}
	}

	mess, ok := messages.([]string)
	if !ok {
		return []string{}
	}

	return mess
}

func (s *Storage) SetRoom(userID, roomID string) {
	s.storage.Store(userID, roomID)
}

func (s *Storage) GetRoom(userID string) string {
	val, exists := s.storage.Load(userID)
	if !exists {
		return "chat_room"
	}

	roomID, ok := val.(string)
	if !ok {
		return "chat_room"
	}

	return roomID
}
