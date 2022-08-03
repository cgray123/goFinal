package websocket

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//handles command called by entering a / in the input
func (c *Client) Command(messageType int, p []byte) {
	//returns the current time
	if string(p[1:]) == "time" {
		message := Message{Type: messageType, Body: "@BOT: The current date and time is " + time.Now().Format("Mon Jan _2 15:04:05 2006")}
		c.singleMes(message)
		//send a direct chat to the user after the @
	} else if string(p)[1:4] == "dc_" {
		var has bool
		for _, b := range string(p[4:]) {
			if b == ':' {
				has = true
			}
		}
		if has {
			//loop through all clients
			for _, client := range c.Bot.AllClients {
				//finds the client and sends a msg to both the sender and the receiver
				if string(p)[4:strings.Index(string(p), ":")] == client.ID {
					if err := c.Conn.WriteJSON(Message{Type: messageType, Body: "To @" + client.ID + ": " + string(p)[strings.Index(string(p), ":")+1:]}); err != nil {
						fmt.Println(err)
						return
					}

					if err := client.Conn.WriteJSON(Message{Type: messageType, Body: "From @" + c.ID + ": " + string(p)[strings.Index(string(p), ":")+1:]}); err != nil {
						fmt.Println(err)
						return
					}
				}

			}

		} else {
			message := Message{Type: messageType, Body: "@BOT: This command must have a : in it and/or be longer"}
			c.singleMes(message)
		}
		//returns all commands and command formats
	} else if string(p[1:]) == "help" {
		mList := make([]Message, 0)
		message := Message{Type: messageType, Body: "/make:chatroomName -makes a chatroom and make you an OP | /join#chatroomName -adds you to the chosen chatroom | /leave#chatroomName -removes you from the chosen chatroom"}
		message1 := Message{Type: messageType, Body: "/dc_userName:msg -send msg just to the chosen user | /chat#chatroomName:msg -send msg to just the chosen chatroom | /time -gives you the current time | /admin:password -makes you an admin"}
		message2 := Message{Type: messageType, Body: "Must be an OP to run: /chlim#chatroomName:limit -sets chatroom max users to the limit | /setop#chatroomName:userName -sets user to OP in chosen chatroom | /kicku#chatroomName:userName -removes user from chatroom"}
		message3 := Message{Type: messageType, Body: "Must be an OP to run: /muteu#chatroomName:userName -user can no longer chat in chatroom | /umute#chatroomName:userName -user is unmuted | /banur#chatroomName:userName -removes user from chatroom and stops them from joining | /unban#chatroomName:userName -user is unbanned"}
		mList = append(mList, message, message1, message2, message3)
		for _, x := range mList {
			c.singleMes(x)
		}

		//makes a new chatroom
	} else if string(p[1:5]) == "make" {
		if string(p[6:]) != "" {
			dup := false
			//checks if the chatroom name has already been made
			for _, x := range c.Bot.AllPool {
				if x.Name == string(p[6:]) {
					dup = true
				}
			}
			if dup {
				message := Message{Type: messageType, Body: "@BOT: Please input a unique chatroom name"}
				c.singleMes(message)
				//makes new chatroom makes the maker an OP
			} else {
				pool := NewPool()
				go pool.Start()
				pool.Name = string(p[6:])
				pool.OpList = append(pool.OpList, c)
				c.Pool = append(c.Pool, pool)

				c.Bot.AllPool = append(c.Bot.AllPool, pool)
				if len(c.Bot.AllAdmin) > 0 {
					c.opAdmin()
				}

				pool.Register <- c
				message := Message{Type: messageType, Body: "@BOT: You have made the chatroom " + pool.Name + " Your are the OP of this chatroom. Use /chat#chatname:msg to send messages just to this chat."}
				c.singleMes(message)
				c.Read()
			}

		} else {

			message := Message{Type: messageType, Body: "@BOT: Please input a unique chatroom name"}
			c.singleMes(message)
		}
		//removes the user from the chatroom
	} else if string(p[1:6]) == "leave" {
		if string(p[7:]) != "" {
			if len(c.Pool) > 1 {
				for i, x := range c.Pool {
					if x.Name == string(p[7:]) {
						x.Unregister <- c
						c.Conn.WriteJSON(Message{Type: 1, Body: "You have left chatroom #" + x.Name})
						c.Pool = append(c.Pool[:i], c.Pool[i+1:]...)
					}
				}
			} else {
				message := Message{Type: messageType, Body: "@BOT: You can not leave you last chatroom"}
				c.singleMes(message)
			}

		} else {
			message := Message{Type: messageType, Body: "@BOT: Please input a chatroom name"}
			c.singleMes(message)
		}
		//adds user to the chatroom
	} else if string(p[1:5]) == "join" {
		join := true

		for _, x := range c.Bot.AllPool {
			if x.Name == string(p[6:]) {
				join = false
				var admin bool
				for _, a := range c.Bot.AllAdmin {
					if c == a {
						admin = true
					}
				}
				//checks if the chatroom is at the limit or if the user is an admin
				if x.PoolSize > len(x.Clients) || admin {
					fmt.Println(x.PoolSize)
					var banned bool
					//checks if the user is not banned
					for _, z := range x.BanList {
						if z == c {
							banned = true
						}
					}
					if !banned {
						c.Pool = append(c.Pool, x)
						x.Register <- c
						c.Read()
						break
					} else {
						message := Message{Type: messageType, Body: "@BOT: Please input a chatroom name that you are not banned in"}
						c.singleMes(message)
					}
				} else {
					message := Message{Type: messageType, Body: "@BOT: Chatroom #" + x.Name + " is full"}
					c.singleMes(message)
				}
			}
		}
		if join {
			message := Message{Type: messageType, Body: "@BOT: Please input a chatroom name"}
			c.singleMes(message)
		}
		//sends a msg to just one chat
	} else if string(p[1:5]) == "chat" {
		if string(p[6:]) != "" {
			var send bool
			var has bool
			for _, b := range string(p[5:]) {
				if b == ':' {
					has = true
				}

			}
			if has {
				send = true
				for _, x := range c.Pool {
					//check for chatroom name
					if x.Name == string(p[6:strings.Index(string(p), ":")]) {
						send = false
						var isMuted bool
						//checks if user is muted
						for _, x := range x.MuteList {
							if x == c {
								isMuted = true
							}
						}
						if !isMuted {
							message := Message{Type: messageType, Body: "#" + x.Name + ": @" + c.ID + ": " + string(p)[strings.Index(string(p), ":")+1:]}
							x.Broadcast <- message
							fmt.Printf("Message Received: %+v\n", message)
						} else {
							message := Message{Type: messageType, Body: "@BOT: You are muted in chatroom #" + x.Name}
							c.singleMes(message)
						}

					}
				}
			} else {
				message := Message{Type: messageType, Body: "@BOT: This command must have a : in it and/or be longer"}
				c.singleMes(message)
			}
			if send {
				message := Message{Type: messageType, Body: "@BOT: Please input a chatroom name that you are in."}
				c.singleMes(message)
			}
		} else {
			message := Message{Type: messageType, Body: "@BOT: Please input a chatroom name"}
			c.singleMes(message)
		}
		//set a chatroom limit
	} else if string(p[1:6]) == "chlim" {
		var isPool bool
		if string(p[7:]) != "" {
			var has bool
			for _, b := range string(p[6:]) {
				if b == ':' {
					has = true
				}

			}
			if has {
				isPool = true
				for _, x := range c.Pool {
					//check for chatroom name
					if x.Name == string(p[7:strings.Index(string(p), ":")]) {
						isPool = false
						isOp := false
						//check if the user is an OP of the chatroom
						for _, i := range x.OpList {
							if c == i {
								isOp = true
							}
						}
						if isOp {
							size, _ := strconv.Atoi(string(p)[strings.Index(string(p), ":")+1:])
							x.PoolSize = size
							message := Message{Type: messageType, Body: "@BOT: Chatroom #" + x.Name + "now has a limit of " + fmt.Sprint(x.PoolSize)}
							x.Broadcast <- message
							break
						}
					}
				}
			} else {
				message := Message{Type: messageType, Body: "@BOT: This command must have a : in it"}
				c.singleMes(message)
			}
			if isPool {
				message := Message{Type: messageType, Body: "@BOT: Please input a chatroom"}
				c.singleMes(message)
			}
		}
		//make the user an admin
	} else if string(p[1:6]) == "admin" {
		if string(p[7:]) != "" {
			var has bool
			for _, b := range string(p[6:]) {
				if b == ':' {
					has = true
				}
			}
			if has {
				//checks if the user entered the right password
				if string(p[strings.Index(string(p), ":")+1:]) == "password" {
					c.Bot.AllAdmin = append(c.Bot.AllAdmin, c)
					c.op1Admin(c)
					message := Message{Type: messageType, Body: "@BOT: You are now an Admin, you are an OP in all chatrooms and can join full chatrooms"}
					c.singleMes(message)
				}
			} else {
				message := Message{Type: messageType, Body: "@BOT: This command must have a : in it and you must enter the password"}
				c.singleMes(message)
			}
		}
		//sets another user to an OP of the chatroom
	} else if string(p[1:6]) == "setop" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		//sets user to an OP
		if a && b && y {
			x.OpList = append(x.OpList, z)
			message := Message{Type: messageType, Body: "@BOT: " + z.ID + " is now an OP of chatroom #" + x.Name}
			x.Broadcast <- message
		}
		//mutes a user
	} else if string(p[1:6]) == "muteu" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		//mutes a user
		if a && b && y {
			x.MuteList = append(x.MuteList, z)
			message := Message{Type: messageType, Body: "@BOT: " + z.ID + " is now muted in chatroom #" + x.Name}
			x.Broadcast <- message
		}
		//unmutes a muted user
	} else if string(p[1:6]) == "umute" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		var index int
		if a && b && y {
			var isMute bool
			//finds the index of the muted user
			for i, o := range x.MuteList {
				if o == z {
					isMute = true
					index = i
				}
			}
			if isMute {
				x.MuteList = append(x.MuteList[:index], x.MuteList[index+1:]...)
				message := Message{Type: messageType, Body: "@BOT: " + z.ID + " is now unmuted in chatroom #" + x.Name}
				x.Broadcast <- message
			} else {
				message := Message{Type: messageType, Body: "@BOT: User was not muted in chatroom #" + x.Name}
				c.singleMes(message)
			}

		}
	} else if string(p[1:6]) == "unban" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		var index int
		if a && b && y {
			var isBan bool
			//finds the index of the banned user
			for i, o := range x.BanList {
				if o == z {
					isBan = true
					index = i
				}
			}
			if isBan {
				x.BanList = append(x.BanList[:index], x.BanList[index+1:]...)
				message1 := Message{Type: messageType, Body: "@BOT: You have been unbanned from chatroom #" + x.Name}
				z.singleMes(message1)
				message := Message{Type: messageType, Body: "@BOT: " + z.ID + " is now unbanned in chatroom #" + x.Name}
				x.Broadcast <- message
			} else {
				message := Message{Type: messageType, Body: "@BOT: User is not banned in chatroom #" + x.Name}
				c.singleMes(message)
			}

		}
	} else if string(p[1:6]) == "banur" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		var index int
		if a && b && y {
			x.BanList = append(x.BanList, z)
			var inPool bool
			//finds the chatroom's index the user is being banned in
			for i, o := range z.Pool {
				if o == x {
					inPool = true
					index = i
				}
			}
			if inPool {
				//removes the user from the chatroom
				for r := range x.Clients {
					if r == z {
						x.Unregister <- z
						z.Pool = append(c.Pool[:index], c.Pool[index+1:]...)
					}
				}
			}

			message1 := Message{Type: messageType, Body: "@BOT: You have been banned from chatroom #" + x.Name}
			z.singleMes(message1)
			message := Message{Type: messageType, Body: "@BOT: " + z.ID + " is now banned in chatroom #" + x.Name}
			x.Broadcast <- message
		}
		//removes the user from the chatroom
	} else if string(p[1:6]) == "kicku" {
		//checks if the user is an OP, checks if the chatroom exist, checks if the target user exist
		a, b, y, z, x := c.opCommand(messageType, p)
		var index int
		if a && b && y {
			var inPool bool
			//finds the chatroom's index the user is being removed from
			for i, o := range z.Pool {
				if o == x {
					inPool = true
					index = i
				}
			}
			if inPool {
				x.Unregister <- z
				z.Pool = append(c.Pool[:index], c.Pool[index+1:]...)
				message1 := Message{Type: messageType, Body: "@BOT: You have been kicked from chatroom #" + x.Name}
				z.singleMes(message1)
				message := Message{Type: messageType, Body: "@BOT: " + z.ID + " has been kicked chatroom #" + x.Name}
				x.Broadcast <- message
			} else {
				message := Message{Type: messageType, Body: "@BOT: User is not in chatroom #" + x.Name}
				c.singleMes(message)
			}

		}
	} else {
		message := Message{Type: messageType, Body: "@BOT: Invalied command, type /help to see all commands"}
		c.singleMes(message)
	}
}

//adds all admins to the OP list of all chatrooms
func (c *Client) opAdmin() {
	for _, y := range c.Bot.AllAdmin {
		for _, x := range c.Bot.AllPool {
			addOp := true
			for _, z := range x.OpList {
				if z == y {
					addOp = false
					break
				}
			}
			if addOp {
				x.OpList = append(x.OpList, y)
			}
		}
	}
}

//adds one admin to the OP list of all chatrooms
func (c *Client) op1Admin(a *Client) {
	for _, x := range c.Bot.AllPool {
		addOp := true
		for _, z := range x.OpList {
			if z == a {
				addOp = false
				break
			}
		}
		if addOp {
			x.OpList = append(x.OpList, a)
		}
	}
}
