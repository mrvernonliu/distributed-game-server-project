package players

import (
	"./actioninterfaces"
	"../connection"
	"../servers/serverinterfaces"
	"../game/gameinterfaces"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)


type ServerResponse serverinterfaces.ServerResponse

type Player struct {
	// Player variables
	stateMux sync.Mutex
	id int
	uniqueIdentifier int
	x int
	y int
	direction int
	Alive bool

	actionMux sync.Mutex
	actionList [] actioninterfaces.ActionUpdate

	// Server state
	players [] gameinterfaces.InGamePlayer
	phase int

	// Connection variables
	rttLogs     [] int
	tick        int
	tickTime    int
	lostPackets int
	conn	    *net.UDPConn
	dst         *net.UDPAddr
}

func createPlayerRequest(player Player) serverinterfaces.PlayerRequest {
	request := serverinterfaces.PlayerRequest{
		Id:               player.id,
		UniqueIdentifier: player.uniqueIdentifier,
		X:                player.x,
		Y:                player.y,
		Direction:        player.direction,
		Alive:            player.Alive,
		ActionList:       player.actionList,
		Tick:             player.tick,
	}
	return request
}

func (player *Player) callServer() {
	before := time.Now()

	request := createPlayerRequest(*player)
	var sendBuf bytes.Buffer
	encoder := gob.NewEncoder(&sendBuf)
	encoder.Encode(request)

	go player.conn.Write(sendBuf.Bytes())

	recvBuf := make([]byte, 4096)
	n, _, err := player.conn.ReadFromUDP(recvBuf[:])
	if err != nil{
		fmt.Println(err)
	}

	dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
	response := ServerResponse{}
	dec.Decode(&response)
	//go fmt.Printf("Client %d: %+v\n\n\n", player.id, response)
	//go fmt.Printf("Client %d - response: %d %t %+v\n\n\n", player.id, response.Id, response.Alive, response)

	after := time.Now()
	player.stateMux.Lock()
	player.rttLogs = append(player.rttLogs, RTT(before, after))
	//go fmt.Printf("player- response: %d %t\n", response.Id, response.Alive)
	if response.Tick != player.tick || response.Id != player.id {player.lostPackets++}
	if response.Alive == false {
		player.Alive = false
		return
	}
	if response.Tick == player.tick && response.Id == player.id {
		player.players = response.Players
	}

	player.stateMux.Unlock()

	//player.actionMux.Lock()
	player.actionList = []actioninterfaces.ActionUpdate{}
	//player.actionMux.Unlock()
	//fmt.Printf("%d %d\n",response.Tick, player.tick)
}

func (player *Player) setAction() {
	//player.actionMux.Lock()
	//fmt.Println("in setAction mux")
	actionUpdate := GetRandomAction(player.id, player.players)
	if actionUpdate.Action == -1 {return}
	player.actionList = append(player.actionList, actionUpdate)
	//go fmt.Println(player.actionList)
	//player.actionMux.Unlock()
}

func (player *Player) Run() {
	sleepTime := time.Duration(player.tickTime)* time.Millisecond
	for {
		//go fmt.Printf("%d running\n", player.id)
		if !player.Alive {
			//go fmt.Printf("%d - eliminated\n", player.id)
			break
		}
		player.tick++
		if len(player.players) == 100 {
			go player.setAction()
		}
		go player.callServer()
		time.Sleep(sleepTime)
		//if player.tick == 5000 {player.Alive = false}
	}
	time.Sleep(5*time.Second)
	//fmt.Print(len(player.rttLogs))
	//fmt.Printf("Loss: %d\n", player.lostPackets)
}

func (player *Player) JoinGame(connection *connection.Connection, refreshRate int) {
	player.tickTime = refreshRate
	player.dst, _ = net.ResolveUDPAddr("udp", connection.Address+":"+connection.Port)
	player.conn, _ = net.DialUDP("udp", nil, player.dst)

	//go fmt.Printf("%d - Connected to %x with tick time %d\n", player.id, player.dst, player.tickTime)
	player.Run()
}

func CreatePlayer(id int) *Player {
	rand.Seed(time.Now().UTC().UnixNano())

	player := &Player{}
	player.id = id
	player.uniqueIdentifier = rand.Int()
	player.x = rand.Intn(200 - 0)
	player.y = rand.Intn(200 - 0)
	player.actionList = []actioninterfaces.ActionUpdate{}
	player.players = []gameinterfaces.InGamePlayer{}
	player.rttLogs = []int{}
	player.direction = GetRandomDirection()
	player.tick = 0
	player.lostPackets = 0
	player.Alive = true
	player.phase = 0

	return player
}


/*
	Random helper methods
 */

func RTT(before time.Time, after time.Time) int {
	return int((after.UnixNano()-before.UnixNano())/int64(time.Millisecond))
}

func (player *Player) GetNetworkStats() (int, int, int) {
	maxRtt := -1
	for _, rtt := range player.rttLogs {
		if rtt > maxRtt {maxRtt = rtt}
	}
	lossRate := int(float32(player.lostPackets) * 100 / float32(len(player.rttLogs)))
	return player.id, maxRtt, lossRate
}

