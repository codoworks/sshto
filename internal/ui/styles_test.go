package ui

import (
	"strings"
	"testing"
)

func TestGroupTag(t *testing.T) {
	tests := []struct {
		name  string
		group string
		color string
	}{
		{"red group", "production", "red"},
		{"green group", "staging", "green"},
		{"yellow group", "dev", "yellow"},
		{"blue group", "test", "blue"},
		{"unknown color", "other", "unknown"},
		{"empty color", "nocolor", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GroupTag(tt.group, tt.color)
			if result == "" {
				t.Error("GroupTag() returned empty string")
			}
			// The tag should contain the group name
			if !strings.Contains(result, tt.group) {
				t.Errorf("GroupTag() = %q, should contain %q", result, tt.group)
			}
		})
	}
}

func TestStylesNotNil(t *testing.T) {
	// Verify all styles are initialized
	styles := []struct {
		name  string
		style interface{}
	}{
		{"TitleStyle", TitleStyle},
		{"ItemStyle", ItemStyle},
		{"SelectedItemStyle", SelectedItemStyle},
		{"DimStyle", DimStyle},
		{"ErrorStyle", ErrorStyle},
		{"SuccessStyle", SuccessStyle},
		{"WarningStyle", WarningStyle},
		{"HelpStyle", HelpStyle},
		{"GroupTagStyle", GroupTagStyle},
	}

	for _, tt := range styles {
		t.Run(tt.name, func(t *testing.T) {
			if tt.style == nil {
				t.Errorf("%s is nil", tt.name)
			}
		})
	}
}

func TestColorsNotEmpty(t *testing.T) {
	if ColorPrimary == "" {
		t.Error("ColorPrimary is empty")
	}
	if ColorSecondary == "" {
		t.Error("ColorSecondary is empty")
	}
	if ColorSuccess == "" {
		t.Error("ColorSuccess is empty")
	}
	if ColorWarning == "" {
		t.Error("ColorWarning is empty")
	}
	if ColorDanger == "" {
		t.Error("ColorDanger is empty")
	}
}

func TestGroupColorsMap(t *testing.T) {
	expectedColors := []string{"red", "green", "yellow", "blue", "magenta", "cyan", "white", "gray"}

	for _, color := range expectedColors {
		if _, ok := GroupColors[color]; !ok {
			t.Errorf("GroupColors missing %q", color)
		}
	}
}
