package server

import (
	"sync"
)

type ChatRoom struct {
	clients *sync.Map
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		broadcast: make(chan []byte, 256),
		register: make(chan *Client),       
		unregister: make(chan *Client),
		clients: &sync.Map{},
	}
}

func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <- cr.register:
			cr.clients.Store(client.nickname, client)
		case client := <- cr.unregister:
			cr.clients.Delete(client)
		case message := <-cr.broadcast:
			cr.clients.Range(func(key, value interface{}) bool {
				client := value.(*Client)
				select {
				case client.send <- []byte(client.nickname + ": " + string(message)):
				default: 
					close(client.send)
					cr.clients.Delete(key)
				}
				return true
			})
		}
	}
}