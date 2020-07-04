package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Data struct {
	Message string
}

func worker() {
	dst, _ := net.ResolveUDPAddr("udp", ":"+"8002")
	conn, _ := net.ListenUDP("udp", dst)

	time.Sleep(3*time.Second)
	fmt.Printf("before\n")

	recvBuf := make([]byte, 1024)
	n, client, _ := conn.ReadFromUDP(recvBuf[:])
	dec := gob.NewDecoder(bytes.NewReader(recvBuf[:n]))
	request := Data{}
	dec.Decode(&request)
	fmt.Printf("after\n")
	fmt.Println(request.Message)
	fmt.Printf("%+v\n", client)
}

func master(dail *net.UDPConn) {

	request := Data{"hello world"}
	var sendBuf bytes.Buffer
	encoder := gob.NewEncoder(&sendBuf)
	encoder.Encode(request)
	fmt.Printf("before send\n")
	dail.Write(sendBuf.Bytes())
	fmt.Printf("after send\n")
}


func main () {
	workeraddress, _ := net.ResolveUDPAddr("udp", ":8002")
	dail, _ := net.DialUDP("udp", nil, workeraddress)
	go master(dail)
	go worker()
	time.Sleep(5*time.Second)
}