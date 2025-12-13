package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codoworks/sshto/internal/config"
)

// ConnectOptions holds optional overrides for SSH connection
type ConnectOptions struct {
	User string
	Port int
	Key  string
}

// Client handles SSH command execution
type Client struct{}

// NewClient creates a new SSH client
func NewClient() *Client {
	return &Client{}
}

// Connect executes an SSH connection to the given server
func (c *Client) Connect(server *config.Server) error {
	args := c.buildArgs(server)

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// buildArgs constructs the SSH command arguments
func (c *Client) buildArgs(server *config.Server) []string {
	var args []string

	// Add identity file if specified
	if server.Key != "" {
		key := config.ExpandPath(server.Key)
		args = append(args, "-i", key)
	}

	// Add port if non-standard
	if server.Port != 0 && server.Port != 22 {
		args = append(args, "-p", strconv.Itoa(server.Port))
	}

	// Build destination
	dest := server.Host
	if server.User != "" {
		dest = server.User + "@" + server.Host
	}
	args = append(args, dest)

	return args
}

// BuildCommand returns the SSH command string for display
func (c *Client) BuildCommand(server *config.Server) string {
	args := c.buildArgs(server)
	return "ssh " + strings.Join(args, " ")
}

// TestConnection tests if an SSH connection can be established
func (c *Client) TestConnection(server *config.Server) error {
	args := c.buildArgs(server)
	args = append(args, "-o", "ConnectTimeout=5", "-o", "BatchMode=yes", "exit")

	cmd := exec.Command("ssh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connection test failed: %s", string(output))
	}
	return nil
}
