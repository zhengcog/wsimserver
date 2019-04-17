package chat

import (
	"log"
	"net/http"
	"sync"
	"wsimserver/websocket"
)

var lockUsers sync.RWMutex

// Chat server.
type Server struct {
	pattern  string
	users    map[string]*Client
	addCh    chan *Client
	delCh    chan *Client
	onlineCh chan *Client
	errCh    chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	Users := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	onlineCh := make(chan *Client)
	errCh := make(chan error)
	return &Server{
		pattern,
		Users,
		addCh,
		delCh,
		onlineCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Online(c *Client) {
	s.onlineCh <- c
}
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Listening server...")
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Server{Handler: onConnected})
	log.Println("Created handler")
	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			//log.Println("Added new client")
			log.Println("Added new client Now", maxId, "clients.")
			_ = c
		// del a client
		case c := <-s.delCh:
			log.Println("Delete client", c.userid)
			lockUsers.Lock()
			delete(s.users, c.userid)
			lockUsers.Unlock()
		case c := <-s.onlineCh:
			lockUsers.Lock()
			s.users[c.userid] = c
			lockUsers.Unlock()
		case err := <-s.errCh:
			log.Println("Error:", err.Error())
			//case <-s.doneCh:
			//	return
		}
	}
}

func (s *Server) GetClient(user string, server *Server) (*Client, bool) {
	lockUsers.RLock()
	client, online := server.users[user]
	lockUsers.RUnlock()
	return client, online
}
