package client

/*
	messages.go implements message types that are sent to the TUI
	The TUI utilizes a switch statement
	based on which messages are received.
*/

// idea - move ErrReply struct to Message

type Message struct {
	Content      string
	Nick         string
	Channel      string
	Time         string
	Notification string
	NumReply     int
	Ping         bool
}

type Command struct {
	Cmd    string
	Params []string
}

var cmd = map[string]string{
	"/nick": "set nickname",
	"/join": "join channel",
	"/part": "leave a channel",
}
