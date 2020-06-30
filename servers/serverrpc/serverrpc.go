package serverrpc

import (
	"../../players/actions"
)

type Action actions.Action

type PlayerRequest struct {
	Id int
	UniqueIdentifier int
	X int
	Y int
	Direction int
	Alive bool

	ActionList [] Action

	Tick int
}

type ServerResponse struct {
	Id int
	UniqueIdentifier int
	X int
	Y int
	Direction int
	Alive bool

	Tick int
}

