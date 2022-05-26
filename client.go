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
	client, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, err
	}

	// RFC 1459 - 4.1.2/3 - NICK and USER messages
	nick := "NICK samej1293871\n"
	user := "USER samej1293871 * * :samej1293871\n"
	client.Write([]byte(nick))

	client.Write([]byte(user))

	return client, nil
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

		fmt.Printf("%s\n", formatMessage(line))

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

func formatMessage(msg string) string {

	timeStamp := time.Now()
	contents := strings.Split(msg, ":")

	formatted := fmt.Sprintf("%s < %s > %s",
		timeStamp.Format("3:04PM"),
		strings.Split(contents[1], "!")[0],
		contents[len(contents)-1:])

	return formatted

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
