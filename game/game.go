package game

import (
	"../servers/serverinterfaces"
	"./gameinterfaces"
	"fmt"
	"math/rand"
	"sync"
	"time"
)


type Game struct {
	Id int

	phase int
	artificialDelay time.Duration

	mux     sync.Mutex
	Players map[int]gameinterfaces.InGamePlayer
}


type PlayerRequest serverinterfaces.PlayerRequest

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

func (game *Game) UpdateGameState(request serverinterfaces.PlayerRequest) serverinterfaces.ServerResponse {
	response := populateServerResponse(request, game.phase)
	game.mux.Lock()
	if game.phase == 0 {
		if len(game.Players) < 100 {
			game.Players[response.Id] =  gameinterfaces.InGamePlayer{
				Id:               response.Id,
				UniqueIdentifier: response.UniqueIdentifier,
				X:                response.X,
				Y:                response.Y,
				Direction:        response.Direction,
				Alive:            true,
			}
			if len(game.Players) == 100 { game.phase = 1}
		}
	} else {
		if !game.Players[response.Id].Alive {
			response.Alive = false
		} else {
			actions := request.ActionList
			if len(actions) != 0 {
				for _, action := range actions {
					if action.Action == 3 {
						player := game.Players[response.Id]
						player.Direction = action.Payload
						game.Players[player.Id] = player
					} else if action.Action == 1 {
						// TODO: make look at direction
						targetPlayer := game.Players[action.Payload]
						targetPlayer.Alive = false
						game.Players[action.Payload] = targetPlayer
					}
				}
			}
			time.Sleep(game.artificialDelay) // Adding artificial delay to simulate action validation
		}
 	}
	var updatedPlayerList []gameinterfaces.InGamePlayer
 	for i := 0; i < 100; i++ {
 		updatedPlayerList = append(updatedPlayerList, game.Players[i])
	}
 	response.Players = updatedPlayerList
	game.mux.Unlock()

 	//fmt.Println(response)
	return response
}

func (game *Game) showState() {
	for {
		//go fmt.Println(game.Players)
		alive := -1
		gameState := ""
		for i := 0; i < 100; i++ {
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
		if alive > 0 {
			fmt.Printf("The winner is: %d\n", alive)
			game.phase = 3
		}
		time.Sleep(500*time.Millisecond)
	}
}

func (game *Game) IsFinished() bool {
	return game.phase == 3
}




func CreateGame(artificialDelay int) *Game {
	fmt.Println("Creating game")
	rand.Seed(time.Now().UTC().UnixNano())
	game := Game{}
	game.Id = rand.Int()
	game.Players = make(map[int]gameinterfaces.InGamePlayer)
	game.phase = 0
	game.mux = sync.Mutex{}
	game.artificialDelay = time.Duration(artificialDelay) * time.Nanosecond
	go game.showState()

	return &game
}
