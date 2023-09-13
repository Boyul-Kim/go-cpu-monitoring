package main

import (
	"cpu-mon/server"
)

func main() {
	err := server.InitializeServer()
	if err != nil {
		panic(err)
	}
}
