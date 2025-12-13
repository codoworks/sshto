package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/codoworks/sshto/internal/config"
)

func TestServerItemFilterValue(t *testing.T) {
	item := ServerItem{
		Server: config.Server{Name: "web1", Host: "192.168.1.1", Group: "production"},
	}

	expected := "web1 192.168.1.1 production"
	if item.FilterValue() != expected {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), expected)
	}
}

func TestServerItemTitle(t *testing.T) {
	item := ServerItem{
		Server: config.Server{Name: "web1"},
	}

	if item.Title() != "web1" {
		t.Errorf("Title() = %q, want %q", item.Title(), "web1")
	}
}

func TestServerItemDescription(t *testing.T) {
	tests := []struct {
		name     string
		server   config.Server
		expected string
	}{
		{
			"host only",
			config.Server{Host: "192.168.1.1"},
			"192.168.1.1",
		},
		{
			"user and host",
			config.Server{Host: "192.168.1.1", User: "admin"},
			"admin@192.168.1.1",
		},
		{
			"user, host and port",
			config.Server{Host: "192.168.1.1", User: "admin", Port: 2222},
			"admin@192.168.1.1:2222",
		},
		{
			"standard port",
			config.Server{Host: "192.168.1.1", Port: 22},
			"192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := ServerItem{Server: tt.server}
			if item.Description() != tt.expected {
				t.Errorf("Description() = %q, want %q", item.Description(), tt.expected)
			}
		})
	}
}

func TestNewServerItemDelegate(t *testing.T) {
	groups := []config.Group{
		{Name: "production", Color: "red"},
		{Name: "staging", Color: "yellow"},
	}

	delegate := NewServerItemDelegate(groups)
	if delegate.groups == nil {
		t.Error("delegate.groups is nil")
	}
	if len(delegate.groups) != 2 {
		t.Errorf("delegate.groups length = %d, want 2", len(delegate.groups))
	}
}

func TestServerItemDelegateHeight(t *testing.T) {
	delegate := ServerItemDelegate{}
	if delegate.Height() != 2 {
		t.Errorf("Height() = %d, want 2", delegate.Height())
	}
}

func TestServerItemDelegateSpacing(t *testing.T) {
	delegate := ServerItemDelegate{}
	if delegate.Spacing() != 0 {
		t.Errorf("Spacing() = %d, want 0", delegate.Spacing())
	}
}

func TestNewListModel(t *testing.T) {
	servers := []config.Server{
		{Name: "web1", Host: "192.168.1.1", Group: "production"},
		{Name: "db1", Host: "192.168.1.2", Group: "staging"},
	}
	groups := []config.Group{
		{Name: "production", Color: "red"},
	}

	model := NewListModel(servers, groups)

	// Verify model is properly initialized
	if model.Selected() != nil {
		t.Error("New model should have no selection")
	}
}

func TestListModelInit(t *testing.T) {
	model := ListModel{}
	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestFilterByGroup(t *testing.T) {
	servers := []config.Server{
		{Name: "web1", Host: "192.168.1.1", Group: "production"},
		{Name: "web2", Host: "192.168.1.2", Group: "production"},
		{Name: "db1", Host: "192.168.1.3", Group: "staging"},
		{Name: "dev1", Host: "192.168.1.4", Group: "dev"},
	}

	tests := []struct {
		name     string
		group    string
		expected int
	}{
		{"production group", "production", 2},
		{"staging group", "staging", 1},
		{"dev group", "dev", 1},
		{"empty group (all)", "", 4},
		{"non-existent group", "nonexistent", 0},
		{"case insensitive", "PRODUCTION", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterByGroup(servers, tt.group)
			if len(filtered) != tt.expected {
				t.Errorf("FilterByGroup(%q) returned %d servers, want %d", tt.group, len(filtered), tt.expected)
			}
		})
	}
}

func TestListModelSelected(t *testing.T) {
	model := ListModel{
		selected: &config.Server{Name: "test"},
	}

	selected := model.Selected()
	if selected == nil {
		t.Error("Selected() returned nil")
	}
	if selected.Name != "test" {
		t.Errorf("Selected().Name = %q, want %q", selected.Name, "test")
	}
}

func TestListModelSelectedNil(t *testing.T) {
	model := ListModel{}
	if model.Selected() != nil {
		t.Error("Selected() should return nil for new model")
	}
}

func TestListModelView(t *testing.T) {
	servers := []config.Server{{Name: "test", Host: "localhost"}}
	model := NewListModel(servers, nil)

	view := model.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestListModelViewQuitting(t *testing.T) {
	model := ListModel{quitting: true}
	view := model.View()
	if view != "" {
		t.Errorf("View() when quitting = %q, want empty", view)
	}
}

func TestListModelUpdate(t *testing.T) {
	servers := []config.Server{{Name: "test", Host: "localhost"}}
	model := NewListModel(servers, nil)

	// Test window size message
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	if newModel == nil {
		t.Error("Update() returned nil model")
	}

	// Test quit key
	m := newModel.(ListModel)
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	m = newModel.(ListModel)
	if !m.quitting {
		t.Error("Model should be quitting after Ctrl+C")
	}
}

func TestListModelUpdateEnter(t *testing.T) {
	servers := []config.Server{{Name: "test", Host: "localhost"}}
	model := NewListModel(servers, nil)

	// Test enter key
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enter should return quit command when item selected")
	}
	m := newModel.(ListModel)
	if m.Selected() == nil {
		t.Error("Should have selected server after Enter")
	}
}

func TestListModelUpdateEsc(t *testing.T) {
	servers := []config.Server{{Name: "test", Host: "localhost"}}
	model := NewListModel(servers, nil)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Error("Esc should return quit command")
	}
	m := newModel.(ListModel)
	if !m.quitting {
		t.Error("Model should be quitting after Esc")
	}
}

func TestListModelUpdateQ(t *testing.T) {
	servers := []config.Server{{Name: "test", Host: "localhost"}}
	model := NewListModel(servers, nil)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("'q' should return quit command")
	}
	m := newModel.(ListModel)
	if !m.quitting {
		t.Error("Model should be quitting after 'q'")
	}
}

func TestServerItemDelegateUpdate(t *testing.T) {
	delegate := ServerItemDelegate{}
	cmd := delegate.Update(nil, nil)
	if cmd != nil {
		t.Error("Delegate Update should return nil")
	}
}

func TestServerItemDelegateRender(t *testing.T) {
	groups := []config.Group{{Name: "production", Color: "red"}}
	_ = NewServerItemDelegate(groups)

	servers := []config.Server{
		{Name: "test", Host: "localhost", Group: "production"},
	}
	model := NewListModel(servers, groups)

	// We can't easily test the render output, but we can verify it doesn't panic
	_ = model.View()
}

func TestServerItemDelegateRenderNoGroup(t *testing.T) {
	_ = NewServerItemDelegate(nil)

	servers := []config.Server{
		{Name: "test", Host: "localhost"},
	}
	model := NewListModel(servers, nil)

	// Verify render doesn't panic
	_ = model.View()
}
