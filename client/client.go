package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Conn     io.ReadWriter
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

	client.Conn = c

	// RFC 1459 - 4.1.2/3 - NICK and USER messages
	nick := "NICK samej1293871\n"
	user := "USER samej1293871 * * :samej1293871\n"

	client.Nick = nick
	// client.Channels = make([]string)

	c.Write([]byte(nick))
	c.Write([]byte(user))

	return c, nil
}

func (client *Client) Write(conn io.Writer, msg string, ping bool) {
	writer := bufio.NewWriter(conn)

	if !ping {
		// formats the message into one acceptable by IRC
		msg = client.FormatMessage(msg)
	}

	// Just makes for easier formatting, as opposed to WriteString()
	fmt.Fprintf(writer, "%s\r\n", msg)
	writer.Flush()
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

	// TODO: handle error messages from the server

	content := strings.Join(fullMsg[3:][:], " ")

	// nickname of who we received the message from
	msgNick := strings.Split(fullMsg[0], "!")[0]

	timeStamp := time.Now()
	cmd := fullMsg[1]
	// if there's a numeric reply from the server, save it to NumReply
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

	if m.Nick == client.Nick {
		client.Channel = m.Channel
		client.Channels = append(client.Channels, client.Channel)
	}

	return m
}

// FormatMessage returns formatted outgoing messages.
func (client *Client) FormatMessage(msg string) string {
	// Check if the message includes a server command

	// TODO: Handle commands and check if the server responds with an error message, before changing state

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

// setNick modifies the clients
// client sends NICK request
func (client *Client) setNick(cmd Command) {

	client.Nick = cmd.Params[0]
}

func removeChannel(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
