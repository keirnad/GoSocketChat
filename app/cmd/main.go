package main

import (
	"gochatroom/controllers"
	"gochatroom/initializers"
	"gochatroom/internal/server"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	r := gin.Default()
	chatRoom := server.NewChatRoom()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.GET("/ws", chatRoom.HandleWebSocket)

	r.Run(":" + os.Getenv("PORT"))

	go chatRoom.Run()
}
