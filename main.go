package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter    int
	listenAddr = "localhost:8080"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err)
		}

		backend := chooseBackend()
		fmt.Printf("counter=%d backend=%s\n", counter, backend)
		go func() {
			err := proxy(backend, conn)
			if err != nil {
				log.Printf("WARNING: proxying failed: %v", err)
			}
		}()
	}
}

func proxy(backend string, c net.Conn) error {
	//Dial connects to the address on the named network.
	bc, err := net.Dial("tcp", backend)
	if err != nil {
		return fmt.Errorf("failed to connect to backend %s: %v", backend, err)
	}

	//for those comes from connection. connection -> backend
	go io.Copy(bc, c)
	//for those comes from backend. backend -> connection
	go io.Copy(c, bc)

	return nil
}

func chooseBackend() string {
	s := server[counter%len(server)]
	counter++
	return s
}
