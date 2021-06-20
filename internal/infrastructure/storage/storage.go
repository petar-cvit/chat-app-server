package storage

import (
	"github.com/petar-cvit/chat-app-server/internal/models"
	"sync"
)

type Storage struct {
	messages sync.Map
	rooms    sync.Map
}

func New() *Storage {
	return &Storage{
		messages: sync.Map{},
		rooms:    sync.Map{},
	}
}

func (s *Storage) SaveMessage(ID string, message string) {
	messages := s.GetMessagesByRoom(ID)
	messages = append(messages, message)

	s.messages.Store(ID, messages)
}

func (s *Storage) GetMessagesByRoom(ID string) []string {
	messages, ok := s.messages.Load(ID)
	if !ok {
		return []string{}
	}

	mess, ok := messages.([]string)
	if !ok {
		return []string{}
	}

	return mess
}

func (s *Storage) SetRoom(userID, roomID string) bool {
	s.messages.Store(userID, roomID)

	data, exists := s.rooms.Load(userID)

	if !exists {
		s.rooms.Store(userID, models.Rooms{map[string]bool{
			roomID: true,
		}})

		return true
	}

	roomsByUser := data.(models.Rooms)

	if _, visitedThisRoom := roomsByUser.Rooms[roomID]; visitedThisRoom {
		return false
	}

	roomsByUser.Rooms[roomID] = true
	s.rooms.Store(userID, roomsByUser)
	return true
}

func (s *Storage) GetRoom(userID string) string {
	val, exists := s.messages.Load(userID)
	if !exists {
		return "chat_room"
	}

	roomID, ok := val.(string)
	if !ok {
		return "chat_room"
	}

	return roomID
}
