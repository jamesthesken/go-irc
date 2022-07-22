package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Host     string
	Nick     string
	Channel  string
	Channels []string
}

// Server operations
func (client *Client) Connect(server string) (net.Conn, error) {

	// SSL
	conf := &tls.Config{}

	c, err := tls.Dial("tcp", server, conf)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, err
	}

	// RFC 1459 - 4.1.2/3 - NICK and USER messages
	nick := "NICK samej1293871\n"
	user := "USER samej1293871 * * :samej1293871\n"

	client.Nick = nick

	c.Write([]byte(nick))
	c.Write([]byte(user))

	return c, nil
}

// parseMessage returns formatted incoming messages
func ParseMessage(msg string) Message {
	// incoming messages are in the leading format:
	// ":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python"
	fullMsg := strings.Split(msg, " ")

	m := Message{}

	// check for ping messages
	if fullMsg[0] == "PING" {
		pong := "PONG " + fullMsg[1]
		return Message{Content: pong, Ping: true}
	} else {

		content := strings.Join(fullMsg[3:][:], " ")

		// nickname of who we received the message from
		msgNick := strings.Split(fullMsg[0], "!")[0]

		timeStamp := time.Now()
		cmd := fullMsg[1]
		if s, err := strconv.Atoi(cmd); err == nil {
			m.NumReply = s
		} else {
			m.Notification = cmd
		}

		m.Channel = fullMsg[2]
		m.Nick = strings.TrimPrefix(msgNick, ":")
		m.Content = strings.TrimPrefix(content, ":")
		m.Time = timeStamp.Format("3:04PM")

		return m
	}

}

// FormatMessage() returns formatted outgoing messages.
func (client *Client) FormatMessage(msg string) string {
	// Check if the message includes a server command
	if strings.HasPrefix(msg, "/") {
		contents := strings.Split(msg, " ")
		if len(contents) >= 1 {
			switch contents[0] {
			case "/nick":
				client.Nick = contents[1]

			case "/join":
				client.Channel = strings.Split(msg, " ")[1]
			}
		}
		return strings.TrimPrefix(msg, "/")
	}

	// otherwise the message is sent to the current channel
	// needs to change to handle private messages of course!
	//  ":james!~james@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #go-nuts :Hello!"
	outMsg := fmt.Sprintf("PRIVMSG %s :%s",
		strings.TrimSuffix(client.Channel, "\n"),
		msg)

	return outMsg

}
