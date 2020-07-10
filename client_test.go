package project

import (
	"./connection"
	"./players"
	"./servers/proposed"
	"./servers/traditional"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
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

func convertRttToString(rttlogs []int) string {
	var output [] string
	for _, i := range rttlogs {
		output = append(output, strconv.Itoa(i))
	}

	return strings.Join(output, ",")
}

func displayPlayerStatistics(playerList []*players.Player) {
	var playerStats []playerStatistic
	fileOutput := ""
	for _, player := range playerList {
		id, maxRtt, lossRate := player.GetNetworkStats()
		resultEntry := strconv.Itoa(id) + "," + convertRttToString(player.RttLogs) + "\n"
		fileOutput += resultEntry
		playerStats = append(playerStats, playerStatistic{
			playerId: id,
			maxRtt: maxRtt,
			lossRate: lossRate,
		})
	}
	fmt.Printf(fileOutput)
	ioutil.WriteFile("results.csv", []byte(fileOutput), 0644)
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

func TestGame_Internal_Distributed(t *testing.T) {
	var serverInfo = ServerInfo{
		protocol: "udp",
		address: "127.0.0.1",
		port: "8000",
	}
	artificialDelay := 1
	var workerList []proposed.Worker
	for i := 5; i < 10; i++ {
		conn := connection.CreateConnection("udp", "127.0.0.1", "800" + strconv.Itoa(i))
		worker := proposed.StartWorker(*conn, artificialDelay)
		workerList = append(workerList, *worker)
	}
	workerPool := proposed.CreateWorkerPool(workerList)

	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	gameServer := proposed.StartServer(*conn, *workerPool, artificialDelay)
	game := gameServer.Game

	tickTime := int(tickToTime(30))
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


func TestGame_External_Distributed(t *testing.T) {
	var serverInfo = ServerInfo{
		protocol: "udp",
		address: "192.168.0.20",
		port: "8000",
	}
	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	tickTime := int(tickToTime(60))

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
