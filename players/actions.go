package players

import (
	"math/rand"
	"time"
	"./actioninterfaces"
	"../game/gameinterfaces"
)


const (
	Shoot actioninterfaces.Action = 1
	Aim actioninterfaces.Action = 2
	Turn actioninterfaces.Action = 3
	Move actioninterfaces.Action = 4
)

const (
	North int = 1
	East int = 2
	West int = 3
	South int = 4
)

func GetRandomDirection() int {
	newDirection := rand.Intn(4-1) + 1
	switch newDirection {
	case 1:
		return North
	case 2:
		return East
	case 3:
		return West
	case 4:
		return South
	default:
		return North
	}
}

func GetRandomPlayer(id int, players [] gameinterfaces.InGamePlayer) int{
	randomPlayer := id
	if len(players) == 0 {
		return -1
	}
	for randomPlayer == id || !players[randomPlayer].Alive {
		randomPlayer++
		if randomPlayer >= 100 {
			randomPlayer = 0
		}
	}
	return randomPlayer
}

func GetRandomAction(id int, players [] gameinterfaces.InGamePlayer) actioninterfaces.ActionUpdate {
	rand.Seed(time.Now().UTC().UnixNano())
	action := rand.Intn(5 - 0)
	if action == 0 { return actioninterfaces.ActionUpdate{Action: Turn, Payload: GetRandomDirection()}}
	if action == 1 { return actioninterfaces.ActionUpdate{ Action: Shoot, Payload: GetRandomPlayer(id, players)}}
	return actioninterfaces.ActionUpdate{Action: -1, Payload: -1}
}
