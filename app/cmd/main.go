package main

import (
	"gochatroom/initializers"
	"gochatroom/internal/server"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	chatRoom := server.NewChatRoom()
	go chatRoom.Run()

	r := gin.Default()
	r.GET("/ws", chatRoom.HandleWebSocket)
	r.Run(":8000")
}
