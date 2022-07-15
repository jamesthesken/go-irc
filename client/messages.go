package client

/*
	messages.go implements message types that are sent to the TUI
	The TUI utilizes a switch statement
	based on which messages are received.
*/

type Message struct {
	msg     string
	nick    string
	channel string
	time    string
}
