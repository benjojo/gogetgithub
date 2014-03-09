package main

import (
	"flag"
	"fmt"
	// "net/http"
	"os"
)

func main() {
	Username := flag.String("username", "", "The username that you are targetting")
	flag.Parse()

	if os.Getenv("GOPATH") == "" {
		fmt.Println("You don't have a GOPATH set, I don't know where to clone to! Please set one.")
		os.Exit(1)
	}

	if *Username == "" {
		fmt.Println("Please give a username to auto clone all their go stars to the gopath to.")
		os.Exit(1)
	}
}
