package project

import (
	"./connection"
	"./players"
	"fmt"
	"testing"
)

type ServerInfo struct {
	protocol string
	address string
	port string
}
var serverInfo = ServerInfo{
	protocol: "",
	address: "",
	port: "",
}

func TestCreatePlayers(t *testing.T) {
	gameServer := connection.CreateConnection(serverInfo.protocol, serverInfo.address, serverInfo.port)
	for i := 0; i < 100; i++ {
		player := players.CreatePlayer(i)
		go player.JoinGame(gameServer)
		fmt.Printf("%+v\n", *player)
	}
}


