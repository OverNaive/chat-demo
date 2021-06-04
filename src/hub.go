package main

// Hub 当前 Websocket 的中心
type Hub struct {
	// 存储所有客户端，Client.id => *Client
	clients map[int64]*Client

	// 登记管道
	checkIn chan *Client

	// 登出管道
	checkOut chan *Client

	// 消息管道，从 redis 那里订阅，分发给客户端
	messages chan *Message
}

func (h *Hub) run() {
	for {
		select {
		case client := <- h.checkIn:
			h.clients[client.id] = client
		case client := <- h.checkOut:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.messages)
			}
		case message := <- h.messages:
			switch message.To {
			case 0:
				// 发送给所有人，除了自己
				for _, client := range h.clients {
					if message.From == client.id {
						continue
					}

					select {
					case client.messages <- message:
					default:
						close(client.messages)
						delete(h.clients, client.id)
					}
				}
			case -1:
				// nothing to do
			default:
				if client, ok := h.clients[message.To]; ok {
					// 发送给指定的人
					client.messages <- message
				}
			}
		}
	}
}
