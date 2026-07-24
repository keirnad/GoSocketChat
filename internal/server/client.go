package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Nickname string `json:"nickname"`
	Body string `json:"body"`
}

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
	nickname string
	send chan []byte
}

func (c *Client) Read(cr *ChatRoom) {
	defer func() {
		cr.unregister <- c 
		c.conn.Close()
	}()

	var msg Message

	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error: ", err)
		}
		payload := Message{Nickname: msg.Nickname, Body: msg.Body}
		jsonmsg, err := json.Marshal(payload)
		if err != nil {
			log.Print("Eror")
		}
		cr.broadcast <- []byte(jsonmsg)
	}
}

func (c *Client) Write() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteJSON(message); err != nil {
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


	nickname := c.Query("nickname")

	if c.Query("nickname") == "" {
		nickname = "Annonymus"
	}

    client := &Client{conn: conn, send: make(chan []byte, 256), nickname: nickname}
    cr.register <- client

    go client.Write()
    go client.Read(cr)
}