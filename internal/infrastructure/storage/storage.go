package storage

import (
	"errors"
	"sync"

	"github.com/petar-cvit/chat-app-server/internal/models"
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

func (s *Storage) SaveMessage(ID string, message *models.Message) error {
	messages, err := s.GetMessagesByRoom(ID)
	if err != nil {
		return err
	}

	messages.AppendMessage(message)
	s.messages.Store(ID, messages)

	return nil
}

func (s *Storage) GetMessagesByRoom(ID string) (*models.Messages, error) {
	messages, ok := s.messages.Load(ID)
	if !ok {
		messages = &models.Messages{Messages: []*models.Message{}}
	}

	mess, ok := messages.(*models.Messages)
	if !ok {
		return nil, errors.New("casting messages error")
	}

	return mess, nil
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

func (s *Storage) GetRoomsByUser(userID string) models.Rooms {
	data, ok := s.rooms.Load(userID)
	if !ok {
		return models.Rooms{}
	}

	roomsByUser := data.(models.Rooms)

	return roomsByUser
}

func (s *Storage) GetRoom(userID string) string {
	val, exists := s.messages.Load(userID)
	if !exists {
		return "chat_room_lobby"
	}

	roomID, ok := val.(string)
	if !ok {
		return "chat_room"
	}

	return roomID
}
