package traditional

import (
	"../../connection"
	"../serverinterfaces"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type PlayerRequest serverinterfaces.PlayerRequest
type ServerResponse serverinterfaces.ServerResponse


type TServer struct {
	Id int

	conn *net.UDPConn
	dst  *net.UDPAddr
}

func (server *TServer) validateAction(request *PlayerRequest) {
	//go fmt.Printf("Server validating request: %d - %d", request.Id, request.Tick)
}

func createResponse(request PlayerRequest) ServerResponse {
	//server.validateAction(request)
	response := ServerResponse{}
	response.Id = request.Id
	response.Tick = request.Tick
	response.Direction = request.Direction
	response.Alive = request.Alive
	response.UniqueIdentifier = request.UniqueIdentifier
	response.X = request.X
	response.Y = request.Y
	return response
}

func (server *TServer) serve() error {
	fmt.Println("Server listening")
	for {
		// Read from UDP
		recvBuf := make([]byte, 1024)
		n, client, _ := server.conn.ReadFromUDP(recvBuf[:])
		dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
		request := PlayerRequest{}
		dec.Decode(&request)

		//Make response
		response := createResponse(request)

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

	go server.serve()
	return &server
}