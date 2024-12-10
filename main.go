package main

import (
	"fmt"

	"src.tiger.place/statsapi/internal/node"
	"src.tiger.place/statsapi/internal/supervisor"
)

func main() {
	fmt.Println("Hi")
	go node.StartNode()
	supervisor.Start()
}
