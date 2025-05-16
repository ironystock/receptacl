package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type server struct {
	msgBuffer   int
	mux         http.ServeMux
	clients     map[*client]struct{}
	clientMutex sync.Mutex
}
type client struct {
	messages chan []byte
}

func NewServer() *server {
	s := &server{
		msgBuffer: 10,
		clients:   make(map[*client]struct{}),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./static")))
	s.mux.HandleFunc("/ws", s.clientHandle)
	return s
}
func (s *server) clientHandle(w http.ResponseWriter, r *http.Request) {
	err := s.clientReg(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (s *server) addClient(client *client) {
	s.clientMutex.Lock()
	s.clients[client] = struct{}{}
	s.clientMutex.Unlock()
}
func (s *server) clientReg(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var wc *websocket.Conn
	client := &client{
		messages: make(chan []byte, s.msgBuffer),
	}
	s.addClient(client)
	wc, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer wc.CloseNow()

	c = wc.CloseRead(c)
	for {
		select {
		case msg := <-client.messages:
			c, cancel := context.WithTimeout(c, time.Second)
			defer cancel()
			err := wc.Write(c, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-c.Done():
			return c.Err()
		}

	}
}
