package main

import (
	"gochatroom/controllers"
	"gochatroom/initializers"
	"gochatroom/internal/server"
	"gochatroom/middleware"
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
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/ws", chatRoom.HandleWebSocket)

	r.Run(":" + os.Getenv("PORT"))

	go chatRoom.Run()
}
