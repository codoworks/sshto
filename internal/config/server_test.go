package config

import "testing"

func TestServerFilterValue(t *testing.T) {
	s := Server{Name: "web1", Host: "192.168.1.1", Group: "production"}
	expected := "web1"
	if s.FilterValue() != expected {
		t.Errorf("FilterValue() = %q, want %q", s.FilterValue(), expected)
	}
}

func TestServerTitle(t *testing.T) {
	s := Server{Name: "web1", Host: "192.168.1.1"}
	if s.Title() != "web1" {
		t.Errorf("Title() = %q, want %q", s.Title(), "web1")
	}
}

func TestServerDescription(t *testing.T) {
	tests := []struct {
		name     string
		server   Server
		expected string
	}{
		{
			"host only",
			Server{Host: "192.168.1.1"},
			"192.168.1.1",
		},
		{
			"user and host",
			Server{Host: "192.168.1.1", User: "admin"},
			"admin@192.168.1.1",
		},
		{
			"user, host and non-standard port",
			Server{Host: "192.168.1.1", User: "admin", Port: 2222},
			"admin@192.168.1.1:2222",
		},
		{
			"host and standard port",
			Server{Host: "192.168.1.1", Port: 22},
			"192.168.1.1",
		},
		{
			"host and port zero",
			Server{Host: "192.168.1.1", Port: 0},
			"192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.server.Description()
			if result != tt.expected {
				t.Errorf("Description() = %q, want %q", result, tt.expected)
			}
		})
	}
}
