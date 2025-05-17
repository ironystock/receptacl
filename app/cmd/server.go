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
	cfg         receptaclCfg
	mux         http.ServeMux
	clients     map[*client]struct{}
	clientMutex sync.Mutex
}
type client struct {
	messages chan []byte
}

func NewServer(cfg receptaclCfg) *server {
	s := &server{
		clients: make(map[*client]struct{}),
		cfg:     cfg,
	}
	//TODO: remove static for module
	s.mux.Handle("/", http.FileServer(http.Dir("./static")))
	s.mux.HandleFunc(s.cfg.LocalPath, s.clientHandle)
	return s
}
func (s *server) clientHandle(w http.ResponseWriter, r *http.Request) {
	err := s.clientReg(r.Context(), w, r)
	if err != nil {
		return
	}
}
func (s *server) addClient(client *client) {
	s.clientMutex.Lock()
	s.clients[client] = struct{}{}
	s.clientMutex.Unlock()
}
func (s *server) delClient(client *client) {
	s.clientMutex.Lock()
	delete(s.clients, client)
	s.clientMutex.Unlock()
}
func (s *server) clientReg(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var wc *websocket.Conn
	client := &client{
		messages: make(chan []byte, s.cfg.MsgBuffer),
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
			c, cancel := context.WithTimeout(c, time.Duration(s.cfg.Timeout)*time.Millisecond)
			defer cancel()
			err := wc.Write(c, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-c.Done():
			s.delClient(client)
			wc.CloseNow()
			return c.Err()
		}

	}
}
func (s *server) broadcast(msg []byte) {
	s.clientMutex.Lock()
	for client := range s.clients {
		fmt.Println(s.clients)
		client.messages <- msg
	}
	s.clientMutex.Unlock()
}
