package main

import (
	"fmt"
	"net/http"

	"github.com/TutorialEdge/realtime-chat-go-react/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request, bot *websocket.Bot) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: make([]*websocket.Pool, 0),
		Bot:  bot,
	}
	client.Pool = append(client.Pool, pool)
	client.Bot.AllClients = append(client.Bot.AllClients, client)
	pool.Register <- client
	client.Read()
}
func makebot(pool *websocket.Pool) *websocket.Bot {

	bot := &websocket.Bot{
		ID:         "BOT",
		AllPool:    make([]*websocket.Pool, 0),
		AllClients: make([]*websocket.Client, 0),
		AllAdmin:   make([]*websocket.Client, 0),
	}
	bot.AllPool = append(bot.AllPool, pool)

	return bot
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()
	makeBot := true
	var bot *websocket.Bot
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		if makeBot {
			bot = makebot(pool)
		}
		makeBot = false
		serveWs(pool, w, r, bot)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()
	http.ListenAndServe("127.0.0.1:8080", nil)
}
