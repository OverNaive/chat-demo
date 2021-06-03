package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWs 处理 Websocket 连接
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	var id int64
	id, err = getClientId()
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &Client{
		hub: hub,
		conn: conn,
		id: id,
		messages: make(chan *Message, 256),
	}

	client.checkIn()
	client.issueId()

	go client.write()
	go client.read()
}

func init() {
	initRedisPool()
}

func main() {
	s := &http.Server{
		Addr: ":8888",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	hub := &Hub{
		clients: make(map[int64]*Client),
		checkIn: make(chan *Client),
		checkOut: make(chan *Client),
		messages: make(chan *Message),
	}
	go hub.run()
	go subscribe(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Fatal(s.ListenAndServe())
}
