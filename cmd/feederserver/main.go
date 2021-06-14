package main

import "feeder-server/server"

const (
	connHost = "localhost"
	connPort = "4000"
)

func main() {
	server.StartServer(connHost+":"+connPort)
}





