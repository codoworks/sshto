package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Defaults holds default values for server connections
type Defaults struct {
	User string `yaml:"user,omitempty"`
	Port int    `yaml:"port,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

// Config represents the full configuration file
type Config struct {
	Groups   []Group  `yaml:"groups,omitempty"`
	Servers  []Server `yaml:"servers"`
	Defaults Defaults `yaml:"defaults,omitempty"`

	path string // internal: path to config file
}

// DefaultPath returns the default config file path
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sshto", "config.yaml")
}

// Load reads the config from the given path
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Return empty config if file doesn't exist
			return &Config{
				path:     path,
				Defaults: Defaults{Port: 22},
			}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	cfg.path = path
	return &cfg, nil
}

// Save writes the config to disk
func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// Path returns the config file path
func (c *Config) Path() string {
	return c.path
}

// FindServer returns a server by name
func (c *Config) FindServer(name string) (*Server, error) {
	for i := range c.Servers {
		if c.Servers[i].Name == name {
			return &c.Servers[i], nil
		}
	}
	return nil, fmt.Errorf("server %q not found", name)
}

// AddServer adds a new server to the config
func (c *Config) AddServer(s Server) error {
	for _, existing := range c.Servers {
		if existing.Name == s.Name {
			return fmt.Errorf("server %q already exists", s.Name)
		}
	}
	c.Servers = append(c.Servers, s)
	return nil
}

// UpdateServer updates an existing server
func (c *Config) UpdateServer(name string, s Server) error {
	for i := range c.Servers {
		if c.Servers[i].Name == name {
			c.Servers[i] = s
			return nil
		}
	}
	return fmt.Errorf("server %q not found", name)
}

// RemoveServer removes a server by name
func (c *Config) RemoveServer(name string) error {
	for i := range c.Servers {
		if c.Servers[i].Name == name {
			c.Servers = append(c.Servers[:i], c.Servers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("server %q not found", name)
}

// FindGroup returns a group by name
func (c *Config) FindGroup(name string) (*Group, error) {
	for i := range c.Groups {
		if c.Groups[i].Name == name {
			return &c.Groups[i], nil
		}
	}
	return nil, fmt.Errorf("group %q not found", name)
}

// AddGroup adds a new group to the config
func (c *Config) AddGroup(g Group) error {
	for _, existing := range c.Groups {
		if existing.Name == g.Name {
			return fmt.Errorf("group %q already exists", g.Name)
		}
	}
	c.Groups = append(c.Groups, g)
	return nil
}

// RemoveGroup removes a group by name
func (c *Config) RemoveGroup(name string) error {
	for i := range c.Groups {
		if c.Groups[i].Name == name {
			c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("group %q not found", name)
}

// ServersByGroup returns servers belonging to a specific group
func (c *Config) ServersByGroup(group string) []Server {
	var servers []Server
	for _, s := range c.Servers {
		if s.Group == group {
			servers = append(servers, s)
		}
	}
	return servers
}
