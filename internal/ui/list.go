package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/codoworks/sshto/internal/config"
)

// ServerItem represents a server in the list
type ServerItem struct {
	Server     config.Server
	GroupColor string
}

func (s ServerItem) FilterValue() string {
	return s.Server.Name + " " + s.Server.Host + " " + s.Server.Group
}

func (s ServerItem) Title() string {
	return s.Server.Name
}

func (s ServerItem) Description() string {
	desc := s.Server.Host
	if s.Server.User != "" {
		desc = s.Server.User + "@" + desc
	}
	if s.Server.Port != 0 && s.Server.Port != 22 {
		desc = fmt.Sprintf("%s:%d", desc, s.Server.Port)
	}
	return desc
}

// ServerItemDelegate handles rendering of server items
type ServerItemDelegate struct {
	groups map[string]*config.Group
}

func NewServerItemDelegate(groups []config.Group) ServerItemDelegate {
	m := make(map[string]*config.Group)
	for i := range groups {
		m[groups[i].Name] = &groups[i]
	}
	return ServerItemDelegate{groups: m}
}

func (d ServerItemDelegate) Height() int                             { return 2 }
func (d ServerItemDelegate) Spacing() int                            { return 0 }
func (d ServerItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ServerItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s, ok := item.(ServerItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	// Build title line
	title := s.Server.Name
	if s.Server.Group != "" {
		color := "gray"
		if g, ok := d.groups[s.Server.Group]; ok && g.Color != "" {
			color = g.Color
		}
		title = GroupTag(s.Server.Group, color) + title
	}

	// Build description
	desc := s.Description()

	// Apply styles
	if isSelected {
		title = SelectedItemStyle.Render("> " + title)
		desc = SelectedItemStyle.Copy().Bold(false).Render("  " + desc)
	} else {
		title = ItemStyle.Render("  " + title)
		desc = DimStyle.Render("    " + desc)
	}

	fmt.Fprintf(w, "%s\n%s\n", title, desc)
}

// ListModel is the bubbletea model for server selection
type ListModel struct {
	list     list.Model
	selected *config.Server
	quitting bool
}

// NewListModel creates a new list model
func NewListModel(servers []config.Server, groups []config.Group) ListModel {
	items := make([]list.Item, len(servers))
	for i, s := range servers {
		items[i] = ServerItem{Server: s}
	}

	delegate := NewServerItemDelegate(groups)
	l := list.New(items, delegate, 80, 20)
	l.Title = "Select a server"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle
	l.Styles.HelpStyle = HelpStyle

	return ListModel{list: l}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(ServerItem); ok {
				m.selected = &item.Server
				m.quitting = true
				return m, tea.Quit
			}
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}

// Selected returns the selected server, if any
func (m ListModel) Selected() *config.Server {
	return m.selected
}

// FilterByGroup returns a new list filtered by group
func FilterByGroup(servers []config.Server, group string) []config.Server {
	if group == "" {
		return servers
	}
	var filtered []config.Server
	for _, s := range servers {
		if strings.EqualFold(s.Group, group) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
