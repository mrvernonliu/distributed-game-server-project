package project

import (
	"./connection"
	"./players"
	"./servers/proposed"
	"./servers/traditional"
	"./servers/proposed-with-distributor"
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
		//fmt.Printf("here %+v\n", player)
		id, maxRtt, lossRate := player.GetNetworkStats()
		resultEntry := convertRttToString(player.RttLogs) + "\n"
		fileOutput += resultEntry
		playerStats = append(playerStats, playerStatistic{
			playerId: id,
			maxRtt: maxRtt,
			lossRate: lossRate,
		})
	}
	//fmt.Printf(fileOutput)

	err := ioutil.WriteFile("results.csv", []byte(fileOutput), 0644)
	if err != nil {
		fmt.Printf("error %+v\n", err)
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
	tickTime := int(tickToTime(20))
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
		address: "192.168.0.20",
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

	tickTime := int(tickToTime(10))
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
	for _, worker := range workerList {
		worker.Kill()
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
	tickTime := int(tickToTime(20))

	var playerList []*players.Player

	for i := 0; i < 100; i++ {
		player := players.CreatePlayer(i)
		playerList = append(playerList, player)
		go player.JoinGame(conn, tickTime)
		time.Sleep(1*time.Millisecond)
	}

	time.Sleep(15*time.Second)
	fmt.Printf("displaying stats")
	displayPlayerStatistics(playerList)
}

func TestDistributor(t *testing.T) {
	/*
		Grabs more workers than necessary, then returns all the workers,
		then grabs more than necessary.
		If it hangs then it fails, should be able to continue working if worker cannot be grabbed
	 */
	artificialDelay := 1
	var workerList []proposedWithDistributor.Worker
	for i := 5; i < 10; i++ {
		conn := connection.CreateConnection("udp", "127.0.0.1", "800" + strconv.Itoa(i))
		worker := proposedWithDistributor.StartWorker(*conn, artificialDelay)
		workerList = append(workerList, *worker)
	}
	workerPool := proposedWithDistributor.CreateWorkerPool(workerList)

	conn := connection.CreateConnection("tcp", "127.0.0.1", "8080")
	proposedWithDistributor.StartDistributor(*conn, *workerPool)

	workerChannel := make(chan proposedWithDistributor.WorkerAddress, 6)
	for i := 0; i < 10; i++ {
		req := proposedWithDistributor.DistributorRequest{
			Request: "get",
		}
		var res proposedWithDistributor.DistributorResponse
		conn.Call("Distributor.GetWorker", &req, &res)
		if res.Response == true {
			fmt.Printf("got worker: %+v\n", res)
			workerChannel <- proposedWithDistributor.WorkerAddress{
				Address: res.Address,
				Port:    res.Port,
			}
		}
	}
	for i := 0; i < 5; i++ {
		worker := <- workerChannel
		req := proposedWithDistributor.DistributorRequest{
			Request: "put",
			Address: worker.Address,
			Port: worker.Port,
		}
		fmt.Printf("returning: %+v\n", req)
		var res proposedWithDistributor.DistributorResponse
		conn.Call("Distributor.ReturnWorker", &req, &res)
	}

	for i := 0; i < 10; i++ {
		req := proposedWithDistributor.DistributorRequest{
			Request: "get",
		}
		var res proposedWithDistributor.DistributorResponse
		conn.Call("Distributor.GetWorker", &req, &res)
		if res.Response == true {
			fmt.Printf("got worker: %+v\n", res)
			workerChannel <- proposedWithDistributor.WorkerAddress{
				Address: res.Address,
				Port:    res.Port,
			}
		}
	}
}

func TestGame_Internal_Distributed_with_Distributor (t *testing.T) {
	// Create Distributor
	artificialDelay := 1
	var workerList []proposedWithDistributor.Worker
	for i := 5; i < 10; i++ {
		conn := connection.CreateConnection("udp", "127.0.0.1", "800" + strconv.Itoa(i))
		worker := proposedWithDistributor.StartWorker(*conn, artificialDelay)
		workerList = append(workerList, *worker)
	}
	workerPool := proposedWithDistributor.CreateWorkerPool(workerList)

	distributorConn := connection.CreateConnection("tcp", "127.0.0.1", "8080")
	proposedWithDistributor.StartDistributor(*distributorConn, *workerPool)

	// Create Game
	var serverInfo = ServerInfo{
		protocol: "udp",
		address: "127.0.0.1",
		port: "8000",
	}
	gameServerConn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	gameServer := proposedWithDistributor.StartServer(*gameServerConn, artificialDelay, *distributorConn)
	game := gameServer.Game

	tickTime := int(tickToTime(10))
	var playerList []*players.Player

	for i := 0; i < 20; i++ {
		player := players.CreatePlayer(i)
		playerList = append(playerList, player)
		go player.JoinGame(gameServerConn, tickTime)
		time.Sleep(1*time.Millisecond)
	}

	for !game.IsFinished() {
		time.Sleep(5*time.Second)
	}
}
