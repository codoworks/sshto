package config

// Group represents a server group for organization
type Group struct {
	Name  string `yaml:"name"`
	Color string `yaml:"color,omitempty"`
}
