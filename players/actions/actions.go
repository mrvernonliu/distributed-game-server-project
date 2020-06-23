package actions

import "math/rand"

type Action int

const (
	Shoot Action = 1
	Aim Action = 2
	Turn Action = 3
	Move Action = 4
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

func GetRandomAction() {
	// Needs to take game state and choose a random thing to do
}
