package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		// Valid IP addresses
		{"valid ipv4", "192.168.1.1", false},
		{"valid ipv4 zeros", "0.0.0.0", false},
		{"valid ipv4 broadcast", "255.255.255.255", false},
		{"valid ipv6", "::1", false},
		{"valid ipv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"valid ipv6 compressed", "2001:db8::1", false},

		// Valid hostnames
		{"valid hostname simple", "localhost", false},
		{"valid hostname domain", "example.com", false},
		{"valid hostname subdomain", "sub.example.com", false},
		{"valid hostname with numbers", "server1.example.com", false},
		{"valid hostname with hyphens", "my-server.example.com", false},
		{"valid hostname deep subdomain", "a.b.c.d.example.com", false},

		// Invalid hosts
		{"empty host", "", true},
		{"invalid hostname starts with hyphen", "-invalid.com", true},
		{"invalid hostname ends with hyphen", "invalid-.com", true},
		{"invalid hostname with underscore", "invalid_host.com", true},
		{"invalid hostname with spaces", "invalid host.com", true},
		{"invalid hostname too long label", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHost(%q) error = %v, wantErr %v", tt.host, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHostLongHostname(t *testing.T) {
	// Generate a hostname that's too long (>253 characters)
	// Each segment is "aaaaaa." (7 chars), 40 segments = 280 chars
	longHost := ""
	for i := 0; i < 40; i++ {
		longHost += "aaaaaa."
	}
	longHost += "com"

	if len(longHost) <= 253 {
		t.Fatalf("Test hostname is only %d chars, need >253", len(longHost))
	}

	err := ValidateHost(longHost)
	if err == nil {
		t.Error("ValidateHost should reject hostnames longer than 253 characters")
	}
}

func TestValidateKeyFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sshto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a readable test file
	readableFile := filepath.Join(tmpDir, "readable_key")
	if err := os.WriteFile(readableFile, []byte("test key"), 0600); err != nil {
		t.Fatalf("Failed to create readable file: %v", err)
	}

	// Create a directory to test dir check
	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		wantWarning bool
		wantErr     bool
	}{
		{"empty path", "", false, false},
		{"readable file", readableFile, false, false},
		{"non-existent file", filepath.Join(tmpDir, "nonexistent"), true, false},
		{"directory", testDir, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning, err := ValidateKeyFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKeyFile(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
			if (warning != "") != tt.wantWarning {
				t.Errorf("ValidateKeyFile(%q) warning = %q, wantWarning %v", tt.path, warning, tt.wantWarning)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"expand tilde", "~/test", home + "/test"},
		{"expand tilde nested", "~/.ssh/id_rsa", home + "/.ssh/id_rsa"},
		{"no tilde", "/absolute/path", "/absolute/path"},
		{"relative path", "relative/path", "relative/path"},
		{"just tilde", "~", "~"},
		{"tilde without slash", "~notapath", "~notapath"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.path)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port 22", 22, false},
		{"valid port 1", 1, false},
		{"valid port 80", 80, false},
		{"valid port 443", 443, false},
		{"valid port 65535", 65535, false},
		{"valid port 0", 0, false},
		{"invalid port negative", -1, true},
		{"invalid port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-server", false},
		{"valid name with numbers", "server123", false},
		{"empty name", "", true},
		{"name too long", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateServer(t *testing.T) {
	tests := []struct {
		name    string
		server  *Server
		wantErr bool
	}{
		{
			"valid server",
			&Server{Name: "test", Host: "192.168.1.1", Port: 22},
			false,
		},
		{
			"valid server with hostname",
			&Server{Name: "test", Host: "example.com", Port: 22},
			false,
		},
		{
			"valid server zero port",
			&Server{Name: "test", Host: "example.com", Port: 0},
			false,
		},
		{
			"invalid empty name",
			&Server{Name: "", Host: "192.168.1.1", Port: 22},
			true,
		},
		{
			"invalid empty host",
			&Server{Name: "test", Host: "", Port: 22},
			true,
		},
		{
			"invalid port",
			&Server{Name: "test", Host: "192.168.1.1", Port: 70000},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServer(tt.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
