package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/codoworks/sshto/internal/config"
)

func TestNewFormModelAdd(t *testing.T) {
	groups := []config.Group{
		{Name: "production", Color: "red"},
	}

	model := NewFormModel(nil, groups)

	if model.isEdit {
		t.Error("isEdit should be false for new server")
	}
	if model.server == nil {
		t.Error("server should not be nil")
	}
	if len(model.inputs) != int(fieldCount) {
		t.Errorf("inputs length = %d, want %d", len(model.inputs), fieldCount)
	}
}

func TestNewFormModelEdit(t *testing.T) {
	server := &config.Server{
		Name:  "test-server",
		Host:  "192.168.1.1",
		User:  "admin",
		Port:  2222,
		Key:   "~/.ssh/key",
		Group: "production",
	}
	groups := []config.Group{
		{Name: "production", Color: "red"},
	}

	model := NewFormModel(server, groups)

	if !model.isEdit {
		t.Error("isEdit should be true for existing server")
	}

	// Verify fields are pre-populated
	if model.inputs[fieldName].Value() != "test-server" {
		t.Errorf("Name field = %q, want %q", model.inputs[fieldName].Value(), "test-server")
	}
	if model.inputs[fieldHost].Value() != "192.168.1.1" {
		t.Errorf("Host field = %q, want %q", model.inputs[fieldHost].Value(), "192.168.1.1")
	}
	if model.inputs[fieldUser].Value() != "admin" {
		t.Errorf("User field = %q, want %q", model.inputs[fieldUser].Value(), "admin")
	}
	if model.inputs[fieldPort].Value() != "2222" {
		t.Errorf("Port field = %q, want %q", model.inputs[fieldPort].Value(), "2222")
	}
	if model.inputs[fieldKey].Value() != "~/.ssh/key" {
		t.Errorf("Key field = %q, want %q", model.inputs[fieldKey].Value(), "~/.ssh/key")
	}
	if model.inputs[fieldGroup].Value() != "production" {
		t.Errorf("Group field = %q, want %q", model.inputs[fieldGroup].Value(), "production")
	}
}

func TestNewFormModelEditZeroPort(t *testing.T) {
	server := &config.Server{
		Name: "test",
		Host: "localhost",
		Port: 0,
	}

	model := NewFormModel(server, nil)

	// Port should be empty for zero port
	if model.inputs[fieldPort].Value() != "" {
		t.Errorf("Port field = %q, want empty for zero port", model.inputs[fieldPort].Value())
	}
}

func TestFormModelInit(t *testing.T) {
	model := NewFormModel(nil, nil)
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a command for blinking cursor")
	}
}

func TestFormModelDone(t *testing.T) {
	model := FormModel{done: false}
	if model.Done() {
		t.Error("Done() should return false initially")
	}

	model.done = true
	if !model.Done() {
		t.Error("Done() should return true after form submitted")
	}
}

func TestFormModelCanceled(t *testing.T) {
	model := FormModel{canceled: false}
	if model.Canceled() {
		t.Error("Canceled() should return false initially")
	}

	model.canceled = true
	if !model.Canceled() {
		t.Error("Canceled() should return true after form canceled")
	}
}

func TestFormModelServer(t *testing.T) {
	server := &config.Server{Name: "test"}
	model := FormModel{server: server}

	if model.Server() != server {
		t.Error("Server() should return the server pointer")
	}
}

func TestFormModelView(t *testing.T) {
	model := NewFormModel(nil, nil)
	view := model.View()

	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestFormModelViewEdit(t *testing.T) {
	server := &config.Server{Name: "test", Host: "localhost"}
	model := NewFormModel(server, nil)
	view := model.View()

	if view == "" {
		t.Error("View() returned empty string for edit mode")
	}
}

func TestFormModelViewWithError(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.err = config.ValidateName("") // This will set an error
	view := model.View()

	if view == "" {
		t.Error("View() returned empty string with error")
	}
}

func TestFormModelViewWithWarning(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.warning = "test warning"
	view := model.View()

	if view == "" {
		t.Error("View() returned empty string with warning")
	}
}

func TestFormModelViewWithGroups(t *testing.T) {
	groups := []config.Group{
		{Name: "production"},
		{Name: "staging"},
	}
	model := NewFormModel(nil, groups)
	model.focused = fieldGroup

	view := model.View()
	if view == "" {
		t.Error("View() returned empty string with groups")
	}
}

func TestFormValidateEmpty(t *testing.T) {
	model := NewFormModel(nil, nil)

	err := model.validate()
	if err == nil {
		t.Error("validate() should return error for empty name")
	}
}

func TestFormValidateValidInput(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test-server")
	model.inputs[fieldHost].SetValue("192.168.1.1")
	model.inputs[fieldPort].SetValue("22")

	err := model.validate()
	if err != nil {
		t.Errorf("validate() error = %v, want nil", err)
	}
}

func TestFormValidateInvalidHost(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test-server")
	model.inputs[fieldHost].SetValue("-invalid")

	err := model.validate()
	if err == nil {
		t.Error("validate() should return error for invalid host")
	}
}

func TestFormValidateInvalidPort(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test-server")
	model.inputs[fieldHost].SetValue("192.168.1.1")
	model.inputs[fieldPort].SetValue("invalid")

	err := model.validate()
	if err == nil {
		t.Error("validate() should return error for invalid port")
	}
}

func TestFormValidatePortOutOfRange(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test-server")
	model.inputs[fieldHost].SetValue("192.168.1.1")
	model.inputs[fieldPort].SetValue("70000")

	err := model.validate()
	if err == nil {
		t.Error("validate() should return error for port out of range")
	}
}

func TestFormBuildServer(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("  test-server  ")
	model.inputs[fieldHost].SetValue("  192.168.1.1  ")
	model.inputs[fieldUser].SetValue("  admin  ")
	model.inputs[fieldPort].SetValue("2222")
	model.inputs[fieldKey].SetValue("  ~/.ssh/key  ")
	model.inputs[fieldGroup].SetValue("  production  ")

	model.buildServer()

	s := model.server
	if s.Name != "test-server" {
		t.Errorf("Name = %q, want %q (trimmed)", s.Name, "test-server")
	}
	if s.Host != "192.168.1.1" {
		t.Errorf("Host = %q, want %q (trimmed)", s.Host, "192.168.1.1")
	}
	if s.User != "admin" {
		t.Errorf("User = %q, want %q (trimmed)", s.User, "admin")
	}
	if s.Port != 2222 {
		t.Errorf("Port = %d, want %d", s.Port, 2222)
	}
	if s.Key != "~/.ssh/key" {
		t.Errorf("Key = %q, want %q (trimmed)", s.Key, "~/.ssh/key")
	}
	if s.Group != "production" {
		t.Errorf("Group = %q, want %q (trimmed)", s.Group, "production")
	}
}

func TestFormBuildServerEmptyPort(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test")
	model.inputs[fieldHost].SetValue("localhost")
	model.inputs[fieldPort].SetValue("")

	model.buildServer()

	if model.server.Port != 0 {
		t.Errorf("Port = %d, want 0 for empty port", model.server.Port)
	}
}

func TestFieldConstants(t *testing.T) {
	// Verify field constants are in expected order
	if fieldName != 0 {
		t.Error("fieldName should be 0")
	}
	if fieldHost != 1 {
		t.Error("fieldHost should be 1")
	}
	if fieldUser != 2 {
		t.Error("fieldUser should be 2")
	}
	if fieldPort != 3 {
		t.Error("fieldPort should be 3")
	}
	if fieldKey != 4 {
		t.Error("fieldKey should be 4")
	}
	if fieldGroup != 5 {
		t.Error("fieldGroup should be 5")
	}
	if fieldCount != 6 {
		t.Error("fieldCount should be 6")
	}
}

func TestFormModelUpdateEsc(t *testing.T) {
	model := NewFormModel(nil, nil)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Error("Esc should return quit command")
	}
	m := newModel.(FormModel)
	if !m.Canceled() {
		t.Error("Model should be canceled after Esc")
	}
}

func TestFormModelUpdateCtrlC(t *testing.T) {
	model := NewFormModel(nil, nil)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	m := newModel.(FormModel)
	if !m.Canceled() {
		t.Error("Model should be canceled after Ctrl+C")
	}
}

func TestFormModelUpdateTab(t *testing.T) {
	model := NewFormModel(nil, nil)

	// Initial focus should be on name field
	if model.focused != fieldName {
		t.Errorf("Initial focus = %d, want %d", model.focused, fieldName)
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	m := newModel.(FormModel)
	if m.focused != fieldHost {
		t.Errorf("After Tab focus = %d, want %d", m.focused, fieldHost)
	}
}

func TestFormModelUpdateShiftTab(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.focused = fieldHost

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m := newModel.(FormModel)
	if m.focused != fieldName {
		t.Errorf("After Shift+Tab focus = %d, want %d", m.focused, fieldName)
	}
}

func TestFormModelUpdateDown(t *testing.T) {
	model := NewFormModel(nil, nil)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := newModel.(FormModel)
	if m.focused != fieldHost {
		t.Errorf("After Down focus = %d, want %d", m.focused, fieldHost)
	}
}

func TestFormModelUpdateUp(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.focused = fieldHost

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m := newModel.(FormModel)
	if m.focused != fieldName {
		t.Errorf("After Up focus = %d, want %d", m.focused, fieldName)
	}
}

func TestFormModelUpdateEnterNextField(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test")
	model.inputs[fieldHost].SetValue("localhost")

	// Enter on non-last field should move to next
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := newModel.(FormModel)
	if m.focused != fieldHost {
		t.Errorf("After Enter on name field, focus = %d, want %d", m.focused, fieldHost)
	}
}

func TestFormModelUpdateEnterSubmit(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test")
	model.inputs[fieldHost].SetValue("localhost")
	model.focused = fieldGroup // Last field

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enter on last field should return quit command")
	}
	m := newModel.(FormModel)
	if !m.Done() {
		t.Error("Model should be done after valid submit")
	}
}

func TestFormModelUpdateEnterSubmitInvalid(t *testing.T) {
	model := NewFormModel(nil, nil)
	// Leave name empty (invalid)
	model.focused = fieldGroup // Last field

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("Enter with invalid data should not quit")
	}
	m := newModel.(FormModel)
	if m.Done() {
		t.Error("Model should not be done with invalid data")
	}
	if m.err == nil {
		t.Error("Model should have error set")
	}
}

func TestFormModelUpdateInputs(t *testing.T) {
	model := NewFormModel(nil, nil)

	// Type a character
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m := newModel.(FormModel)
	if m.inputs[fieldName].Value() != "a" {
		t.Errorf("Input value = %q, want %q", m.inputs[fieldName].Value(), "a")
	}
}

func TestFormValidateKeyFileWarning(t *testing.T) {
	model := NewFormModel(nil, nil)
	model.inputs[fieldName].SetValue("test")
	model.inputs[fieldHost].SetValue("localhost")
	model.inputs[fieldKey].SetValue("/nonexistent/key/file")

	err := model.validate()
	if err != nil {
		t.Errorf("validate() error = %v, want nil (warning only)", err)
	}
	if model.warning == "" {
		t.Error("validate() should set warning for non-existent key file")
	}
}
