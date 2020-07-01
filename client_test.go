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

func tickToTime(tickRate int) float32 {
	return float32(1000.0 / tickRate)
}



func TestCreatePlayers(t *testing.T) {
	conn := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	gameServer := traditional.StartServer(*conn)
	game := gameServer.Game
	tickTime := int(tickToTime(90))
	time.Sleep(1*time.Second)
	fmt.Println(gameServer)
	for i := 0; i < 100; i++ {
		player := players.CreatePlayer(i)
		go player.JoinGame(conn, tickTime)
		time.Sleep(1*time.Millisecond)
	}

	for !game.IsFinished() {
		time.Sleep(5*time.Second)
	}

}


