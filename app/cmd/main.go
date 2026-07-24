package main

import (
	"gochatroom/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	chatRoom := server.NewChatRoom()
	go chatRoom.Run()

	r := gin.Default()
	r.GET("/ws", chatRoom.HandleWebSocket)
    r.Run(":8000")
}