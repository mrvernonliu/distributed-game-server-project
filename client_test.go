package project

import (
	"./players"
	"fmt"
	"testing"
	"./servers/traditional"
	"./connection"
	"time"
)

type ServerInfo struct {
	protocol string
	address string
	port string
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

func TestGame_Internal_Traditional(t *testing.T) {
	var serverInfo = ServerInfo{
		protocol: "udp",
		address: "127.0.0.1",
		port: "8000",
	}
	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	artificialDelay := 1
	gameServer := traditional.StartServer(*conn, artificialDelay)
	game := gameServer.Game
	tickTime := int(tickToTime(60))
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

func TestGame_External_Traditional(t *testing.T) {
	var serverInfo = ServerInfo{
		protocol: "udp",
		address: "10.0.0.55",
		port: "8000",
	}
	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	tickTime := int(tickToTime(20))

	var playerList []*players.Player

	for i := 0; i < 100; i++ {
		player := players.CreatePlayer(i)
		playerList = append(playerList, player)
		go player.JoinGame(conn, tickTime)
		time.Sleep(1*time.Millisecond)
	}

	time.Sleep(30*time.Second)
	displayPlayerStatistics(playerList)
}

