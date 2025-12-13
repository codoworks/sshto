package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "sshto", "config.yaml")
	if path != expected {
		t.Errorf("DefaultPath() = %q, want %q", path, expected)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Load() error = %v, want nil for non-existent file", err)
	}
	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
	if cfg.Defaults.Port != 22 {
		t.Errorf("Default port = %d, want 22", cfg.Defaults.Port)
	}
	if len(cfg.Servers) != 0 {
		t.Errorf("Servers = %d, want 0", len(cfg.Servers))
	}
}

func TestLoadAndSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create initial config
	cfg := &Config{
		path: configPath,
		Groups: []Group{
			{Name: "production", Color: "red"},
			{Name: "staging", Color: "yellow"},
		},
		Servers: []Server{
			{Name: "web1", Host: "192.168.1.1", User: "admin", Port: 22, Group: "production"},
			{Name: "db1", Host: "192.168.1.2", User: "root", Port: 2222, Key: "~/.ssh/db_key", Group: "staging"},
		},
		Defaults: Defaults{User: "deploy", Port: 22},
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify
	if len(loaded.Groups) != 2 {
		t.Errorf("Groups count = %d, want 2", len(loaded.Groups))
	}
	if len(loaded.Servers) != 2 {
		t.Errorf("Servers count = %d, want 2", len(loaded.Servers))
	}
	if loaded.Defaults.User != "deploy" {
		t.Errorf("Defaults.User = %q, want %q", loaded.Defaults.User, "deploy")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("Failed to write invalid YAML: %v", err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Error("Load() should return error for invalid YAML")
	}
}

func TestConfigPath(t *testing.T) {
	cfg := &Config{path: "/test/path/config.yaml"}
	if cfg.Path() != "/test/path/config.yaml" {
		t.Errorf("Path() = %q, want %q", cfg.Path(), "/test/path/config.yaml")
	}
}

func TestFindServer(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "web1", Host: "192.168.1.1"},
			{Name: "db1", Host: "192.168.1.2"},
		},
	}

	// Find existing server
	server, err := cfg.FindServer("web1")
	if err != nil {
		t.Fatalf("FindServer(web1) error = %v", err)
	}
	if server.Host != "192.168.1.1" {
		t.Errorf("FindServer(web1).Host = %q, want %q", server.Host, "192.168.1.1")
	}

	// Find non-existent server
	_, err = cfg.FindServer("nonexistent")
	if err == nil {
		t.Error("FindServer(nonexistent) should return error")
	}
}

func TestAddServer(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "existing", Host: "192.168.1.1"},
		},
	}

	// Add new server
	err := cfg.AddServer(Server{Name: "new", Host: "192.168.1.2"})
	if err != nil {
		t.Fatalf("AddServer() error = %v", err)
	}
	if len(cfg.Servers) != 2 {
		t.Errorf("Servers count = %d, want 2", len(cfg.Servers))
	}

	// Try to add duplicate
	err = cfg.AddServer(Server{Name: "existing", Host: "192.168.1.3"})
	if err == nil {
		t.Error("AddServer() should return error for duplicate name")
	}
}

func TestUpdateServer(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "web1", Host: "192.168.1.1", Port: 22},
		},
	}

	// Update existing server
	err := cfg.UpdateServer("web1", Server{Name: "web1", Host: "10.0.0.1", Port: 2222})
	if err != nil {
		t.Fatalf("UpdateServer() error = %v", err)
	}
	if cfg.Servers[0].Host != "10.0.0.1" {
		t.Errorf("Updated host = %q, want %q", cfg.Servers[0].Host, "10.0.0.1")
	}
	if cfg.Servers[0].Port != 2222 {
		t.Errorf("Updated port = %d, want %d", cfg.Servers[0].Port, 2222)
	}

	// Update non-existent server
	err = cfg.UpdateServer("nonexistent", Server{Name: "nonexistent", Host: "10.0.0.2"})
	if err == nil {
		t.Error("UpdateServer() should return error for non-existent server")
	}
}

func TestRemoveServer(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "web1", Host: "192.168.1.1"},
			{Name: "web2", Host: "192.168.1.2"},
		},
	}

	// Remove existing server
	err := cfg.RemoveServer("web1")
	if err != nil {
		t.Fatalf("RemoveServer() error = %v", err)
	}
	if len(cfg.Servers) != 1 {
		t.Errorf("Servers count = %d, want 1", len(cfg.Servers))
	}
	if cfg.Servers[0].Name != "web2" {
		t.Errorf("Remaining server = %q, want %q", cfg.Servers[0].Name, "web2")
	}

	// Remove non-existent server
	err = cfg.RemoveServer("nonexistent")
	if err == nil {
		t.Error("RemoveServer() should return error for non-existent server")
	}
}

func TestFindGroup(t *testing.T) {
	cfg := &Config{
		Groups: []Group{
			{Name: "production", Color: "red"},
			{Name: "staging", Color: "yellow"},
		},
	}

	// Find existing group
	group, err := cfg.FindGroup("production")
	if err != nil {
		t.Fatalf("FindGroup(production) error = %v", err)
	}
	if group.Color != "red" {
		t.Errorf("FindGroup(production).Color = %q, want %q", group.Color, "red")
	}

	// Find non-existent group
	_, err = cfg.FindGroup("nonexistent")
	if err == nil {
		t.Error("FindGroup(nonexistent) should return error")
	}
}

func TestAddGroup(t *testing.T) {
	cfg := &Config{
		Groups: []Group{
			{Name: "production", Color: "red"},
		},
	}

	// Add new group
	err := cfg.AddGroup(Group{Name: "staging", Color: "yellow"})
	if err != nil {
		t.Fatalf("AddGroup() error = %v", err)
	}
	if len(cfg.Groups) != 2 {
		t.Errorf("Groups count = %d, want 2", len(cfg.Groups))
	}

	// Try to add duplicate
	err = cfg.AddGroup(Group{Name: "production", Color: "blue"})
	if err == nil {
		t.Error("AddGroup() should return error for duplicate name")
	}
}

func TestRemoveGroup(t *testing.T) {
	cfg := &Config{
		Groups: []Group{
			{Name: "production", Color: "red"},
			{Name: "staging", Color: "yellow"},
		},
	}

	// Remove existing group
	err := cfg.RemoveGroup("production")
	if err != nil {
		t.Fatalf("RemoveGroup() error = %v", err)
	}
	if len(cfg.Groups) != 1 {
		t.Errorf("Groups count = %d, want 1", len(cfg.Groups))
	}

	// Remove non-existent group
	err = cfg.RemoveGroup("nonexistent")
	if err == nil {
		t.Error("RemoveGroup() should return error for non-existent group")
	}
}

func TestServersByGroup(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "web1", Host: "192.168.1.1", Group: "production"},
			{Name: "web2", Host: "192.168.1.2", Group: "production"},
			{Name: "db1", Host: "192.168.1.3", Group: "staging"},
		},
	}

	prodServers := cfg.ServersByGroup("production")
	if len(prodServers) != 2 {
		t.Errorf("Production servers = %d, want 2", len(prodServers))
	}

	stagingServers := cfg.ServersByGroup("staging")
	if len(stagingServers) != 1 {
		t.Errorf("Staging servers = %d, want 1", len(stagingServers))
	}

	devServers := cfg.ServersByGroup("dev")
	if len(devServers) != 0 {
		t.Errorf("Dev servers = %d, want 0", len(devServers))
	}
}

func TestLoadWithEmptyPath(t *testing.T) {
	// When path is empty, it should use DefaultPath
	// We can't really test this fully without mocking, but we can test it doesn't panic
	cfg, err := Load("")
	if err != nil {
		// This may fail if the default path doesn't exist, which is fine
		return
	}
	if cfg == nil {
		t.Error("Load() returned nil config")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use a nested path that doesn't exist
	configPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")

	cfg := &Config{
		path:    configPath,
		Servers: []Server{{Name: "test", Host: "localhost"}},
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file was not created: %v", err)
	}
}
