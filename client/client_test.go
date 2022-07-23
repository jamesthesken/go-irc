package client

import (
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("test successful join message", func(t *testing.T) {

		client := Client{
			Nick:    "Guest56",
			Channel: "#python",
		}

		client.Channels = []string{}

		// if a client joins a channel successfully, update the client's state
		msg := ParseMessage(":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com JOIN #go-nuts", &client)

		got := client.Channels[0]

		want := "#go-nuts"

		if got != want {
			t.Errorf("got %q with message: %T, wanted %q", got, msg, want)
		}
	})

	t.Run("test command execution", func(t *testing.T) {

		client := Client{}

		client.FormatMessage("/nick james")

		got := client.Nick
		want := "james"

		if want != got {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})

	t.Run("test outgoing message formatting", func(t *testing.T) {

		client := Client{
			Nick:    "james",
			Channel: "#go-nuts",
		}

		got := client.FormatMessage("Hello!")

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
