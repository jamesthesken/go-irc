package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
	listener net.Listener
}

func (s *Server) Start(network, address string) (err error) {
	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}

	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConnection(conn)
	}
}

func handleConnection(conn io.ReadWriter) {
	// Incoming message from client
	for {
		s, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errors.New("oh no!")
		}
		// cmd := strings.Trim(s, "\n")
		fmt.Print(s)
		// // broadcast message to clients
		writer := bufio.NewWriter(conn)

		writer.WriteString(s)

		writer.Flush()
	}

}
