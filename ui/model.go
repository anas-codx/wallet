package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	Keys       []string
	ListCursor int
	Err        error
	StatusMsg  string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Password..."
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return Model{
		Choices:   []string{"Add Secret", "List Secrets", "Update Secret", "Delete Secret"},
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
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.State == StateMenu || m.State == StateList {
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
						m.Err = fmt.Errorf("Access Denied: Invalid Password")
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
				case "Add Secret":
					m.State = StateKeyEntry
					m.TextInput.Placeholder = "Secret Name (Key)"
					m.TextInput.EchoMode = textinput.EchoNormal
					m.TextInput.Focus()
				case "List Secrets":
					m.refreshKeys()
					m.State = StateList
					m.ListCursor = 0
				case "Update Secret":
					m.refreshKeys()
					if len(m.Keys) == 0 {
						m.Err = fmt.Errorf("Vault is empty")
						return m, nil
					}
					m.State = StateUpdateKey
					m.ListCursor = 0
				case "Delete Secret":
					m.refreshKeys()
					if len(m.Keys) == 0 {
						m.Err = fmt.Errorf("Vault is empty")
						return m, nil
					}
					m.State = StateDeleteConfirm
					m.ListCursor = 0
				}
				return m, nil

			case StateKeyEntry:
				m.PendingKey = strings.TrimSpace(m.TextInput.Value())
				if m.PendingKey == "" {
					m.Err = fmt.Errorf("Key required")
					return m, nil
				}
				m.State = StateValueEntry
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Secret Value"
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
						m.Err = fmt.Errorf("Clipboard error: %v", err)
					} else {
						m.StatusMsg = fmt.Sprintf("Copied: %s", keyToCopy)
					}
				}

			case StateUpdateKey:
				m.PendingKey = m.Keys[m.ListCursor]
				m.State = StateValueEntry
				m.TextInput.Reset()
				m.TextInput.Placeholder = "New Secret Value"
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
	var content strings.Builder

	// Header
	content.WriteString(HeaderStyle.Render(" "+IconLock+" SECURE VAULT CLI "))
	content.WriteString("\n\n")

	// Status & Error area
	if m.Err != nil {
		content.WriteString(ErrorStyle.Render("! "+m.Err.Error()) + "\n\n")
	}
	if m.StatusMsg != "" {
		content.WriteString(StatusStyle.Render(IconCheck+" "+m.StatusMsg) + "\n\n")
	}

	switch m.State {
	case StatePassword:
		content.WriteString(InputLabelStyle.Render("Vault Authentication"))
		content.WriteString("\n" + m.TextInput.View() + "\n")

	case StateMenu:
		content.WriteString(TitleStyle.Render("Main Menu"))
		content.WriteString("\n")
		for i, choice := range m.Choices {
			if m.Cursor == i {
				content.WriteString(SelectedItemStyle.Render(IconSelected+" "+choice) + "\n")
			} else {
				content.WriteString(ItemStyle.Render(choice) + "\n")
			}
		}

	case StateKeyEntry:
		content.WriteString(InputLabelStyle.Render("New Secret Name"))
		content.WriteString("\n" + m.TextInput.View() + "\n")

	case StateValueEntry:
		content.WriteString(InputLabelStyle.Render(fmt.Sprintf("Secret for [%s]", m.PendingKey)))
		content.WriteString("\n" + m.TextInput.View() + "\n")

	case StateList, StateUpdateKey, StateDeleteConfirm:
		title := "Vault Secrets"
		if m.State == StateUpdateKey {
			title = "Select to Update"
		} else if m.State == StateDeleteConfirm {
			title = "Select to Delete"
		}
		content.WriteString(TitleStyle.Render(title) + "\n")

		if len(m.Keys) == 0 {
			content.WriteString(ItemStyle.Render("Vault is currently empty."))
		} else {
			for i, key := range m.Keys {
				if m.ListCursor == i {
					content.WriteString(SelectedItemStyle.Render(IconSelected+" "+key) + "\n")
				} else {
					content.WriteString(ItemStyle.Render(key) + "\n")
				}
			}
		}
	}

	// Footer
	var help string
	switch m.State {
	case StatePassword:
		help = "enter: unlock vault • ctrl+c: quit"
	case StateMenu:
		help = "↑/↓: navigate • enter: select • q: quit"
	case StateList:
		help = "enter: copy secret • esc: back • q: quit"
	case StateKeyEntry, StateValueEntry:
		help = "enter: confirm • esc: cancel"
	case StateUpdateKey:
		help = "enter: update value • esc: back"
	case StateDeleteConfirm:
		help = "enter: DELETE PERMANENTLY • esc: back"
	}

	content.WriteString(FooterStyle.Render(help))

	return WindowStyle.Render(content.String())
}
