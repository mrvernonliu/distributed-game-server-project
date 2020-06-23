package players

import (
	"../connection"
	"./actions"
	"fmt"
	"math/rand"
	"time"
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

	connection connection.Connection
}

func (player *Player) Run() {
	for i := 0; i < 100; i++ {
		go fmt.Printf("Running %d\n", player.id)
	}
}

func (player *Player) JoinGame(connection *connection.Connection) {
	player.connection = *connection
	go player.Run()
}

func CreatePlayer(id int) *Player {
	rand.Seed(time.Now().UTC().UnixNano())

	player := &Player{}
	player.id = id
	player.uniqueIdentifier = rand.Int()
	player.x = rand.Intn(200 - 0)
	player.y = rand.Intn(200 - 0)
	player.direction = actions.GetRandomDirection()
	player.alive = true

	return player
}

