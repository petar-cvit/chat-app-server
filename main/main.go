package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	_ "github.com/joho/godotenv/autoload"

	"github.com/petar-cvit/chat-app-server/internal/infrastructure/storage"
	"github.com/petar-cvit/chat-app-server/internal/models"
)

func main() {
	router := gin.New()
	storage := storage.New()

	server, _ := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		room := storage.GetRoom(s.ID())

		storage.SetRoom(s.ID(), room)

		s.SetContext("")
		s.Join(room)
		s.Emit("new_room", "lobby")

		msgs, err := storage.GetMessagesByRoom(room)
		if err != nil {
			return err
		}

		for _, msg := range msgs.Messages {
			s.Emit("reply", msg)
		}

		return nil
	})

	server.OnEvent("/", "joinRoom", func(conn socketio.Conn, roomID string) error {
		roomName := roomName(roomID)
		room := strings.Join([]string{"chat_room", roomName}, "_")

		conn.Leave(fmt.Sprintf("chat_room_%v", storage.GetRoom(conn.ID())))
		conn.Emit("clear", "clear")

		conn.Join(room)

		if storage.SetRoom(conn.ID(), room) {
			conn.Emit("new_room", roomName)
		}

		msgs, err := storage.GetMessagesByRoom(room)
		if err != nil {
			return err
		}

		for _, msg := range msgs.Messages {
			marshalledMessage, err := json.Marshal(msg)
			if err != nil {
				return err
			}

			conn.Emit("reply", string(marshalledMessage))
		}

		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, text string) {
		parts := strings.Split(text, ";")

		message := &models.Message{
			Text:   parts[0],
			Time:   time.Now().Format(time.Stamp),
			Issuer: parts[1],
		}

		storage.SaveMessage(storage.GetRoom(s.ID()), message)

		marshalledMessage, err := json.Marshal(message)
		if err != nil {
			return
		}

		server.BroadcastToRoom("/", storage.GetRoom(s.ID()), "reply", string(marshalledMessage))
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Leave(storage.GetRoom(s.ID()))
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	router.LoadHTMLGlob("./templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index_chat.html", gin.H{
			"title": "Chat app",
		})
	})

	router.Static("./css", "./css")

	router.Run()
}

func roomName(name string) string {
	if len(name) == 0 {
		return "lobby"
	}

	return name
}
