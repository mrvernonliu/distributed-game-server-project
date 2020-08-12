package proposedWithDistributor

import (
	"../../connection"
	"../../game/gameinterfaces"
	"../../servers/serverinterfaces"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

const LOBBY_SIZE = 20
const WORKER_COUNT = 3

type DistributedGame struct {
	Id int

	phase int
	artificialDelay time.Duration

	mux     sync.Mutex
	Players map[int]gameinterfaces.InGamePlayer

	workerPool WorkerPool
	distributor connection.Connection
	conn *net.UDPConn
	dst  *net.UDPAddr
}

func populateServerResponse(request serverinterfaces.PlayerRequest, phase int) serverinterfaces.ServerResponse {
	return serverinterfaces.ServerResponse{
		Id:               request.Id,
		UniqueIdentifier:  request.UniqueIdentifier,
		X:                request.X,
		Y:                request.Y,
		Direction:        request.Direction,
		Alive:            request.Alive,
		Players:          nil,
		GamePhase:		  phase,
		Tick:             request.Tick,
	}
}

func (game *DistributedGame) getWorkerPoolFromDistributor() {
	var workerList []WorkerAddress
	for i := 0; i < WORKER_COUNT; i++ {
		req := DistributorRequest{
			Request: "get",
		}
		var res DistributorResponse

		game.distributor.Call("Distributor.GetWorker", &req, &res)
		if res.Response == true {
			fmt.Printf("got worker %+v\n", res)
			workerList = append(workerList, WorkerAddress{
				Address: res.Address,
				Port:    res.Port,
			})
		} else {
			i--
		}
		game.workerPool = *CreateWorkerAddressPool(workerList)
		//fmt.Printf("got %d worker\n", i)
	}
}

func (game *DistributedGame) sendValidationToWorker(request serverinterfaces.PlayerRequest) serverinterfaces.ServerResponse {
	response := populateServerResponse(request, game.phase)
	//go fmt.Printf("Game got request: %+v\n", request)
	game.mux.Lock()
	if game.phase == 0 {
		if len(game.Players) < LOBBY_SIZE {
			game.Players[response.Id] = gameinterfaces.InGamePlayer{
				Id:               response.Id,
				UniqueIdentifier: response.UniqueIdentifier,
				X:                response.X,
				Y:                response.Y,
				Direction:        response.Direction,
				Alive:            true,
			}
			if len(game.Players) == LOBBY_SIZE {
				// Get Workers
				game.getWorkerPoolFromDistributor()
				game.phase = 1
				go game.listenToWorkers()
				go game.showState()
			}
		}
	} else {
		worker := game.workerPool.popIdle()
		address, _ := net.ResolveUDPAddr("udp", worker.Address + ":" + worker.Port)
		//fmt.Printf("Server sending to: %+v\n", address)
		workerRequest := WorkerRequest{
			PlayerRequest:   request,
			GameState:       game.Players,
			ArtificialDelay: game.artificialDelay,
		}
		//go fmt.Printf("Server - Formatted Worker request: %+v\n", workerRequest.PlayerRequest.ActionList)
		var sendBuf bytes.Buffer
		encoder := gob.NewEncoder(&sendBuf)
		encoder.Encode(workerRequest)
		game.conn.WriteToUDP(sendBuf.Bytes(), address)
		game.workerPool.pushIdle(worker) // Worker pool not properly implemented, just place at the end of queue
		var updatedPlayerList []gameinterfaces.InGamePlayer
		for i := 0; i < 100; i++ {
			updatedPlayerList = append(updatedPlayerList, game.Players[i])
		}
		response.Players = updatedPlayerList
	}
	if !game.Players[response.Id].Alive {
		response.Alive = false
	}
	game.mux.Unlock()
	return response
}

func (game *DistributedGame) updateState(request WorkerResponse) {
	game.mux.Lock()
	if game.Players[request.Player.Id].Alive {
		for _, affectedPlayer := range request.AffectedPlayers {
			//go fmt.Printf("%d Updating players %+v\n", request.Player.Id, request.AffectedPlayers)
			game.Players[affectedPlayer.Id] = affectedPlayer
		}
	}
	game.mux.Unlock()

}

func (game *DistributedGame) listenToWorkers() {
	for {
		recvBuf := make([]byte, 1024)
		n, _, _ := game.conn.ReadFromUDP(recvBuf[:])
		dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
		response := WorkerResponse{}
		dec.Decode(&response)
		//go fmt.Printf("Server got from worker: %+v\n", response)
		go game.updateState(response)
	}
}


func (game *DistributedGame) showState() {
	for {
		//go fmt.Println(game.Players)
		alive := -1
		gameState := ""
		for i := 0; i < LOBBY_SIZE; i++ {
			if i%10 == 0 {
				gameState += "\n"
			}
			if game.Players[i].Alive {
				gameState += " o "
				if alive == -1 {
					alive = i
				} else if alive > 0 {
					alive = -2
					continue
				}
			} else {
				gameState += " x "
			}
		}
		go fmt.Println(gameState)
		go fmt.Println(game.phase)
		if alive > 0 {
			fmt.Printf("The winner is: %d\n", alive)
			game.phase = 3
		}
		time.Sleep(500*time.Millisecond)
	}
}

func (game *DistributedGame) IsFinished() bool {
	return game.phase == 3
}

func (game *DistributedGame) ResetGame() {

}

func CreateGame(artificialDelay int, distributor connection.Connection) *DistributedGame {
	fmt.Println("Creating distributed game")
	rand.Seed(time.Now().UTC().UnixNano())
	game := DistributedGame{}
	game.Id = rand.Int()
	game.Players = make(map[int]gameinterfaces.InGamePlayer)
	game.phase = 0
	game.mux = sync.Mutex{}
	game.artificialDelay = time.Duration(artificialDelay) * time.Nanosecond
	game.distributor = distributor

	game.dst, _ = net.ResolveUDPAddr("udp", ":"+"8001")
	game.conn, _ = net.ListenUDP("udp", game.dst)
	go game.showState()

	return &game
}