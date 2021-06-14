package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sync"
	"time"
)

const (
	connType             = "tcp"
	maxConcurrentClients = 5
	executionTime        = 60
	fileName             = "product-sku-list"
)

type Server struct {
	listener             net.Listener
	listenerClosed       chan interface{}
	maxConcurrentClients chan struct{}
	timer                *time.Timer
	wg                   *sync.WaitGroup
	m                    *sync.RWMutex
	file                 *os.File
	report               *Report
}

type Report struct {
	unique     map[string]struct{}
	duplicated []string
	invalid    []string
}

func StartServer(address string) {
	s := &Server{
		listenerClosed:       make(chan interface{}),
		maxConcurrentClients: make(chan struct{}, maxConcurrentClients),
		timer:                time.NewTimer(executionTime * time.Second),
		wg:                   &sync.WaitGroup{},
		m:                    &sync.RWMutex{},
		report: &Report{
			unique: make(map[string]struct{}),
		},
	}
	//Start listening for incoming connections
	l, err := net.Listen(connType, address)
	if err != nil {
		log.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	s.listener = l
	log.Println("Starting " + connType + " server on " + s.listener.Addr().String())

	f, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file: ", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	s.file = f
	log.Println("File created: " + f.Name())

	//Wait for timer to finish before closing Server connection
	go func() {
		<-s.timer.C
		fmt.Println("Timeout")
		s.closeServer()
	}()

	s.serve()
	//Wait for all open connections to finish
	s.wg.Wait()
	s.file.Close()
	s.printReport()
}

func (s *Server) serve() {
	defer close(s.maxConcurrentClients)
	for {
		//If maxConcurrentClients is 5, wait until one is released from channel
		s.maxConcurrentClients <- struct{}{}
		//Wait until accepting a new connection
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.listenerClosed:
				fmt.Println("Server has been closed")
				return
			default:
				log.Println("Error accepting connection: ", err)
				<-s.maxConcurrentClients
				continue
			}
		} else {
			log.Println("Accepted client connection: " + conn.RemoteAddr().String())
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				defer func() { <-s.maxConcurrentClients }()

				message, err := s.handleConnection(conn)
				if err != nil {
					return
				}
				if message != "" {
					err = s.handleMessage(message)
					if err != nil {
						return
					}
				}
			}()
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) (message string, err error) {
	defer conn.Close()
	var buffer []byte
	for {
		select {
		case <-s.listenerClosed:
			//If Server listener has been closed, exit connection loop
			return
		default:
			// Wait 2 seconds until receiving a message from client
			conn.SetDeadline(time.Now().Add(2 * time.Second))
			buffer, err = bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					//If haven't received a message in 2 seconds, try again
					continue
				} else {
					log.Println("Connection reading error: ", err)
					return
				}
			}
			message = string(buffer)
			return
		}
	}

}

func (s *Server) handleMessage(message string) (err error) {
	if message == "terminate\n" {
		log.Println("Client terminate")
		s.closeServer()
		return
	}
	// Check Product SKU format
	matched, _ := regexp.MatchString(`[a-zA-Z]{4}[-][0-9]{4}[\n]`, message)
	if !matched {
		s.m.Lock()
		s.report.invalid = append(s.report.invalid, message)
		log.Println("Invalid Product Sku: ", message)
		s.m.Unlock()
		return
	}
	s.m.RLock()
	_, exists := s.report.unique[message]
	s.m.RUnlock()

	if exists {
		s.m.Lock()
		s.report.duplicated = append(s.report.duplicated, message)
		s.m.Unlock()
		log.Println("Duplicated Product Sku: ", message)
	} else {
		s.m.Lock()
		s.report.unique[message] = struct{}{}
		_, err = fmt.Fprint(s.file, message)
		s.m.Unlock()
		if err != nil {
			log.Println("Error writing to file: ", err)
			return
		}
		log.Println("Unique Product Sku: ", message)
	}
	return
}

func (s *Server) closeServer() {
	fmt.Println("Closing Server")
	close(s.listenerClosed)
	s.listener.Close()
	s.timer.Stop()
}

func (s *Server) printReport() {
	fmt.Printf("Received %d unique product skus, %d duplicates, %d discard values.\n",
		len(s.report.unique), len(s.report.duplicated), len(s.report.invalid))
}
