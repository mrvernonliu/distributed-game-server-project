package connection

import (
	"fmt"
	"log"
	"net/rpc"
)

type Connection struct {
	Protocol string
	Address string
	Port string
}

func (connection *Connection) call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP(connection.Protocol, connection.Address+connection.Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func CreateConnection(protocol string, address string, port string) *Connection {
	return &Connection{protocol, address, port }
}