package main

import (
	"fmt"
	"tic-tac-toe-backend/internal/server"
)

func main() {

	server := server.NewServer()

	err := server.Run()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
