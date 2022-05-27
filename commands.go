package main

import "strings"

var commands = map[string]func(client *Client, userInput string){
	"/nick": func(client *Client, userInput string) { nick(client, userInput) },
}

func nick(client *Client, userInput string) {
	msg := strings.Split(userInput, " ")
	client.nick = msg[1]
}
