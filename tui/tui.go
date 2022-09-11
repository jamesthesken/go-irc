package tui

/*
TODO:
- Set viewport content based on selected channel in list
- Set message target based on selected channel in list
*/

import (
	"flag"
	"fmt"
	"gopherchatv2/client"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	irc "github.com/fluffle/goirc/client"
)

var wg sync.WaitGroup

var host *string = flag.String("host", "irc.libera.chat", "IRC server")
var channel *string = flag.String("channel", "#test", "IRC channel")

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

	p := *tea.NewProgram(m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	// Connect to IRC
	flag.Parse()

	// create new IRC connection
	c := irc.SimpleClient("GoTest1234124123", "gotest")
	c.EnableStateTracking()
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { conn.Join(*channel) })

	// Set up a handler to notify of disconnect events.
	quit := make(chan bool)
	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	c.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			msgRcv := client.Message{
				Time:    time.Now().Format("3:04PM"),
				Channel: line.Target(),
				Nick:    line.Nick,
				Content: line.Text(),
			}
			p.Send(msgRcv)

		})

	m.ircClient = c

	// connect to server
	if err := c.ConnectTo(*host); err != nil {
		fmt.Printf("Connection error: %s\n", err)
		return
	}

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
	channel     string
	channels    []string
	messages    []string
	ircClient   *irc.Conn
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
	list.SetStatusBarItemName("Channel", "Channels")

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
