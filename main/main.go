package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	_ "github.com/joho/godotenv/autoload"

	"github.com/petar-cvit/chat-app-server/internal/infrastructure/storage"
)

func main() {
	router := gin.New()
	storage := storage.New()

	server, _ := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		room := fmt.Sprintf("chat_room_%v", storage.GetRoom(s.ID()))

		storage.SetRoom(s.ID(), room)

		s.SetContext("")
		s.Join(room)

		msgs := storage.GetMessagesByRoom(room)
		for _, msg := range msgs {
			s.Emit("reply", msg)
		}

		return nil
	})

	server.OnEvent("/", "joinRoom", func(conn socketio.Conn, roomID string) error {
		room := fmt.Sprintf("chat_room_%v", roomID)

		conn.Leave(fmt.Sprintf("chat_room_%v", storage.GetRoom(conn.ID())))
		conn.Emit("clear", "clear")

		conn.Join(room)

		if storage.SetRoom(conn.ID(), room) {
			conn.Emit("new_room", roomID)
		}

		msgs := storage.GetMessagesByRoom(room)
		for _, msg := range msgs {
			conn.Emit("reply", msg)
		}

		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		storage.SaveMessage(storage.GetRoom(s.ID()), msg)

		server.BroadcastToRoom("/", storage.GetRoom(s.ID()), "reply", msg)
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
