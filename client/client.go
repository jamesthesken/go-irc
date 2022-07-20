package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: Implement a Client interface which contains methods read, write, connect, configure, etc.

type CLI struct {
	in io.Reader
}

type Client struct {
	host    string
	channel string
	nick    string
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

	client.nick = nick

	c.Write([]byte(nick))
	c.Write([]byte(user))

	return c, nil
}

// TODO: implement multi-line messages
func (cli *CLI) Write(conn io.Writer) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(cli.in)
	client := &Client{}

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Print(str)
			log.Fatalf("Error reading input: %s", err)
		}

		msg := client.FormatMessage(str)

		// Just makes for easier formatting, as opposed to WriteString()
		fmt.Fprintf(writer, "%s\r\n", msg)
		writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// This may be better implemented as two separate functions
// One to read incoming messages
// Another to format messages
func (cli *CLI) Read(client io.ReadWriter) {
	// Tests pass if you return a string, this works with the server though
	s := bufio.NewScanner(client)
	for s.Scan() {
		line := s.Text()
		fmt.Printf("%s\n", ParseMessage(line))

	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
}

// Possibly create a map that includes characters to remove in the expected msg format

// parseMessage returns formatted incoming messages
func ParseMessage(msg string) Message {
	// kinda wonky, this assumes incoming messages are always in the leading format:
	// ":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python"
	// with that assumption, after the split we only need
	// everything after index 3:
	fullMsg := strings.Split(msg, " ")

	m := Message{}

	// check for ping messages
	if fullMsg[0] == "PING" {
		pong := "PONG" + fullMsg[1]
		return Message{Content: pong}
	} else {

		content := strings.Join(fullMsg[3:][:], " ")

		// nickname of who we received the message from
		msgNick := strings.Split(fullMsg[0], "!")[0]

		timeStamp := time.Now()
		cmd := fullMsg[1]
		if s, err := strconv.Atoi(cmd); err == nil {
			m.Command = s
		}

		m.Channel = fullMsg[2]
		m.Nick = strings.TrimPrefix(msgNick, ":")
		m.Content = strings.TrimPrefix(content, ":")
		m.Time = timeStamp.Format("3:04PM")

		return m
	}

}

// formatMessage returns formatted outgoing messages
func (client *Client) FormatMessage(msg string) string {

	// Check if the message includes a server command
	if strings.HasPrefix(msg, "/") {
		// Maybe a bad assumption that we can properly split the string
		// e.g. - '/nick james' will become [/nick, james]
		cmd := strings.Split(msg, " ")[0]
		// see commands.go for details
		err := executeCommand(client, msg, cmd)
		if err != nil {
			log.Fatal(err)
		}

		return strings.TrimPrefix(msg, "/")
	}

	// otherwise the message is sent to the current channel
	// needs to change to handle private messages of course!
	//  ":james!~james@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #go-nuts :Hello!"
	outMsg := fmt.Sprintf("PRIVMSG %s :%s",
		strings.TrimSuffix(client.channel, "\n"),
		msg)

	return outMsg

}

var wg sync.WaitGroup
