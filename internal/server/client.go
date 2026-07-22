package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	  // Allow all origins for development
  	CheckOrigin: func(r *http.Request) bool {
    	return true
  	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) Read(cr *ChatRoom) {
	defer func() {
		cr.unregister <- c 
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		cr.broadcast <- []byte(string(message))
	}
}

func (c *Client) Write() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func (cr *ChatRoom) HandleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }

    client := &Client{conn: conn, send: make(chan []byte, 256)}
    cr.register <- client

    go client.Write()
    go client.Read(cr)
}