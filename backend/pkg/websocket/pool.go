package websocket

import (
	"fmt"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	Name       string
	OpList     []*Client
	MuteList   []*Client
	BanList    []*Client
	PoolSize   int
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Name:       "Default",
		OpList:     make([]*Client, 0),
		MuteList:   make([]*Client, 0),
		BanList:    make([]*Client, 0),
		PoolSize:   100,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		//adds a user to a chatroom
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			//tells all users in the chatroom a new user joined
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined #" + pool.Name})
			}
			break
		//removes a user from a chatroom
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			//tells all users in the chatroom a user has left
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected from #" + pool.Name})
			}
			break
		//sends all user in a chatroom the msg
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				fmt.Println(client.ID + " " + pool.Name)
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
