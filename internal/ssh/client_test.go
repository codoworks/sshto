package ssh

import (
	"os"
	"strings"
	"testing"

	"github.com/codoworks/sshto/internal/config"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient() returned nil")
	}
}

func TestBuildArgs(t *testing.T) {
	client := NewClient()
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		server   *config.Server
		expected []string
	}{
		{
			"basic host only",
			&config.Server{Host: "192.168.1.1"},
			[]string{"192.168.1.1"},
		},
		{
			"host with user",
			&config.Server{Host: "192.168.1.1", User: "admin"},
			[]string{"admin@192.168.1.1"},
		},
		{
			"host with non-standard port",
			&config.Server{Host: "192.168.1.1", Port: 2222},
			[]string{"-p", "2222", "192.168.1.1"},
		},
		{
			"host with standard port",
			&config.Server{Host: "192.168.1.1", Port: 22},
			[]string{"192.168.1.1"},
		},
		{
			"host with key",
			&config.Server{Host: "192.168.1.1", Key: "~/.ssh/id_rsa"},
			[]string{"-i", home + "/.ssh/id_rsa", "192.168.1.1"},
		},
		{
			"full config",
			&config.Server{Host: "192.168.1.1", User: "admin", Port: 2222, Key: "~/.ssh/mykey"},
			[]string{"-i", home + "/.ssh/mykey", "-p", "2222", "admin@192.168.1.1"},
		},
		{
			"absolute key path",
			&config.Server{Host: "192.168.1.1", Key: "/etc/ssh/key"},
			[]string{"-i", "/etc/ssh/key", "192.168.1.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := client.buildArgs(tt.server)
			if len(args) != len(tt.expected) {
				t.Errorf("buildArgs() = %v, want %v", args, tt.expected)
				return
			}
			for i := range args {
				if args[i] != tt.expected[i] {
					t.Errorf("buildArgs()[%d] = %q, want %q", i, args[i], tt.expected[i])
				}
			}
		})
	}
}

func TestBuildCommand(t *testing.T) {
	client := NewClient()
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		server   *config.Server
		expected string
	}{
		{
			"basic command",
			&config.Server{Host: "192.168.1.1"},
			"ssh 192.168.1.1",
		},
		{
			"command with user",
			&config.Server{Host: "192.168.1.1", User: "admin"},
			"ssh admin@192.168.1.1",
		},
		{
			"command with port",
			&config.Server{Host: "192.168.1.1", Port: 2222},
			"ssh -p 2222 192.168.1.1",
		},
		{
			"full command",
			&config.Server{Host: "192.168.1.1", User: "admin", Port: 2222, Key: "~/.ssh/mykey"},
			"ssh -i " + home + "/.ssh/mykey -p 2222 admin@192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := client.BuildCommand(tt.server)
			if cmd != tt.expected {
				t.Errorf("BuildCommand() = %q, want %q", cmd, tt.expected)
			}
		})
	}
}

func TestConnectOptionsStruct(t *testing.T) {
	opts := ConnectOptions{
		User: "testuser",
		Port: 2222,
		Key:  "/path/to/key",
	}

	if opts.User != "testuser" {
		t.Errorf("User = %q, want %q", opts.User, "testuser")
	}
	if opts.Port != 2222 {
		t.Errorf("Port = %d, want %d", opts.Port, 2222)
	}
	if opts.Key != "/path/to/key" {
		t.Errorf("Key = %q, want %q", opts.Key, "/path/to/key")
	}
}

func TestBuildArgsPreservesOrder(t *testing.T) {
	client := NewClient()
	server := &config.Server{
		Host: "example.com",
		User: "user",
		Port: 2222,
		Key:  "/tmp/key",
	}

	args := client.buildArgs(server)
	argsStr := strings.Join(args, " ")

	// Key should come before port, port before destination
	keyIdx := strings.Index(argsStr, "-i")
	portIdx := strings.Index(argsStr, "-p")
	destIdx := strings.Index(argsStr, "user@example.com")

	if keyIdx == -1 || portIdx == -1 || destIdx == -1 {
		t.Fatalf("Missing expected args: got %q", argsStr)
	}

	if keyIdx > portIdx {
		t.Error("Key flag should come before port flag")
	}
	if portIdx > destIdx {
		t.Error("Port flag should come before destination")
	}
}

func TestBuildArgsPortZero(t *testing.T) {
	client := NewClient()
	server := &config.Server{
		Host: "example.com",
		Port: 0,
	}

	args := client.buildArgs(server)
	for _, arg := range args {
		if arg == "-p" || arg == "0" {
			t.Error("Port 0 should not add -p flag")
		}
	}
}
