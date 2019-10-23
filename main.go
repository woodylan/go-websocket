package main

import (
	"go-websocket/src"
)

func main() {
	server := src.NewServer(":666")

	if err := server.ListenAndServer(); err != nil {
		panic(err)
	}
}
