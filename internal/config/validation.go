package config

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// hostnameRegex validates RFC 1123 hostnames
var hostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

// ValidateHost validates that the host is a valid IP address or hostname
func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}

	// Check if it's a valid IP address (v4 or v6)
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	// Check if it's a valid hostname (RFC 1123)
	if len(host) > 253 {
		return fmt.Errorf("hostname too long (max 253 characters)")
	}

	if !hostnameRegex.MatchString(host) {
		return fmt.Errorf("invalid hostname format")
	}

	// Check each label length (max 63 characters)
	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) > 63 {
			return fmt.Errorf("hostname label too long (max 63 characters)")
		}
	}

	return nil
}

// ValidateKeyFile checks if the key file exists
// Returns a warning message if file doesn't exist, nil otherwise
func ValidateKeyFile(path string) (warning string, err error) {
	if path == "" {
		return "", nil
	}

	expanded := ExpandPath(path)

	info, err := os.Stat(expanded)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Sprintf("key file does not exist: %s", path), nil
		}
		return "", fmt.Errorf("cannot access key file: %w", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("key file is a directory: %s", path)
	}

	// Check if file is readable (basic permission check)
	mode := info.Mode()
	if mode.Perm()&0400 == 0 {
		return fmt.Sprintf("key file may not be readable: %s", path), nil
	}

	return "", nil
}

// ExpandPath expands ~ to home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}

// ValidatePort validates that the port is in valid range
func ValidatePort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}
	return nil
}

// ValidateName validates the server name
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) > 64 {
		return fmt.Errorf("name too long (max 64 characters)")
	}
	return nil
}

// ValidateServer validates all fields of a server
func ValidateServer(s *Server) error {
	if err := ValidateName(s.Name); err != nil {
		return err
	}
	if err := ValidateHost(s.Host); err != nil {
		return err
	}
	if err := ValidatePort(s.Port); err != nil {
		return err
	}
	return nil
}
