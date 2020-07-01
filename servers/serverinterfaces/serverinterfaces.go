package serverinterfaces

import (
	"../../players/actioninterfaces"
	"../../game/gameinterfaces"
)

type PlayerRequest struct {
	Id int
	UniqueIdentifier int
	X int
	Y int
	Direction int
	Alive bool

	ActionList [] actioninterfaces.ActionUpdate

	Tick int
}

type ServerResponse struct {
	Id int
	UniqueIdentifier int
	X int
	Y int
	Direction int
	Alive bool

	Players []gameinterfaces.InGamePlayer
	GamePhase int

	Tick int
}

