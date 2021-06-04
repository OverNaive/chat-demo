# chat demo

## 如何运行

### 构建镜像

切换至 demo 项目的根目录，执行命令： `docker build --no-cache -t chat:demo .`

### docker-compose

切换至 demo 项目的根目录，执行命令： `docker-compose up -d`

### 聊天页面

浏览器打开 `http://127.0.0.1:9999/chat.html` 即可进行聊天

## 如何实现

### id 生成

利用 Redis 的 incr 生成递增 id,。incr 具有原子性，可保证 id 唯一

### 跨服务消息

采用 Redis 的 pub/sub 方式，Websocket 服务收到消息后立马发布到指定 channel。

每个 Websocket 服务都订阅了该 channel，获取信息后进行判断，在发送给客户端。

### 消息包

采用 json 的格式，具体如下：

```json
{
  "to": 0,
  "from": 0,
  "content": "Welcome"
}
```