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
	FullName string
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
	// client.Channels = make([]string)

	c.Write([]byte(nick))
	c.Write([]byte(user))

	return c, nil
}

// ParseMessage returns formatted incoming messages
func ParseMessage(msg string, client *Client) Message {
	// incoming messages are in the leading format:
	// ":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python"
	fullMsg := strings.Split(msg, " ")

	m := Message{}

	// check for ping messages
	if fullMsg[0] == "PING" {
		pong := "PONG " + fullMsg[1]
		return Message{Content: pong, Ping: true}
	}

	content := strings.Join(fullMsg[3:][:], " ")

	// nickname of who we received the message from
	msgNick := strings.Split(fullMsg[0], "!")[0]

	timeStamp := time.Now()
	cmd := fullMsg[1]
	if s, err := strconv.Atoi(cmd); err == nil {
		m.NumReply = s
	} else {
		m.Notification = cmd
		m.NumReply = 0
	}

	m.Channel = fullMsg[2]
	m.Nick = strings.TrimPrefix(msgNick, ":")
	m.Content = strings.TrimPrefix(content, ":")
	m.Time = timeStamp.Format("3:04PM")

	if m.NumReply == RplEndOfNames {
		client.Channels = append(client.Channels, fullMsg[3])
	}

	return m
}

// UpdateClient
func (client *Client) UpdateClient(msg Message) {
	if msg.NumReply < 401 {
		client.Nick = msg.Nick
		client.Channel = msg.Channel
		client.Channels = append(client.Channels, msg.Channel)
	}

}

// FormatMessage returns formatted outgoing messages.
func (client *Client) FormatMessage(msg string) string {
	// Check if the message includes a server command
	if strings.HasPrefix(msg, "/") {
		contents := strings.Split(msg, " ")
		if len(contents) >= 1 {
			switch contents[0] {
			case "/part":
				client.Channels = removeChannel(client.Channels, contents[1])
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

func removeChannel(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
