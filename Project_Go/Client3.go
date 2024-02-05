package main

import (
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("Départ : {1,1}, Arrivée : {5,1}"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	received := ""
	var builder strings.Builder

	for {
		n, err := conn.Read(buf)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		received += string(buf[:n])

		builder.Write(buf[:n])
		if strings.Contains(builder.String(), "\n") {
			println("Received message:", received)
			conn.Close()
			break
		}
	}

}
