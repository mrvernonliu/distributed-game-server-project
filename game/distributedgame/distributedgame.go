package distributedgame

import "time"

func doValidation() {
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
}