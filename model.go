package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	appName       = "windows-shell-cracker"
	visibleRows   = 11
	centerRow     = 6 // 1-based line index of highlighted row
	footerText    = "MouseWheel=Scroll  -  Enter=Run Highlighted  -  Esc=reset  -  Ctrl+C=Exit  "
)

type model struct {
	all            []string
	filtered       []string
	filter         textinput.Model
	selected       int
	lastHighlight  string
	width          int
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "filter: "
	ti.Placeholder = "type to filter commands…"
	ti.CharLimit = 128
	ti.Width = 48
	ti.Focus()

	all := sortedShellCommands()
	last := loadLastCommand()
	filtered := filterCommands(all, "")

	sel := 0
	if last != "" {
		if i := indexOf(filtered, last); i >= 0 {
			sel = i
		}
	}

	return model{
		all:           all,
		filtered:      filtered,
		filter:        ti,
		selected:      sel,
		lastHighlight: highlightedCommand(filtered, sel),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func highlightedCommand(filtered []string, sel int) string {
	if sel < 0 || sel >= len(filtered) {
		return ""
	}
	return filtered[sel]
}

func (m *model) applyFilter(previousFilter string) {
	prevHighlighted := m.lastHighlight
	newFiltered := filterCommands(m.all, m.filter.Value())
	prevInNew := indexOf(newFiltered, prevHighlighted)

	switch {
	case len(m.filter.Value()) > len(previousFilter):
		// Added characters: if current highlight dropped out, jump to first match.
		if prevInNew < 0 {
			m.selected = 0
		} else {
			m.selected = prevInNew
		}
	case len(m.filter.Value()) < len(previousFilter):
		// Removed characters: keep same highlighted text on line 6 when possible.
		if prevInNew >= 0 {
			m.selected = prevInNew
		} else if len(newFiltered) > 0 {
			m.selected = 0
		} else {
			m.selected = 0
		}
	default:
		if prevInNew >= 0 {
			m.selected = prevInNew
		}
	}

	if m.selected >= len(newFiltered) {
		if len(newFiltered) == 0 {
			m.selected = 0
		} else {
			m.selected = len(newFiltered) - 1
		}
	}
	if m.selected < 0 {
		m.selected = 0
	}

	m.filtered = newFiltered
	m.lastHighlight = highlightedCommand(m.filtered, m.selected)
}

func (m *model) moveSelection(delta int) {
	if len(m.filtered) == 0 {
		return
	}
	next := m.selected + delta
	if next < 0 {
		return
	}
	if next >= len(m.filtered) {
		return
	}
	m.selected = next
	m.lastHighlight = highlightedCommand(m.filtered, m.selected)
	_ = saveLastCommand(m.lastHighlight)
}

func (m *model) resetFilter() {
	m.filter.SetValue("")
	m.filtered = filterCommands(m.all, "")
	m.selected = 0
	if len(m.filtered) > 0 {
		m.lastHighlight = m.filtered[0]
		_ = saveLastCommand(m.lastHighlight)
	} else {
		m.lastHighlight = ""
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		filterW := msg.Width - len(appName) - 6
		if filterW < 20 {
			filterW = 20
		}
		m.filter.Width = filterW
		return m, nil

	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp {
			m.moveSelection(-1)
			return m, nil
		}
		if msg.Button == tea.MouseButtonWheelDown {
			m.moveSelection(1)
			return m, nil
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			_ = saveLastCommand(m.lastHighlight)
			return m, tea.Quit
		case "esc":
			m.resetFilter()
			return m, nil
		case "enter":
			_ = launchExplorerShell(m.lastHighlight)
			_ = saveLastCommand(m.lastHighlight)
			return m, nil
		case "up":
			m.moveSelection(-1)
			return m, nil
		case "down":
			m.moveSelection(1)
			return m, nil
		}
	}

	prev := m.filter.Value()
	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	if m.filter.Value() != prev {
		m.applyFilter(prev)
	}

	return m, cmd
}

var (
	styleDarkDim   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	styleLightDim  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	styleNormal    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	styleHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(true)
	styleHeaderL   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	styleHeaderR   = lipgloss.NewStyle()
	styleFooter    = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true)
)

func styleForOffset(offset int) lipgloss.Style {
	switch {
	case offset <= -5:
		return styleDarkDim
	case offset == -4 || offset == -3:
		return styleLightDim
	case offset == -2 || offset == -1:
		return styleNormal
	case offset == 0:
		return styleHighlight
	case offset == 1 || offset == 2:
		return styleNormal
	case offset == 3 || offset == 4:
		return styleLightDim
	default:
		return styleDarkDim
	}
}

func (m model) View() string {
	var b strings.Builder

	headerLeft := styleHeaderL.Render(appName)
	filterLine := m.filter.View()
	header := lipgloss.JoinHorizontal(lipgloss.Top, headerLeft, "  ", styleHeaderR.Render(filterLine))
	b.WriteString(header)
	b.WriteString("\n")

	for line := 1; line <= visibleRows; line++ {
		offset := line - centerRow
		idx := m.selected + offset
		text := ""
		if idx >= 0 && idx < len(m.filtered) {
			text = m.filtered[idx]
		} else {
			text = " "
		}
		b.WriteString(styleForOffset(offset).Render(text))
		b.WriteString("\n")
	}

	footerW := m.width
	if footerW <= 0 {
		footerW = len(footerText)
	}
	b.WriteString(styleFooter.Render(lipgloss.Place(footerW, 1, lipgloss.Left, lipgloss.Bottom, footerText)))

	return b.String()
}

func runTUI() error {
	p := tea.NewProgram(initialModel(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func main() {
	if err := runTUI(); err != nil {
		fmt.Println("error:", err)
	}
}
