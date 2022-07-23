package tui

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"bufio"
	"fmt"
	"gopherchatv2/client"
	"gopherchatv2/tui/constants"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var wg sync.WaitGroup

func StartTea() {

	// start logging
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	m := initialModel()

	wg.Add(1)
	// Connect to IRC
	client := &client.Client{}
	conn, err := client.Connect("irc.libera.chat:6697")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	m.client = client
	m.conn = conn

	p := *tea.NewProgram(m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	go m.Read(conn, &p)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error
type item string
type itemDelegate struct{}
type mode int

func (i item) FilterValue() string { return "" }

const (
	nav mode = iota
	msgMode
)

type Model struct {
	mode        mode
	viewport    viewport.Model
	textarea    textarea.Model
	list        list.Model
	senderStyle lipgloss.Style
	notifStyle  lipgloss.Style
	conn        io.Writer
	client      *client.Client
	messages    []string
	serverMsg   string
	err         error
}

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

func initialModel() *Model {

	// text area
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(2)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// viewport
	vp := viewport.New(5, 2)
	vp.SetContent(`Welcome to IRC!`)
	// TODO - apply a new list of keybindings ...
	vp.KeyMap.PageDown.SetEnabled(false)
	vp.KeyMap.PageUp.SetEnabled(false)
	vp.KeyMap.HalfPageDown.SetEnabled(false)
	vp.KeyMap.HalfPageUp.SetEnabled(false)
	vp.KeyMap.Up.SetEnabled(false)
	vp.KeyMap.Down.SetEnabled(false)

	// channel list
	items := []list.Item{item("Test")}
	const defaultWidth = 20

	list := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	list.SetFilteringEnabled(false)
	list.DisableQuitKeybindings()
	list.Title = "Channels"

	return &Model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		notifStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")),
		list:        list,
		err:         nil,
	}
}

// Init() is the first function called by BubbleTea.
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		liCmd tea.Cmd
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// Implement different tea messages sent by the clients.
	// i.e., Message interface
	case client.Message:
		// if this message is a ping message, don't show it
		if !msg.Ping {

			switch msg.Notification {
			case "JOIN":
				notif := fmt.Sprintf("%s has joined %s", msg.Nick, msg.Channel)
				m.messages = append(m.messages, m.notifStyle.Render(msg.Time)+" "+m.notifStyle.Render(""+notif))
				m.setContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			case "QUIT":
				notif := fmt.Sprintf("%s has left -> %s", msg.Nick, msg.Content)
				m.messages = append(m.messages, m.notifStyle.Render(msg.Time)+" "+m.notifStyle.Render(""+notif))
				m.setContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			case "NICK":
				notif := fmt.Sprintf("%s changed their nick to %s", msg.Nick, msg.Channel)
				m.messages = append(m.messages, m.notifStyle.Render(msg.Time)+" "+m.notifStyle.Render(""+notif))
				m.setContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			case "PRIVMSG":
				m.messages = append(m.messages, m.senderStyle.Render(msg.Time)+" "+m.senderStyle.Render(msg.Nick+": ")+" "+msg.Content)
				m.setContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			case "MODE":
			default:
				m.messages = append(m.messages, m.senderStyle.Render(msg.Time)+" "+m.senderStyle.Render(msg.Nick+" ")+" "+msg.Content)
				m.setContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			}
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - msg.Width/4
		m.viewport.Height = msg.Height - msg.Height/4
		m.setContent(strings.Join(m.messages, "\n"))
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Tab):
			m.toggleBox()
		case key.Matches(msg, constants.Keymap.Enter):
			if m.textarea.Focused() {
				timeStamp := time.Now()
				m.messages = append(m.messages, m.senderStyle.Render(timeStamp.Format("3:04PM"+" < You > "))+m.textarea.Value())
				m.Write(m.textarea.Value(), false)
				m.textarea.Reset()
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	channels := channelsToItems(m.client.Channels)
	m.list.SetItems(channels)

	m.list, liCmd = m.list.Update(msg)
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd, liCmd)
}

func (m Model) View() string {
	m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("212")).Width(m.viewport.Width)

	// channel pane
	left := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("212")).Height(m.viewport.Height).Width(m.viewport.Width / 7).Padding(1).Render(m.list.View())
	right := m.viewport.View()
	bottomRight := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(1).Width(m.viewport.Width).BorderForeground(lipgloss.Color("212")).Padding(1).Render(m.textarea.View())

	// chat window and input
	rightPane := lipgloss.JoinVertical(lipgloss.Center, right, bottomRight)

	formatted := lipgloss.JoinHorizontal(lipgloss.Left, left, rightPane)

	return constants.DocStyle.Render(formatted)
}

// Read() receives messages from the IRC server and outputs to the Bubbletea program
func (m Model) Read(conn io.ReadWriter, p *tea.Program) {
	s := bufio.NewScanner(conn)
	for s.Scan() {
		line := s.Text()
		log.Println(line)
		msgRcv := client.ParseMessage(line, m.client)

		if msgRcv.Ping {
			m.Write(msgRcv.Content, true)
		}

		// TODO: Move to a dedicated "Handler" module of some sort
		switch msgRcv.NumReply {
		case client.RplHelp:
			m.serverMsg += msgRcv.Content + "\n"
		case client.RplHelpEnd:
			msgRcv.Content = m.serverMsg
			m.serverMsg = "" // reset server message
			p.Send(msgRcv)
		case client.RplNamReply:
		case client.RplMotdStart:
		case client.RplMotd:
			m.serverMsg += msgRcv.Content + "\n"
		case client.RplMotdEnd:
			msgRcv.Content = m.serverMsg
			m.serverMsg = "" // reset server message
			p.Send(msgRcv)
		default:
			// send the received message up to the Bubble Tea Program
			p.Send(msgRcv)
		}
	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
}

func (m *Model) Write(msg string, ping bool) {
	writer := bufio.NewWriter(m.conn)

	if !ping {
		// formats the message into one acceptable by IRC
		msg = m.client.FormatMessage(msg)
	}

	// Just makes for easier formatting, as opposed to WriteString()
	fmt.Fprintf(writer, "%s\r\n", msg)
	writer.Flush()
}

func (m *Model) setContent(text string) {
	// Perform text wrapping before setting the content in the viewport
	wrap := lipgloss.NewStyle().Width(m.viewport.Width)
	m.viewport.SetContent(wrap.Render(text))
}

// toggleBox toggles between the message entry and channels list
func (m *Model) toggleBox() {
	m.mode = (m.mode + 1) % 2
	if m.mode == 0 {
		m.textarea.Blur()
	} else {
		m.textarea.Focus()
	}
}

func channelsToItems(channels []string) []list.Item {
	items := make([]list.Item, len(channels))

	for i := range channels {
		items[i] = item(channels[i])
		i++
	}

	return items
}
