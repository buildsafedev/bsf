package search

import (
	"context"
	"time"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

var (
	frameHeight, frameWidth int
)

type pkgitem struct {
	name string
}

func (i pkgitem) Title() string       { return i.name }
func (i pkgitem) Description() string { return i.name }
func (i pkgitem) FilterValue() string {
	return i.name
}

func convLPR2Items(packages *buildsafev1.ListPackagesResponse) []list.Item {
	items := make([]list.Item, 0, len(packages.Packages))
	for _, name := range packages.Packages {
		items = append(items, pkgitem{name: name})
	}

	return items
}

type (
	errMsg struct{ error }
)

type mode int

const (
	modeSearch mode = iota
	modeVersion
	modeOption
)

var (
	currentMode mode = modeSearch
)

// Model the search model definition
type Model struct {
	// mode       mode
	searchList list.Model
	vlModel    *versionListModel
	// input      textinput.Model
	quitting bool
}

// InitSearch initialize the search model for your program
func InitSearch(items []list.Item) tea.Model {
	currentMode = modeSearch
	m := Model{searchList: list.New(items, list.NewDefaultDelegate(), 8, 8)}

	if frameHeight != 0 || frameWidth != 0 {
		m.searchList.SetSize(frameWidth, frameHeight)
	}
	m.searchList.Title = "packages"
	m.searchList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			KeyMap.Enter,
			KeyMap.Back,
			KeyMap.Quit,
		}
	}
	return m
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, KeyMap.Quit) {
			return m, tea.Quit
		}

		if key.Matches(msg, KeyMap.Enter) {
			item := m.searchList.SelectedItem()
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			pvr, err := sc.FetchPackages(ctx, &buildsafev1.FetchPackagesRequest{
				Name: item.FilterValue(),
			})
			if err != nil {
				return m, tea.Quit
			}

			m.vlModel = initVersionTable(item.FilterValue(), m.searchList, pvr)
			nm, cmd := m.vlModel.Update(msg)
			return nm, cmd
		}
	case tea.WindowSizeMsg:
		frameHeight, frameWidth = styles.DocStyle.GetFrameSize()
		m.searchList.SetSize(msg.Width-frameHeight, msg.Height-frameWidth)
		frameHeight = msg.Height - frameWidth
		frameWidth = msg.Width - frameHeight
	}
	m.searchList, cmd = m.searchList.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return styles.DocStyle.Render(m.searchList.View() + "\n")
}

type keymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Space key.Binding
	Up    key.Binding
	Down  key.Binding
}

// KeyMap reusable key mappings shared across models
var KeyMap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),

	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
}
