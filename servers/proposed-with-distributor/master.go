package proposedWithDistributor

import (
	"../../connection"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"time"
	"../serverinterfaces"
)

type Master struct {
	Id   int
	Game *DistributedGame

	conn *net.UDPConn
	dst  *net.UDPAddr
}

func (server *Master) serve() error {
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
		response := server.Game.sendValidationToWorker(request)

		//go fmt.Printf("server- response: %+v\n", response)
		//go fmt.Printf("server- response: %d %t\n", response.Id, response.Alive)
		//go fmt.Printf("%+v\n", client)
		//send response
		var sendBuf bytes.Buffer
		encoder := gob.NewEncoder(&sendBuf)
		encoder.Encode(response)
		server.conn.WriteToUDP(sendBuf.Bytes(), client)

	}
	return nil
}


func StartServer(connection connection.Connection, artificalDelay int, distributor connection.Connection) *Master {
	rand.Seed(time.Now().UTC().UnixNano())
	server := Master{}
	server.Id = rand.Int()
	server.dst, _ = net.ResolveUDPAddr("udp", ":"+connection.Port)
	server.conn, _ = net.ListenUDP("udp", server.dst)
	server.Game = CreateGame(artificalDelay, distributor)

	go server.serve()
	return &server
}
