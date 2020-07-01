package traditional

import (
	"../../connection"
	"../serverinterfaces"
	"../../game"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type TServer struct {
	Id   int
	Game *game.Game

	conn *net.UDPConn
	dst  *net.UDPAddr
}

func (server *TServer) validateAction(request *serverinterfaces.PlayerRequest) {
	//go fmt.Printf("Server validating request: %d - %d", request.Id, request.Tick)
}


func (server *TServer) serve() error {
	fmt.Println("Server listening")
	for {
		// Read from UDP
		recvBuf := make([]byte, 1024)
		n, client, _ := server.conn.ReadFromUDP(recvBuf[:])
		dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
		request := serverinterfaces.PlayerRequest{}
		dec.Decode(&request)
		//go fmt.Printf("Server - request: %+v %+v\n", request, request.ActionList)

		//Make response
		response := server.Game.UpdateGameState(request)

		//go fmt.Printf("server- response: %+v\n", response)
		//go fmt.Printf("server- response: %d %t\n", response.Id, response.Alive)

		//send response
		var sendBuf bytes.Buffer
		encoder := gob.NewEncoder(&sendBuf)
		encoder.Encode(response)
		server.conn.WriteToUDP(sendBuf.Bytes(), client)

	}
	return nil
}



func StartServer(connection connection.Connection) *TServer {
	rand.Seed(time.Now().UTC().UnixNano())
	server := TServer{}
	server.Id = rand.Int()
	server.dst, _ = net.ResolveUDPAddr("udp", connection.Address+":"+connection.Port)
	server.conn, _ = net.ListenUDP("udp", server.dst)
	server.Game = game.CreateGame()

	go server.serve()
	return &server
}