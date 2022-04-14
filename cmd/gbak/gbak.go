package main

import (
	"fmt"
	"g.bak/internal/cmd/credentials"
	"g.bak/internal/cmd/server"
	"g.bak/internal/cmd/sync"
	"g.bak/internal/cmd/tokens"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please specify a command")
		return
	}

	cmd := os.Args[1]
	os.Args = os.Args[2:]

	switch cmd {
	case "serve":
		server.Server()
	case "credentials":
		credentials.Credentials()
	case "tokens":
		tokens.Tokens()
	case "sync":
		sync.Sync()
	default:
		fmt.Printf("'%s' is not a valid command", cmd)
	}
}