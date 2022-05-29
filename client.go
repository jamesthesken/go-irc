package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

var commands = map[string]func(client *Client, userInput string){
	"/nick": func(client *Client, userInput string) { nick(client, userInput) },
}

func nick(client *Client, userInput string) {
	msg := strings.Split(userInput, " ")
	client.nick = msg[1]
}

// TODO: implement multi-line messages
func (cli *CLI) Write(client io.Writer) {
	writer := bufio.NewWriter(client)
	reader := bufio.NewReader(cli.in)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Print(str)
			log.Fatalf("Error reading input: %s", err)
		}

		// Just makes for easier formatting, as opposed to WriteString()
		fmt.Fprintf(writer, "%s\r\n", str)
		writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// Server operations
func Connect(server string) (net.Conn, error) {
	c, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, err
	}

	// RFC 1459 - 4.1.2/3 - NICK and USER messages
	nick := "NICK samej1293871\n"
	user := "USER samej1293871 * * :samej1293871\n"

	c.Write([]byte(nick))

	c.Write([]byte(user))

	return c, nil
}

var wg sync.WaitGroup

// This may be better implemented as two separate functions
// One to read incoming messages
// Another to format messages
func (cli *CLI) Read(client io.ReadWriter) {
	// Tests pass if you return a string, this works with the server though
	s := bufio.NewScanner(client)
	for s.Scan() {
		line := s.Text()

		fmt.Printf("%s\n", parseMessage(line))

		if strings.Contains(line, "PING") {
			msg := strings.TrimPrefix(line, "PING")

			pong := "PONG" + msg
			client.Write([]byte(pong))
			log.Printf("Client sent: %s", pong)
		}
	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
}

// Possibly create a map that includes characters to remove in the expected msg format

// parseMessage returns formatted incoming messages
func parseMessage(msg string) string {
	timeStamp := time.Now()
	// kinda wonky, this assumes incoming messages are always in the leading format:
	// ":Guest56!~Guest56@cpe-10.21.123.13.1.foo.bar.com PRIVMSG #python"
	// with that assumption, after the split we only need
	// everything after index 3:
	fullMsg := strings.Split(msg, " ")

	content := strings.Join(fullMsg[3:][:], " ")

	// nickname of who we received the message from
	msgNick := strings.Split(fullMsg[0], "!")[0]

	formatted := fmt.Sprintf("%s < %s > %s",
		timeStamp.Format("3:04PM"),
		strings.TrimPrefix(msgNick, ":"),
		strings.TrimPrefix(content, ":"))

	return formatted
}

// formatMessage returns formatted outgoing messages
func (client *Client) formatMessage(msg string) string {

	// Check if the message includes a server command
	if strings.HasPrefix(msg, "/") {
		// Maybe a bad assumption that we can properly split the string
		// e.g. - '/nick james' will become [/nick, james]
		cmd := strings.Split(msg, " ")[0]

		// If the string value of the command exists in the commands map,
		// execute the state-changing function -> see commands.go for details
		if val, exists := commands[cmd]; exists {
			val(client, msg)
		} else {
			return cmd
		}
	}

	return "formatted"

}

func main() {
	wg.Add(1)

	cli := CLI{os.Stdin}

	client, err := Connect("irc.libera.chat:6667")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	go cli.Read(client)
	go cli.Write(client)
	wg.Wait()

}
