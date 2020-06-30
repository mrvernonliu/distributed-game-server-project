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
	// Player variables
	id int
	uniqueIdentifier int
	x int
	y int
	direction int
	alive bool

	actionList [] Action

	// Connection variables
	rttLogs [] int64
	tick       int
	tickTime   int
	lostPackets int
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

func (player *Player) callServer() {
	before := time.Now()
	response := serverrpc.ServerResponse{}
	request := createPlayerRequest(*player)
	player.connection.Call("TServer.UpdatePlayerState", request, &response)
	after := time.Now()
	player.rttLogs = append(player.rttLogs, RTT(before, after))
	if response.Tick != player.tick {player.lostPackets++}
	//fmt.Printf("%d %d\n",response.Tick, player.tick)
}

func (player *Player) Run() {
	sleepTime := time.Duration(player.tickTime)* time.Millisecond
	for i := 0; i < 100; i++ {
		player.tick++
		go player.callServer()
		time.Sleep(sleepTime)
	}
	time.Sleep(5*time.Second)
	fmt.Println(player.rttLogs)
	fmt.Printf("Loss: %d\n", player.lostPackets)
}

func (player *Player) JoinGame(connection *connection.Connection, refreshRate int) {
	player.connection = *connection
	player.tickTime = refreshRate
	go fmt.Printf("%d - Connected to %x with tick time %d\n", player.id, player.connection, player.tickTime)
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
	player.lostPackets = 0
	player.alive = true

	return player
}


/*
	Random helper methods
 */

func RTT(before time.Time, after time.Time) int64 {
	return (after.UnixNano()-before.UnixNano())/int64(time.Millisecond)
}

