# Gopher Chat

![build-badge](https://app.travis-ci.com/jamesthesken/go-irc.svg?branch=master)

This is an IRC client written in Go, using [Bubble Tea](https://github.com/charmbracelet/bubbletea) for a built out TUI experience.

The mission is to be the user-friendliest IRC tui on the market.

![screenshot](./assets/ui.png)

### Current Features
* The usual IRC Commands (join, quit, nick, privmsg, etc.)
* Send/Receive messages

### Roadmap
In order of priority:
* Toggle view between active channels
* User configuration page
* Support auto-completion for commands

### How to Run:

This software is tested on `go1.18.1`

First, edit the config file located in the `config` directory. Replace the values for nick, host, and channel with your own. Upon startup you join the `default_channel`, currently set to #test.

Then:

```
git clone https://github.com/jamesthesken/go-irc.git
cd go-irc
go build
./gopherchatv2
```

### There are a lot of bugs!
This code needs work, so expect some instability.

