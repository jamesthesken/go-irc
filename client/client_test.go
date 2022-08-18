package client

import (
	"bytes"
	"testing"
)

var b bytes.Buffer

func TestClient(t *testing.T) {
	t.Run("test client writing", func(t *testing.T) {
		client := Client{
			Nick:    "james",
			Channel: "#go-nuts",
		}

		client.Write(&b, "Hello, world!", false)

		got := b.String()
		want := "PRIVMSG #go-nuts :Hello, world!\r\n"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}

	})

	t.Run("test command handling", func(t *testing.T) {
		client := Client{
			Conn: &b,
			Nick: "Guest56",
		}

		// cmd := Command{
		// 	Cmd:    "/nick",
		// 	Params: []string{"james"},
		// }

		client.Write(&b, "/nick james", false)
		// err := client.HandleCommand(cmd)
		// assertError(t, err)

	})

	t.Run("test successful join message", func(t *testing.T) {

		client := Client{
			Nick:     "Guest56",
			Channel:  "#python",
			Channels: []string{"#python"},
		}

		// if a client joins a channel successfully, update the client's state
		msg := ParseMessage(":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com JOIN #go-nuts", &client)

		got := client.Channels[1]

		want := "#go-nuts"

		if got != want {
			t.Errorf("got %q with message: %T, wanted %q", got, msg, want)
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

func assertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}
