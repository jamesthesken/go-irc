package gopherchatv2

import (
	"bufio"
	"errors"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
}

func (s *Server) Start(network, address string) (err error) {
	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}
	defer s.listener.Close()

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

func handleConnection(c net.Conn) {
	// Incoming message from client
	s, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		errors.New("oh no!")
	}
	cmd := strings.Trim(s, "\n")

	// echo message back
	c.Write([]byte(cmd))
	c.Close()
}
