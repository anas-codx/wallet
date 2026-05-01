package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateMenu State = iota
	StateInput
)

type Model struct {
	Choices   []string
	Cursor    int
	Selected  map[int]struct{}
	State     State
	TextInput textinput.Model
	Err       error
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Amount"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		Choices:   []string{"Check Balance", "Add Income", "Add Expense", "History"},
		Selected:  make(map[int]struct{}),
		State:     StateMenu,
		TextInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.State == StateMenu && m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.State == StateMenu && m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter":
			if m.State == StateMenu {
				choice := m.Choices[m.Cursor]
				if choice == "Add Income" || choice == "Add Expense" {
					m.State = StateInput
					m.TextInput.Focus()
					return m, nil
				}
				if _, ok := m.Selected[m.Cursor]; ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = struct{}{}
				}
			} else if m.State == StateInput {
				m.State = StateMenu
				m.TextInput.Blur()
				m.TextInput.Reset()
			}

		case "esc":
			if m.State == StateInput {
				m.State = StateMenu
				m.TextInput.Blur()
				m.TextInput.Reset()
			}
		}
	}

	if m.State == StateInput {
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render("Wallet CLI"))
	s.WriteString("\n\n")

	if m.State == StateMenu {
		for i, choice := range m.Choices {
			cursor := " "
			var style lipgloss.Style

			if m.Cursor == i {
				cursor = ">"
				style = SelectedItemStyle
			} else {
				style = ItemStyle
			}

			checked := " "
			if _, ok := m.Selected[i]; ok {
				checked = CheckMark.String()
			}

			s.WriteString(style.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, choice)))
			s.WriteString("\n")
		}

		s.WriteString(HelpStyle.Render("\n↑/↓: navigate • enter: select • q: quit"))
	} else if m.State == StateInput {
		s.WriteString(fmt.Sprintf(
			"Enter amount for %s:\n\n%s\n\n%s",
			m.Choices[m.Cursor],
			m.TextInput.View(),
			HelpStyle.Render("esc: back • enter: submit"),
		))
	}

	return s.String()
}
