package websocket

import "strings"

func (c *Client) opCommand(messageType int, p []byte) (bool, bool, bool, *Client, *Pool) {
	var isIn, isOp, isPool = true, true, false
	var client *Client
	var pool *Pool

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

				if x.Name == string(p[7:strings.Index(string(p), ":")]) {

					pool = x
					isPool = false
					isOp = true
					for _, i := range x.OpList {
						if c == i {

							isOp = false
							isIn = true
							for _, z := range c.Bot.AllClients {
								if string(p)[strings.Index(string(p), ":")+1:] == z.ID {

									client = z
									isIn = false
									break
								}
							}
							if isIn {
								message := Message{Type: messageType, Body: "@BOT: Please input a real user"}
								c.singleMes(message)
							}
						}
					}
					if isOp {
						message := Message{Type: messageType, Body: "@BOT: You are not an OP of this chatroom"}
						c.singleMes(message)
						break
					}
				}
			}
			if isPool {
				message := Message{Type: messageType, Body: "@BOT: Please input a chatroom"}
				c.singleMes(message)
			}
		}
	} else {
		message := Message{Type: messageType, Body: "@BOT: This command must have a : in it and/or be longer"}
		c.singleMes(message)
	}
	return !isPool, !isOp, !isIn, client, pool
}
