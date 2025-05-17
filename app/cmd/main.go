package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// server-facing ws
type proxy struct {
	msgBuffer int
}

// client-facing ws

func main() {
	//192.168.1.171 / Triskaphile

	cfg, err := configure()
	if err != nil {
		panic("No config")
	}
	fmt.Println(cfg)
	s := NewServer(cfg)
	go func(s *server) {

		for {
			var tpl bytes.Buffer
			HxText("henlo", "catchtext", "").Render(context.Background(), &tpl)
			s.broadcast([]byte(tpl.String()))
			tpl.Reset()
			HxBlock("henlo", "catchblock", "").Render(context.Background(), &tpl)
			s.broadcast([]byte(tpl.String()))
			tpl.Reset()
			time.Sleep(1 * time.Second)
		}

	}(s)
	err = http.ListenAndServe(":"+s.cfg.LocalPort, &s.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// go func() {
	// 	//Websocket to proxy
	// }
}
