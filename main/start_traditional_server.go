package main

import (
	"../connection"
	"../servers/traditional"
	"fmt"
	"time"
)



func main() {
	conn := connection.CreateConnection("udp", "0.0.0.0", "8000")
	fmt.Printf("Starting server on %+v\n", conn)
	artificialDelay := 200000
	gameServer := traditional.StartServer(*conn, artificialDelay)
	game := gameServer.Game
	for !game.IsFinished() {
		time.Sleep(5*time.Second)
	}
	println("Finished game")
}