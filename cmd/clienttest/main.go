package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// Application constants, defining host, port, and protocol.
const (
	connHost = "localhost"
	connPort = "4000"
	connType = "tcp"
)

func main() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("PRESS INTRO TO TERMINATE THE SERVER.")
		reader.ReadString('\n')
		fmt.Println("TERMINATING THE SERVER")
		StartClientConnection("terminate\n")
	}()

	wg := &sync.WaitGroup{}
	productSkuList := CreateProductSkuList()
	for _, productSku := range productSkuList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			StartClientConnection(productSku)
		}()
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
}

func StartClientConnection(message string) {
	//Start the client and connect to the server.
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		log.Println("Error connecting to server: ", err.Error())
		os.Exit(1)
	}
	log.Println("Connected to", connType, "server", connHost+":"+connPort)
	defer conn.Close()
	//Send message to the server
	conn.Write([]byte(message + "\n"))
	log.Println("Sent to server: ", message)
}

func CreateProductSkuList() []string {
	duplicated := 50
	invalid := 20
	var productSkuList []string
	for i := 1000; i < 10000; i++ {
		productSkuList = append(productSkuList, "ABCD-"+strconv.Itoa(i))
		if duplicated > 0 {
			productSkuList = append(productSkuList, "ABCD-"+strconv.Itoa(i))
			duplicated--
		}
		if invalid > 0 {
			productSkuList = append(productSkuList, "ZZ"+strconv.Itoa(i))
			invalid--
		}
	}
	return productSkuList
}
