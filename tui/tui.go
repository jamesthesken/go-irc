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
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var wg sync.WaitGroup

func StartTea() {

	m := initialModel()

	wg.Add(1)
	// Connect to IRC
	client := &client.Client{}
	conn, err := client.Connect("irc.libera.chat:6697")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	go m.Read(conn)

	p := tea.NewProgram(m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type tickMsg struct{}
type errMsg error

type Model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel() *Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 10)
	vp.SetContent(`Welcome to IRC!`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &Model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

// Init() is the first function called by BubbleTea.
func (m Model) Init() tea.Cmd {

	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// Implement different tea messages sent by the clients.
	// i.e., Message interface
	case RcvMessage:
		// this does not work:
		m.messages = append(m.messages, m.senderStyle.Render("Server: ")+string(msg.content))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		fmt.Print(msg.content) // but this does ??
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit), msg.String() == "ctrl+c":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			timeStamp := time.Now()
			m.messages = append(m.messages, m.senderStyle.Render("< You > "+timeStamp.Format("3:04PM: "))+m.textarea.Value())
			print(m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

type RcvMessage struct {
	content string
}

func (model *Model) Read(conn io.ReadWriter) {
	s := bufio.NewScanner(conn)
	for s.Scan() {
		line := s.Text()
		msgRcv := fmt.Sprintf(client.ParseMessage(line), "\n")
		model.Update(RcvMessage{content: msgRcv})
	}
	if s.Err() != nil {
		log.Fatalf("Error occured: %s", s.Err())
	}
}
