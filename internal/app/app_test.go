package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codoworks/sshto/internal/config"
	"github.com/codoworks/sshto/internal/ssh"
)

func TestNew(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Test with non-existent config (should create empty config)
	app, err := New(configPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if app == nil {
		t.Fatal("New() returned nil app")
	}
	if app.Config == nil {
		t.Error("App.Config is nil")
	}
	if app.SSHClient == nil {
		t.Error("App.SSHClient is nil")
	}
}

func TestNewWithExistingConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create a config file
	configContent := `servers:
  - name: test-server
    host: 192.168.1.1
    user: admin
    port: 22
defaults:
  user: deploy
  port: 2222
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	app, err := New(configPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if len(app.Config.Servers) != 1 {
		t.Errorf("Servers count = %d, want 1", len(app.Config.Servers))
	}
	if app.Config.Servers[0].Name != "test-server" {
		t.Errorf("Server name = %q, want %q", app.Config.Servers[0].Name, "test-server")
	}
}

func TestResolveServer(t *testing.T) {
	app := &App{
		Config: &config.Config{
			Defaults: config.Defaults{
				User: "default-user",
				Port: 2222,
				Key:  "~/.ssh/default_key",
			},
		},
		SSHClient: ssh.NewClient(),
	}

	tests := []struct {
		name         string
		server       *config.Server
		expectedUser string
		expectedPort int
		expectedKey  string
	}{
		{
			"all defaults applied",
			&config.Server{Name: "test", Host: "localhost"},
			"default-user",
			2222,
			"~/.ssh/default_key",
		},
		{
			"server values override defaults",
			&config.Server{Name: "test", Host: "localhost", User: "custom", Port: 3333, Key: "~/.ssh/custom"},
			"custom",
			3333,
			"~/.ssh/custom",
		},
		{
			"partial override",
			&config.Server{Name: "test", Host: "localhost", User: "custom"},
			"custom",
			2222,
			"~/.ssh/default_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved := app.resolveServer(tt.server)
			if resolved.User != tt.expectedUser {
				t.Errorf("User = %q, want %q", resolved.User, tt.expectedUser)
			}
			if resolved.Port != tt.expectedPort {
				t.Errorf("Port = %d, want %d", resolved.Port, tt.expectedPort)
			}
			if resolved.Key != tt.expectedKey {
				t.Errorf("Key = %q, want %q", resolved.Key, tt.expectedKey)
			}
		})
	}
}

func TestResolveServerNoDefaults(t *testing.T) {
	app := &App{
		Config: &config.Config{
			Defaults: config.Defaults{},
		},
		SSHClient: ssh.NewClient(),
	}

	server := &config.Server{Name: "test", Host: "localhost"}
	resolved := app.resolveServer(server)

	// Should default to port 22 when no defaults.Port
	if resolved.Port != 22 {
		t.Errorf("Port = %d, want 22 (default)", resolved.Port)
	}
}

func TestSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	app, err := New(configPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add a server
	app.Config.AddServer(config.Server{Name: "test", Host: "localhost"})

	// Save
	if err := app.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file not created: %v", err)
	}
}

func TestConnectServerNotFound(t *testing.T) {
	app := &App{
		Config: &config.Config{
			Servers: []config.Server{},
		},
		SSHClient: ssh.NewClient(),
	}

	err := app.Connect("nonexistent", ssh.ConnectOptions{})
	if err == nil {
		t.Error("Connect() should return error for non-existent server")
	}
}

// TestConnectWithOverrides tests that overrides are applied correctly
// Note: We can't test the actual SSH connection without mocking
func TestConnectResolution(t *testing.T) {
	// Create a mock-like test by checking the resolution logic
	app := &App{
		Config: &config.Config{
			Servers: []config.Server{
				{Name: "test", Host: "localhost", User: "original", Port: 22},
			},
			Defaults: config.Defaults{},
		},
		SSHClient: ssh.NewClient(),
	}

	// Find the server
	server, err := app.Config.FindServer("test")
	if err != nil {
		t.Fatalf("FindServer() error = %v", err)
	}

	// Test resolution
	resolved := app.resolveServer(server)
	if resolved.User != "original" {
		t.Errorf("Resolved user = %q, want %q", resolved.User, "original")
	}

	// Test with overrides (simulating what Connect does)
	opts := ssh.ConnectOptions{User: "override", Port: 2222, Key: "/tmp/key"}

	if opts.User != "" {
		resolved.User = opts.User
	}
	if opts.Port != 0 {
		resolved.Port = opts.Port
	}
	if opts.Key != "" {
		resolved.Key = opts.Key
	}

	if resolved.User != "override" {
		t.Errorf("Override user = %q, want %q", resolved.User, "override")
	}
	if resolved.Port != 2222 {
		t.Errorf("Override port = %d, want %d", resolved.Port, 2222)
	}
	if resolved.Key != "/tmp/key" {
		t.Errorf("Override key = %q, want %q", resolved.Key, "/tmp/key")
	}
}
