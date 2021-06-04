package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// Client 客户端信息
type Client struct {
	// 当前 Websocket 中心
	hub *Hub

	// 客户端的连接
	conn *websocket.Conn

	// 客户端的 id
	id int64

	// 消息管道，由 Hub 分发
	messages chan *Message
}

func (c *Client) checkIn() {
	c.hub.checkIn <- c
}

func (c *Client) checkOut() {
	c.hub.checkOut <- c
}

// 客户端连接上后，将 id 发给客户端
func (c *Client) issueId()  {
	message := &Message{
		To: c.id,
		From: -1,
		Content: "issue id",
	}

	c.messages <- message
}

// read 用来接收客户端发送过来的消息，消息为 json 格式
// 将接收到的消息解码后，发布到 redis 里
func (c *Client) read() {
	defer func() {
		c.checkOut()
		_ = c.conn.Close()
	}()

	for {
		message := &Message{}
		err := c.conn.ReadJSON(message)
		if err != nil {
			break
		}

		if message.To == c.id {
			continue
		}

		message.From = c.id
		err = publish(message)
		if err != nil {
			break
		}
	}
}

// write 用来发送消息给客户端，消息为 json 格式
// 从 messages 管道里获取消息，专属于本客户端的消息
func (c *Client) write() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message := <- c.messages:
			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <- ticker.C:
			message := &Message{
				To: c.id,
				From: -1,
				Content: "PING",
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}
