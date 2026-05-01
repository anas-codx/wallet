package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codingdestro/wallet-go/pkg/utils"
)

type State int

const (
	StatePassword State = iota
	StateMenu
	StateKeyEntry
	StateValueEntry
	StateList
	StateUpdateKey
	StateDeleteConfirm
)

type Model struct {
	Choices    []string
	Cursor     int
	State      State
	TextInput  textinput.Model
	Data       map[string]string
	Password   string
	Filename   string
	PendingKey string
	Keys       []string // Sorted keys for listing
	ListCursor int
	Err        error
	StatusMsg  string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return Model{
		Choices:   []string{"Add Key", "List Keys", "Update Key", "Delete Key"},
		State:     StatePassword,
		TextInput: ti,
		Data:      make(map[string]string),
		Filename:  "wallet.enc",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) refreshKeys() {
	m.Keys = []string{}
	for k := range m.Data {
		m.Keys = append(m.Keys, k)
	}
	sort.Strings(m.Keys)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.State == StateMenu {
				return m, tea.Quit
			}
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case "up", "k":
			m.StatusMsg = ""
			if m.State == StateMenu {
				if m.Cursor > 0 {
					m.Cursor--
				}
			} else if m.State == StateList || m.State == StateUpdateKey || m.State == StateDeleteConfirm {
				if m.ListCursor > 0 {
					m.ListCursor--
				}
			}

		case "down", "j":
			m.StatusMsg = ""
			if m.State == StateMenu {
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}
			} else if m.State == StateList || m.State == StateUpdateKey || m.State == StateDeleteConfirm {
				if m.ListCursor < len(m.Keys)-1 {
					m.ListCursor++
				}
			}

		case "enter":
			m.StatusMsg = ""
			switch m.State {
			case StatePassword:
				m.Password = m.TextInput.Value()
				if utils.FileExists(m.Filename) {
					err := utils.LoadEncryptedJSON(m.Filename, m.Password, &m.Data)
					if err != nil {
						m.Err = fmt.Errorf("wrong password or corrupt file")
						m.TextInput.Reset()
						return m, nil
					}
				} else {
					m.Data = make(map[string]string)
				}
				m.refreshKeys()
				m.State = StateMenu
				m.TextInput.Reset()
				m.TextInput.Blur()
				m.Err = nil

			case StateMenu:
				choice := m.Choices[m.Cursor]
				switch choice {
				case "Add Key":
					m.State = StateKeyEntry
					m.TextInput.Placeholder = "Key"
					m.TextInput.EchoMode = textinput.EchoNormal
					m.TextInput.Focus()
				case "List Keys":
					m.refreshKeys()
					m.State = StateList
					m.ListCursor = 0
				case "Update Key":
					m.refreshKeys()
					if len(m.Keys) == 0 {
						m.Err = fmt.Errorf("no keys to update")
						return m, nil
					}
					m.State = StateUpdateKey
					m.ListCursor = 0
				case "Delete Key":
					m.refreshKeys()
					if len(m.Keys) == 0 {
						m.Err = fmt.Errorf("no keys to delete")
						return m, nil
					}
					m.State = StateDeleteConfirm
					m.ListCursor = 0
				}
				return m, nil

			case StateKeyEntry:
				m.PendingKey = strings.TrimSpace(m.TextInput.Value())
				if m.PendingKey == "" {
					m.Err = fmt.Errorf("key cannot be empty")
					return m, nil
				}
				m.State = StateValueEntry
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Value"
				m.TextInput.EchoMode = textinput.EchoPassword
				m.TextInput.EchoCharacter = '•'
				m.Err = nil

			case StateValueEntry:
				value := m.TextInput.Value()
				m.Data[m.PendingKey] = value
				err := utils.SaveEncryptedJSON(m.Filename, m.Password, m.Data)
				if err != nil {
					m.Err = err
				}
				m.refreshKeys()
				m.State = StateMenu
				m.TextInput.Blur()
				m.TextInput.Reset()
				m.PendingKey = ""
				m.Err = nil

			case StateList:
				if len(m.Keys) > 0 {
					keyToCopy := m.Keys[m.ListCursor]
					valToCopy := m.Data[keyToCopy]
					err := clipboard.WriteAll(valToCopy)
					if err != nil {
						m.Err = fmt.Errorf("failed to copy to clipboard: %v", err)
					} else {
						m.StatusMsg = fmt.Sprintf("Copied value for [%s] to clipboard!", keyToCopy)
					}
				}

			case StateUpdateKey:
				m.PendingKey = m.Keys[m.ListCursor]
				m.State = StateValueEntry
				m.TextInput.Reset()
				m.TextInput.Placeholder = "New Value"
				m.TextInput.EchoMode = textinput.EchoPassword
				m.TextInput.EchoCharacter = '•'
				m.TextInput.Focus()
				m.Err = nil

			case StateDeleteConfirm:
				keyToDelete := m.Keys[m.ListCursor]
				delete(m.Data, keyToDelete)
				err := utils.SaveEncryptedJSON(m.Filename, m.Password, m.Data)
				if err != nil {
					m.Err = err
				}
				m.refreshKeys()
				m.State = StateMenu
				m.Err = nil
			}

		case "esc":
			m.StatusMsg = ""
			if m.State != StateMenu && m.State != StatePassword {
				m.State = StateMenu
				m.TextInput.Blur()
				m.TextInput.Reset()
				m.PendingKey = ""
				m.Err = nil
			}
		}
	}

	if m.State == StatePassword || m.State == StateKeyEntry || m.State == StateValueEntry {
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render("Wallet CLI"))
	s.WriteString("\n\n")

	if m.Err != nil {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v\n\n", m.Err)))
	}

	if m.StatusMsg != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(fmt.Sprintf("%s\n\n", m.StatusMsg)))
	}

	switch m.State {
	case StatePassword:
		s.WriteString("Enter Wallet Password:\n\n")
		s.WriteString(m.TextInput.View())
		s.WriteString(HelpStyle.Render("\n\npress enter to unlock/create"))

	case StateMenu:
		for i, choice := range m.Choices {
			cursor := " "
			var style lipgloss.Style

			if m.Cursor == i {
				cursor = ">"
				style = SelectedItemStyle
			} else {
				style = ItemStyle
			}

			s.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, choice)))
			s.WriteString("\n")
		}
		s.WriteString(HelpStyle.Render("\n↑/↓: navigate • enter: select • q: quit"))

	case StateKeyEntry:
		s.WriteString(fmt.Sprintf(
			"Enter Key:\n\n%s\n\n%s",
			m.TextInput.View(),
			HelpStyle.Render("esc: back • enter: next"),
		))

	case StateValueEntry:
		s.WriteString(fmt.Sprintf(
			"Enter Value for [%s]:\n\n%s\n\n%s",
			m.PendingKey,
			m.TextInput.View(),
			HelpStyle.Render("esc: back • enter: submit"),
		))

	case StateList, StateUpdateKey, StateDeleteConfirm:
		title := "List of Keys:"
		help := "esc: back"
		if m.State == StateUpdateKey {
			title = "Select Key to Update:"
			help = "↑/↓: navigate • enter: select • esc: back"
		} else if m.State == StateDeleteConfirm {
			title = "Select Key to Delete:"
			help = "↑/↓: navigate • enter: DELETE • esc: back"
		}

		s.WriteString(title + "\n\n")
		if len(m.Keys) == 0 {
			s.WriteString(ItemStyle.Render("No keys found."))
		} else {
			for i, key := range m.Keys {
				cursor := " "
				style := ItemStyle
				if (m.State == StateList || m.State == StateUpdateKey || m.State == StateDeleteConfirm) && m.ListCursor == i {
					cursor = ">"
					style = SelectedItemStyle
				}
				s.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, key)))
				s.WriteString("\n")
			}
		}
		s.WriteString(HelpStyle.Render("\n" + help))
	}

	return s.String()
}
