package client

import (
	"log"
	"strings"
)

// TODO - implement a command interface

var commands = map[string]func(client *Client, userInput string){
	"/nick": func(client *Client, userInput string) { nick(client, userInput) },
	"/join": func(client *Client, userInput string) { join(client, userInput) },
}

func nick(client *Client, userInput string) {
	msg := strings.Split(userInput, " ")
	client.nick = msg[1]
}

// As the user joins more channels, client.channel should be
// an arry of current channels
func join(client *Client, userInput string) {
	msg := strings.Split(userInput, " ")
	client.channel = msg[1]
}

func part(client *Client, userInput string) {
	client.channel = ""
}

// executeCommand takes a cmd and messaage and executes a command
func executeCommand(client *Client, msg, cmd string) (err error) {
	// If the string value of the command exists in the commands map,
	// execute the state-changing function
	if val, exists := commands[cmd]; exists {
		val(client, msg)
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
