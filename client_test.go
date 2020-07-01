package project

import (
	"./connection"
	"./players"
	"./servers/traditional"
	"fmt"
	"testing"
	"time"
)

type ServerInfo struct {
	protocol string
	address string
	port string
}
var serverInfo = ServerInfo{
	protocol: "tcp",
	address: "127.0.0.1",
	port: "8000",
}

type playerStatistic struct {
	playerId int
	maxRtt int
	lossRate int
}

func tickToTime(tickRate int) float32 {
	return float32(1000.0 / tickRate)
}

func displayPlayerStatistics(playerList []*players.Player) {
	var playerStats []playerStatistic
	for _, player := range playerList {
		id, maxRtt, lossRate := player.GetNetworkStats()
		playerStats = append(playerStats, playerStatistic{
			playerId: id,
			maxRtt: maxRtt,
			lossRate: lossRate,
		})
	}
	fmt.Printf("Results: %+v\n", playerStats)
	maxRtt := -1
	maxLoss := -1
	for _, playerStatistic := range playerStats {
		if playerStatistic.maxRtt > maxRtt {maxRtt = playerStatistic.maxRtt}
		if playerStatistic.lossRate > maxLoss {maxLoss = playerStatistic.lossRate}
	}
	fmt.Printf("MaxRTT: %d - MaxLoss: %d\n", maxRtt, maxLoss)
}

func TestGame(t *testing.T) {
	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	gameServer := traditional.StartServer(*conn)
	game := gameServer.Game
	tickTime := int(tickToTime(71))
	time.Sleep(1*time.Second)
	fmt.Println(gameServer)
	var playerList []*players.Player

	for i := 0; i < 100; i++ {
		player := players.CreatePlayer(i)
		playerList = append(playerList, player)
		go player.JoinGame(conn, tickTime)
		time.Sleep(1*time.Millisecond)
	}

	for !game.IsFinished() {
		time.Sleep(5*time.Second)
	}

	displayPlayerStatistics(playerList)
}


