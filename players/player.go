package players

import (
	"../connection"
	"./actions"
	"fmt"
	"math/rand"
	"time"
	"../servers/serverrpc"
)

type Action actions.Action

type Player struct {
	id int
	uniqueIdentifier int
	x int
	y int
	direction int
	alive bool

	actionList [] Action
	rttLogs [] int64

	tick int
	connection connection.Connection
}

func createPlayerRequest(player Player) serverrpc.PlayerRequest {
	request := serverrpc.PlayerRequest{
		Id:               player.id,
		UniqueIdentifier: player.uniqueIdentifier,
		X:                player.x,
		Y:                player.y,
		Direction:        player.direction,
		Alive:            player.alive,
		//ActionList:       player.actionList,
		Tick:             player.tick,
	}
	return request
}

func RTT(before time.Time, after time.Time) int64 {
	return (after.UnixNano()-before.UnixNano())/int64(time.Millisecond)
}

func (player *Player) Run() {
	for i := 0; i < 100; i++ {
		response := serverrpc.ServerResponse{}
		request := createPlayerRequest(*player)
		before := time.Now()
		player.connection.Call("TServer.UpdatePlayerState", request, &response)
		after := time.Now()
		player.rttLogs = append(player.rttLogs, RTT(before, after))
		time.Sleep(16*time.Millisecond)
	}
	time.Sleep(5*time.Second)
	go fmt.Println(player.rttLogs)
}

func (player *Player) JoinGame(connection *connection.Connection) {
	player.connection = *connection
	player.Run()
}

func CreatePlayer(id int) *Player {
	rand.Seed(time.Now().UTC().UnixNano())

	player := &Player{}
	player.id = id
	player.uniqueIdentifier = rand.Int()
	player.x = rand.Intn(200 - 0)
	player.y = rand.Intn(200 - 0)
	player.actionList = []Action{}
	player.rttLogs = []int64{}
	player.direction = actions.GetRandomDirection()
	player.tick = 0
	player.alive = true

	return player
}

