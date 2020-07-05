package main

import (
	"../connection"
	"../servers/proposed"
	"strconv"
	"time"
)

func main() {
	artificialDelay := 1
	var workerList []proposed.Worker
	for i := 0; i < 10; i++ {
		conn := connection.CreateConnection("udp", "127.0.0.1", "801" + strconv.Itoa(i))
		worker := proposed.StartWorker(*conn, artificialDelay)
		workerList = append(workerList, *worker)
	}
	time.Sleep(20*time.Minute)
}
