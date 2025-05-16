package main

import (
	"fmt"
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
		fmt.Println("config issue")
		panic("No config")
	}
	fmt.Println(cfg)

	// go func() {
	// 	//Websocket to proxy
	// }
}
