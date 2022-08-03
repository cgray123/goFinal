package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

//struct of users
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool []*Pool
	Bot  *Bot
}

//helper struct that holds all chatrooms, admins, and users
type Bot struct {
	ID         string
	AllPool    []*Pool
	AllClients []*Client
	AllAdmin   []*Client
}

//msg struct stores data is a json
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	//remove user from all chatrooms and closes the user websocket if the browser is closed
	defer func() {
		for i := 0; i < len(c.Pool); i++ {
			c.Pool[i].Unregister <- c
		}
		c.Conn.Close()
	}()

	for {
		//catchs all msg sent
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
			//handles empty inputs
		} else if string(p) == "" {
			message := Message{Type: messageType, Body: "@BOT: Message must have text"}
			c.singleMes(message)
			//handle setting names
		} else if p[0] == '~' {
			var notU bool
			var count int
			var name string
			if string(p)[1:] == "" {
				name = "anonymous"
			} else {
				name = string(p)[1:]
			}

			//checks for dup names and adds a number to make it unique
			tName := name
			for _, y := range c.Bot.AllClients {
				if y.ID == tName {
					notU = true
					count++
					tName = name + fmt.Sprint(count)
				}
			}
			if notU {
				name += fmt.Sprint(count)
			}
			c.ID = name
			message := Message{Type: messageType, Body: "@BOT: Name set to @" + c.ID}
			c.singleMes(message)
			fmt.Printf("Name Set: %+v\n", name)

			//handles all commands
		} else if p[0] == '/' {
			c.Command(messageType, p)
		} else {
			//sends msg to all chatrooms the user is in
			for i := 0; i < len(c.Pool); i++ {
				var isMuted bool
				//checks if user is muted
				for _, x := range c.Pool[i].MuteList {
					if x == c {
						isMuted = true
					}
				}
				if !isMuted {
					message := Message{Type: messageType, Body: "#" + c.Pool[i].Name + ": @" + c.ID + ": " + string(p)}
					c.Pool[i].Broadcast <- message
					fmt.Printf("Message Received: %+v\n", message)
				} else {
					message := Message{Type: messageType, Body: "@BOT: You are muted in chatroom #" + c.Pool[i].Name}
					c.singleMes(message)
				}

			}

		}

	}
}

//sends a msg just to the user
func (c *Client) singleMes(m Message) {
	if err := c.Conn.WriteJSON(m); err != nil {
		fmt.Println(err)
	}
}
