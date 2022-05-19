package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

// Only testing on localhost (for now)
const (
	network = "tcp"
	address = "localhost:3001"
)

func init() {
	server := Server{}

	// Start up the server so it doesn't block the tests
	go func() {
		server.Start(network, address)
	}()

}

// Test: server is up and able to receive client connections
func TestServer(t *testing.T) {
	t.Run("test if the server can accept connections", func(t *testing.T) {
		client, err := net.Dial(network, address)
		assertNoError(t, err)
		defer client.Close()
	})
}

func TestClient(t *testing.T) {
	t.Run("test if the client write to the server", func(t *testing.T) {
		client, err := net.Dial(network, address)
		assertNoError(t, err)
		defer client.Close()

		// Set timeout so we can test cli.Read()
		client.SetDeadline(time.Now().Add(time.Second))

		// Simulate user input
		in := strings.NewReader("Hello, world!\n")
		cli := &CLI{in}

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			cli.Write(client)
		}()

		got := cli.testRead(client)
		want := "Hello, world!"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}

	})
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("Received an error, %s", got)
	}
}

// Test client read by returning a string
func (cli *CLI) testRead(client net.Conn) string {
	s := bufio.NewScanner(client)
	for s.Scan() {
		line := s.Text()
		return line
	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
	return ""
}