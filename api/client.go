package api

import (
	"github.com/gorilla/websocket"
	"log"
)

type client struct {

	// Socket is the web socket for this client
	socket *websocket.Conn

	receive chan []byte

	room *room

	username string

	roomId int64
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading the message: %v", err)
			return
		}
		log.Printf("Message read: %s", string(msg))
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		log.Printf("Writing message: %s", string(msg))
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error writing the message: %v", err)
			return
		}
	}
}
