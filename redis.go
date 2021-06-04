package main

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

// initRedisPool 初始化连接池
func initRedisPool() {
	redisPool = &redis.Pool{
		MaxIdle: 30,
		MaxActive: 30,
		IdleTimeout: 240 *time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// subscribe 订阅 channel: "chat:message"
// 将获取到消息，解码后传输到 Hub.messages 管道里
func subscribe(hub *Hub) {
	for {
		conn := redisPool.Get()

		psc := redis.PubSubConn{Conn: conn}
		_ = psc.Subscribe("chat:message")

		for conn.Err() == nil {
			switch reply := psc.Receive().(type) {
			case redis.Message:
				message := &Message{}
				err := json.Unmarshal(reply.Data, message)
				if err != nil {
					break
				}

				hub.messages <- message
			case redis.Subscription:
			case error:
				break
			}
		}

		_ = conn.Close()
	}
}

// publish 发布消息到 channel: "chat:message"
// 将消息 json 化再发布
func publish(message *Message) (err error) {
	conn := redisPool.Get()
	defer func() {
		_ = conn.Close()
	}()

	var msg []byte
	msg, err = json.Marshal(message)
	_, err = conn.Do("PUBLISH", "chat:message", msg)
	return
}

// getClientId 生成客户端 Client.id
// 利用 incr 的原子性，获取递增的数字来作为 id
func getClientId() (id int64, err error) {
	conn := redisPool.Get()
	defer func() {
		err = conn.Close()
	}()

	id, err = redis.Int64(conn.Do("INCR", "chat:client_id"))
	return
}


