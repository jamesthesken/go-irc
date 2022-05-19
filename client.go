package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type CLI struct {
	in io.Reader
}

func (cli *CLI) Write(client net.Conn) {
	writer := bufio.NewWriter(client)
	reader := bufio.NewReader(cli.in)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %s", err)
		}

		writer.WriteString(str)
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

var wg sync.WaitGroup

func (cli *CLI) Read(client net.Conn) {
	// Tests pass if you return a string, this works with the server though
	s := bufio.NewScanner(client)
	for s.Scan() {
		line := s.Text()
		log.Printf("Server: %s", line)
	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
}

// func main() {
// 	wg.Add(1)

// 	client, err := net.Dial("tcp", "127.0.0.1:3001")
// 	if err != nil {
// 		log.Fatalf("Error: %s", err)
// 	}
// 	cli := CLI{os.Stdin}

// 	go cli.Read(client)
// 	go cli.Write(client)
// 	wg.Wait()

// }
