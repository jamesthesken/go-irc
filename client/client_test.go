package client

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Run("test ingress message formatting", func(t *testing.T) {
		got := ParseMessage(":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python :Please use a test channel https://foo.bar.com")

		now := time.Now()

		want := Message{
			Content: "Please use a test channel https://foo.bar.com",
			Nick:    "Guest56",
			Channel: "#python",
			Time:    now.Format("3:04PM"),
		}

		if want != got {
			t.Errorf("got %s, wanted %s", got.Time, want.Time)
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
