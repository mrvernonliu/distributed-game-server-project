package proposedWithDistributor

import (
	"../../connection"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Distributor struct {
	Id   int
	WorkerQueue WorkerPool
}

type DistributorRequest struct {
	Request string
	Address string
	Port string
}

type DistributorResponse struct {
	Response bool
	Address string
	Port string
}

func (distributor *Distributor) GetWorker(request *DistributorRequest, response *DistributorResponse) error {
	worker := distributor.WorkerQueue.popIdle()
	if worker.Address == "-1" {
		response.Response = false
	} else {
		response.Response = true
		response.Address = worker.Address
		response.Port = worker.Port
	}
	return nil
}

func (distributor *Distributor) ReturnWorker(request *DistributorRequest, response *DistributorResponse) error {
	worker := WorkerAddress{
		Address: request.Address,
		Port:    request.Port,
	}
	distributor.WorkerQueue.pushIdle(worker)
	return nil
}

func (distributor *Distributor) serve(connection connection.Connection) {
	fmt.Printf("distributor connection: %+v\n", connection)
	rpc.Register(distributor)
	rpc.HandleHTTP()
	l, e := net.Listen(connection.Protocol, connection.Address + ":" + connection.Port)
	if e != nil {
		log.Fatal("listen error:",e)
	}
	go http.Serve(l, nil)
}



func StartDistributor(connection connection.Connection, workerPool WorkerPool) *Distributor {
	rand.Seed(time.Now().UTC().UnixNano())
	distributor := Distributor{}
	distributor.Id = rand.Int()
	distributor.WorkerQueue = workerPool

	distributor.serve(connection)
	return &distributor
}