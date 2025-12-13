package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/codoworks/sshto/internal/config"
)

type formField int

const (
	fieldName formField = iota
	fieldHost
	fieldUser
	fieldPort
	fieldKey
	fieldGroup
	fieldCount
)

// FormModel is the bubbletea model for server add/edit form
type FormModel struct {
	inputs   []textinput.Model
	focused  formField
	server   *config.Server
	isEdit   bool
	groups   []config.Group
	done     bool
	canceled bool
	err      error
	warning  string
}

// NewFormModel creates a new form model
func NewFormModel(server *config.Server, groups []config.Group) FormModel {
	inputs := make([]textinput.Model, fieldCount)

	inputs[fieldName] = textinput.New()
	inputs[fieldName].Placeholder = "server-name"
	inputs[fieldName].CharLimit = 64
	inputs[fieldName].Width = 40
	inputs[fieldName].Prompt = "Name: "
	inputs[fieldName].Focus()

	inputs[fieldHost] = textinput.New()
	inputs[fieldHost].Placeholder = "192.168.1.1 or hostname.com"
	inputs[fieldHost].CharLimit = 256
	inputs[fieldHost].Width = 40
	inputs[fieldHost].Prompt = "Host: "

	inputs[fieldUser] = textinput.New()
	inputs[fieldUser].Placeholder = "optional"
	inputs[fieldUser].CharLimit = 64
	inputs[fieldUser].Width = 40
	inputs[fieldUser].Prompt = "User: "

	inputs[fieldPort] = textinput.New()
	inputs[fieldPort].Placeholder = "22"
	inputs[fieldPort].CharLimit = 5
	inputs[fieldPort].Width = 40
	inputs[fieldPort].Prompt = "Port: "

	inputs[fieldKey] = textinput.New()
	inputs[fieldKey].Placeholder = "~/.ssh/id_rsa (optional)"
	inputs[fieldKey].CharLimit = 256
	inputs[fieldKey].Width = 40
	inputs[fieldKey].Prompt = "Key:  "

	inputs[fieldGroup] = textinput.New()
	inputs[fieldGroup].Placeholder = "optional"
	inputs[fieldGroup].CharLimit = 64
	inputs[fieldGroup].Width = 40
	inputs[fieldGroup].Prompt = "Group:"

	isEdit := server != nil
	if server == nil {
		server = &config.Server{}
	} else {
		// Pre-populate fields for edit mode
		inputs[fieldName].SetValue(server.Name)
		inputs[fieldHost].SetValue(server.Host)
		inputs[fieldUser].SetValue(server.User)
		if server.Port != 0 {
			inputs[fieldPort].SetValue(strconv.Itoa(server.Port))
		}
		inputs[fieldKey].SetValue(server.Key)
		inputs[fieldGroup].SetValue(server.Group)
	}

	return FormModel{
		inputs:  inputs,
		focused: fieldName,
		server:  server,
		isEdit:  isEdit,
		groups:  groups,
	}
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit

		case "tab", "down":
			m.focused = (m.focused + 1) % fieldCount
			return m, m.focusField()

		case "shift+tab", "up":
			m.focused = (m.focused - 1 + fieldCount) % fieldCount
			return m, m.focusField()

		case "enter":
			if m.focused == fieldCount-1 {
				// Last field, submit form
				if err := m.validate(); err != nil {
					m.err = err
					return m, nil
				}
				m.buildServer()
				m.done = true
				return m, tea.Quit
			}
			// Move to next field
			m.focused = (m.focused + 1) % fieldCount
			return m, m.focusField()
		}
	}

	// Update the focused input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *FormModel) focusField() tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	return m.inputs[m.focused].Focus()
}

func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return cmd
}

func (m *FormModel) validate() error {
	name := strings.TrimSpace(m.inputs[fieldName].Value())
	if err := config.ValidateName(name); err != nil {
		return err
	}

	host := strings.TrimSpace(m.inputs[fieldHost].Value())
	if err := config.ValidateHost(host); err != nil {
		return err
	}

	portStr := strings.TrimSpace(m.inputs[fieldPort].Value())
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("port must be a valid number (1-65535)")
		}
		if err := config.ValidatePort(port); err != nil {
			return err
		}
	}

	// Check key file (warning only, not error)
	keyPath := strings.TrimSpace(m.inputs[fieldKey].Value())
	if warning, err := config.ValidateKeyFile(keyPath); err != nil {
		return err
	} else if warning != "" {
		m.warning = warning
	}

	return nil
}

func (m *FormModel) buildServer() {
	m.server.Name = strings.TrimSpace(m.inputs[fieldName].Value())
	m.server.Host = strings.TrimSpace(m.inputs[fieldHost].Value())
	m.server.User = strings.TrimSpace(m.inputs[fieldUser].Value())

	portStr := strings.TrimSpace(m.inputs[fieldPort].Value())
	if portStr != "" {
		m.server.Port, _ = strconv.Atoi(portStr)
	}

	m.server.Key = strings.TrimSpace(m.inputs[fieldKey].Value())
	m.server.Group = strings.TrimSpace(m.inputs[fieldGroup].Value())
}

func (m FormModel) View() string {
	var b strings.Builder

	title := "Add Server"
	if m.isEdit {
		title = "Edit Server"
	}
	b.WriteString(TitleStyle.Render(title))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(WarningStyle.Render("Warning: " + m.warning))
		b.WriteString("\n")
	}

	// Show available groups hint
	if len(m.groups) > 0 && m.focused == fieldGroup {
		b.WriteString("\n")
		b.WriteString(DimStyle.Render("Available groups: "))
		var groupNames []string
		for _, g := range m.groups {
			groupNames = append(groupNames, g.Name)
		}
		b.WriteString(DimStyle.Render(strings.Join(groupNames, ", ")))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("tab/shift+tab: navigate • enter: next/submit • esc: cancel"))

	return b.String()
}

// Server returns the built server config
func (m FormModel) Server() *config.Server {
	return m.server
}

// Done returns true if the form was submitted
func (m FormModel) Done() bool {
	return m.done
}

// Canceled returns true if the form was canceled
func (m FormModel) Canceled() bool {
	return m.canceled
}
