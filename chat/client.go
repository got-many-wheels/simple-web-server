package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room where the client is chatting in
	room *room
	// userData holds informnation about the user
	userData map[string]any
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.Email = c.userData["email"].(string)
		msg.AvatarURL, err = c.room.avatar.GetAvatarURL(c)
		if err != nil {
			c.room.tracer.Trace(err)
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
