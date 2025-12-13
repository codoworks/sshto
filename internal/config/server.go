package config

import "fmt"

// Server represents an SSH server configuration
type Server struct {
	Name  string `yaml:"name"`
	Host  string `yaml:"host"`
	User  string `yaml:"user,omitempty"`
	Port  int    `yaml:"port,omitempty"`
	Key   string `yaml:"key,omitempty"`
	Group string `yaml:"group,omitempty"`
}

// FilterValue implements list.Item for bubbles list
func (s Server) FilterValue() string {
	return s.Name
}

// Title implements list.DefaultItem
func (s Server) Title() string {
	return s.Name
}

// Description implements list.DefaultItem
func (s Server) Description() string {
	desc := s.Host
	if s.User != "" {
		desc = s.User + "@" + desc
	}
	if s.Port != 0 && s.Port != 22 {
		desc = fmt.Sprintf("%s:%d", desc, s.Port)
	}
	return desc
}
