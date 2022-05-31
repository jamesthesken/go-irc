package main

import (
	"fmt"
	"net"
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
	t.Run("test ingress message formatting", func(t *testing.T) {
		got := parseMessage(":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python :Please use a test channel https://foo.bar.com")

		now := time.Now()

		want := fmt.Sprintf("%s %s", now.Format("3:04PM"), "< Guest56 > Please use a test channel https://foo.bar.com")

		if want != got {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})

	t.Run("test command execution", func(t *testing.T) {

		client := Client{}

		client.formatMessage("/nick james")

		got := client.nick
		want := "james"

		if want != got {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})

	t.Run("test outgoing message formatting", func(t *testing.T) {

		client := Client{
			nick:    "james",
			channel: "#go-nuts",
		}

		got := client.formatMessage("Hello!")

		want := "PRIVMSG #go-nuts :Hello!"

		if want != got {
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
