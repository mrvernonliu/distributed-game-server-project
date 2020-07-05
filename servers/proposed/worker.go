package proposed

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"time"
	"../../connection"
	"../../servers/serverinterfaces"
	"../../game/gameinterfaces"
)

type Worker struct {
	Id   int

	Conn    *net.UDPConn
	Dst     *net.UDPAddr
	Address string
	Port    string
}


type WorkerRequest struct {
	PlayerRequest   serverinterfaces.PlayerRequest
	GameState       map[int]gameinterfaces.InGamePlayer
	ArtificialDelay time.Duration
}

type WorkerResponse struct {
	Player serverinterfaces.PlayerRequest
	AffectedPlayers []gameinterfaces.InGamePlayer
}

func updatePlayers(workerRequest WorkerRequest) WorkerResponse {
	request := workerRequest.PlayerRequest
	gameState := workerRequest.GameState
	var affectedPlayers []gameinterfaces.InGamePlayer

	//go fmt.Printf("Worker got actions: %+v\n", workerRequest.PlayerRequest)

	if !gameState[request.Id].Alive {
		affectedPlayers = append(affectedPlayers, gameinterfaces.InGamePlayer{
			Id:               request.Id,
			UniqueIdentifier: request.UniqueIdentifier,
			Alive:            false,
		})
	} else {
		actions := request.ActionList
		if len(actions) != 0 {
			for _, action := range actions {
				if action.Action == 3 {
					player := gameState[request.Id]
					player.Direction = action.Payload
					affectedPlayers = append(affectedPlayers, player)
				} else if action.Action == 1 {
					// TODO: make look at direction
					targetPlayer := gameState[action.Payload]
					targetPlayer.Alive = false
					affectedPlayers = append(affectedPlayers, targetPlayer)
				}
			}
		}
	}
	time.Sleep(workerRequest.ArtificialDelay)

	return WorkerResponse{workerRequest.PlayerRequest, affectedPlayers}
}

func (worker *Worker) serve() {
	fmt.Printf("Worker %d listening\n", worker.Id)
	for {
		// Read from UDP
		recvBuf := make([]byte, 4096)
		n, client, _ := worker.Conn.ReadFromUDP(recvBuf[:])
		dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
		request := WorkerRequest{}
		dec.Decode(&request)
		//fmt.Printf("Worker %d received: %+v\n", worker.Id, request)
		response := updatePlayers(request)

		var sendBuf bytes.Buffer
		encoder := gob.NewEncoder(&sendBuf)
		encoder.Encode(response)
		worker.Conn.WriteToUDP(sendBuf.Bytes(), client)
	}
}


func StartWorker(connection connection.Connection, artificalDelay int) *Worker {
	rand.Seed(time.Now().UTC().UnixNano())
	worker := Worker{}
	worker.Id = rand.Int()
	worker.Address = connection.Address
	worker.Port = connection.Port
	worker.Dst, _ = net.ResolveUDPAddr("udp", ":"+connection.Port)
	worker.Conn, _ = net.ListenUDP("udp", worker.Dst)

	go worker.serve()
	return &worker
}