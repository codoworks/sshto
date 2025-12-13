package config

import "testing"

func TestGroupStruct(t *testing.T) {
	g := Group{Name: "production", Color: "red"}

	if g.Name != "production" {
		t.Errorf("Name = %q, want %q", g.Name, "production")
	}
	if g.Color != "red" {
		t.Errorf("Color = %q, want %q", g.Color, "red")
	}
}

func TestGroupEmpty(t *testing.T) {
	g := Group{}

	if g.Name != "" {
		t.Errorf("Empty group Name = %q, want empty", g.Name)
	}
	if g.Color != "" {
		t.Errorf("Empty group Color = %q, want empty", g.Color)
	}
}
