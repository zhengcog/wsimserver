package chat

import (
	"fmt"
	"io"
	"log"
	"wsimserver/models"
	"wsimserver/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	id     int
	userid string
	ws     *websocket.Conn
	server *Server
	ch     chan Message
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if server == nil {
		panic("server cannot be nil")
	}
	maxId++
	ch := make(chan Message, channelBufSize)
	doneCh := make(chan bool)
	return &Client{maxId, "", ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg, c.userid)
			_, err := c.ws.Write(msg)
			if err != nil {
				c.server.Del(c)
				c.doneCh <- true // for listenRead method
				log.Println("Error:", err.Error())
				return
			}
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message = make(Message, 1024)
			i, err := c.ws.Read(msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
			} else {
				msg = msg[:i]
				input := ParseMessage(msg)
				log.Printf("Receive: %s\n", msg[:])
				if input == nil { //登录 websocket服务器 token=Bearer agag.....
					action, login, user := WhetherLogin(msg)
					log.Println(action, login, user)
					if action {
						if login {
							c.userid = user
							c.server.Online(c)
							c.Write([]byte("ok"))
							messages := models.SendoutOfflineMsg(user)
							for _, v := range messages {
								c.Write([]byte(v)) //offline message
							}
						} else {
							c.server.Del(c)
						}
					}
				} else { //send message to client
					lockUsers.RLock()
					touser, online := c.server.users[input.Target]
					lockUsers.RUnlock()
					output := NewOutput(input)
					if online {
						touser.Write([]byte(output.String()))
					} else {
						//offline....
						models.PushOfflineMsg(input.Target, output.String())
					}

				}
			}
		}
	}
}
