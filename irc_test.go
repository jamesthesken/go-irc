package gopherchatv2

import (
	"net"
	"testing"
)

func init() {
	network := "tcp"
	address := "localhost:3000"
	server := Server{}

	// Start up the server so it doesn't block the tests
	go func() {
		server.Start(network, address)
	}()

}

// Test: server is up and able to receive client connections
func TestServer(t *testing.T) {
	t.Run("test if the server can accept connections", func(t *testing.T) {
		network := "tcp"
		address := "localhost:3000"

		client, err := net.Dial(network, address)
		assertNoError(t, err)
		defer client.Close()
	})
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("Received an error")
	}
}
