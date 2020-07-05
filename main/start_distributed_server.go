package main

import (
	"../connection"
	"../servers/proposed"
	"fmt"
	"strconv"
	"time"
)

func main() {

	artificialDelay := 1
	var workerList []proposed.Worker
	for i := 0; i < 10; i++ {
		worker := &proposed.Worker{
			Address: "192.168.0.18", // NOTE: place worker server address here
			Port:    "801" + strconv.Itoa(i),
		}
		workerList = append(workerList, *worker)
	}
	workerPool := proposed.CreateWorkerPool(workerList)

	conn := connection.CreateConnection("udp", "0.0.0.0", "8000")
	gameServer := proposed.StartServer(*conn, *workerPool, artificialDelay)
	game := gameServer.Game

	for !game.IsFinished() {
		time.Sleep(5*time.Second)
	}
	fmt.Println("Finished game")
}
