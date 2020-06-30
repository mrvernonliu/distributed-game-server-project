package traditional

import (
	"../serverrpc"
	"../../connection"
	//"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"net/http"
	"time"
)

type PlayerRequest serverrpc.PlayerRequest
type ServerResponse serverrpc.ServerResponse

type TServer struct {
	Id int
}

func (server *TServer) validateAction(request *PlayerRequest) {
	//go fmt.Printf("Server validating request: %d - %d", request.Id, request.Tick)
}

func (server *TServer) UpdatePlayerState(request *PlayerRequest, response *ServerResponse) error {
	server.validateAction(request);
	response.Id = request.Id
	response.Tick = request.Tick
	response.Direction = request.Direction
	response.Alive = request.Alive
	response.UniqueIdentifier = request.UniqueIdentifier
	response.X = request.X
	response.Y = request.Y
	return nil
}

func (server *TServer) serve(connection connection.Connection) {
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen(connection.Protocol, connection.Address + ":" + connection.Port)
	if e != nil {
		log.Fatal("listen error:",e)
	}
	go http.Serve(l, nil)
}


func StartServer(connection connection.Connection) *TServer {
	rand.Seed(time.Now().UTC().UnixNano())
	server := TServer{}
	server.Id = rand.Int()
	server.serve(connection)
	return &server
}