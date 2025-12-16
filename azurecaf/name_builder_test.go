package azurecaf

import (
	"testing"
)

func TestNewNameBuilder(t *testing.T) {
	builder := NewNameBuilder(63, "-")

	if builder.MaxLength != 63 {
		t.Errorf("expected MaxLength 63, got %d", builder.MaxLength)
	}
	if builder.Separator != "-" {
		t.Errorf("expected Separator '-', got %s", builder.Separator)
	}
	if len(builder.content) != 0 {
		t.Errorf("expected empty content, got %d elements", len(builder.content))
	}
}

func TestNameBuilder_Append(t *testing.T) {
	tests := []struct {
		name         string
		maxLength    int
		separator    string
		segments     []string
		wantName     string
		wantFullName string
	}{
		{
			name:         "single segment",
			maxLength:    63,
			separator:    "-",
			segments:     []string{"myapp"},
			wantName:     "myapp",
			wantFullName: "myapp",
		},
		{
			name:         "multiple segments",
			maxLength:    63,
			separator:    "-",
			segments:     []string{"rg", "myapp", "prod"},
			wantName:     "rg-myapp-prod",
			wantFullName: "rg-myapp-prod",
		},
		{
			name:         "empty separator",
			maxLength:    63,
			separator:    "",
			segments:     []string{"st", "myapp", "prod"},
			wantName:     "stmyappprod",
			wantFullName: "stmyappprod",
		},
		{
			name:         "all segments fit",
			maxLength:    20,
			separator:    "-",
			segments:     []string{"rg", "app", "dev"},
			wantName:     "rg-app-dev",
			wantFullName: "rg-app-dev",
		},
		{
			name:         "last segment excluded",
			maxLength:    10,
			separator:    "-",
			segments:     []string{"rg", "app", "development"},
			wantName:     "rg-app",
			wantFullName: "rg-app-development",
		},
		{
			name:         "extra separator fits",
			maxLength:    10,
			separator:    "-",
			segments:     []string{"rg", "app", "de"},
			wantName:     "rg-app-de",
			wantFullName: "rg-app-de",
		},
		{
			name:         "empty builder",
			maxLength:    63,
			separator:    "-",
			segments:     []string{},
			wantName:     "",
			wantFullName: "",
		},
		{
			name:         "empty segment",
			maxLength:    63,
			separator:    "-",
			segments:     []string{"", "myapp"},
			wantName:     "-myapp",
			wantFullName: "-myapp",
		},
		{
			name:         "zero max length",
			maxLength:    0,
			separator:    "-",
			segments:     []string{"rg"},
			wantName:     "",
			wantFullName: "rg",
		},
		{
			name:         "segment exceeds max length",
			maxLength:    5,
			separator:    "-",
			segments:     []string{"abcdef"},
			wantName:     "",
			wantFullName: "abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewNameBuilder(tt.maxLength, tt.separator)
			for _, seg := range tt.segments {
				builder.Append(seg)
			}

			if got := builder.GetTrimmedName(); got != tt.wantName {
				t.Errorf("GetTrimmedName() = %q, want %q", got, tt.wantName)
			}
			if got := builder.GetName(); got != tt.wantFullName {
				t.Errorf("GetName() = %q, want %q", got, tt.wantFullName)
			}
		})
	}
}

func TestNameBuilder_Prepend(t *testing.T) {
	tests := []struct {
		name         string
		maxLength    int
		separator    string
		prepends     []string
		wantName     string
		wantFullName string
	}{
		{
			name:         "single prepend",
			maxLength:    63,
			separator:    "-",
			prepends:     []string{"rg"},
			wantName:     "rg",
			wantFullName: "rg",
		},
		{
			name:         "multiple prepends",
			maxLength:    63,
			separator:    "-",
			prepends:     []string{"prod", "myapp", "rg"},
			wantName:     "rg-myapp-prod",
			wantFullName: "rg-myapp-prod",
		},
		{
			name:         "prepends with empty separator",
			maxLength:    63,
			separator:    "",
			prepends:     []string{"prod", "myapp", "st"},
			wantName:     "stmyappprod",
			wantFullName: "stmyappprod",
		},
		{
			name:         "all segments fit",
			maxLength:    20,
			separator:    "-",
			prepends:     []string{"dev", "app", "rg"},
			wantName:     "rg-app-dev",
			wantFullName: "rg-app-dev",
		},
		{
			name:         "last segment excluded",
			maxLength:    10,
			separator:    "-",
			prepends:     []string{"development", "app", "rg"},
			wantName:     "rg-app",
			wantFullName: "rg-app-development",
		},
		{
			name:         "extra separator fits",
			maxLength:    10,
			separator:    "-",
			prepends:     []string{"de", "app", "rg"},
			wantName:     "rg-app-de",
			wantFullName: "rg-app-de",
		},
		{
			name:         "empty builder",
			maxLength:    63,
			separator:    "-",
			prepends:     []string{},
			wantName:     "",
			wantFullName: "",
		},
		{
			name:         "empty segment",
			maxLength:    63,
			separator:    "-",
			prepends:     []string{"myapp", ""},
			wantName:     "-myapp",
			wantFullName: "-myapp",
		},
		{
			name:         "zero max length",
			maxLength:    0,
			separator:    "-",
			prepends:     []string{"rg"},
			wantName:     "",
			wantFullName: "rg",
		},
		{
			name:         "segment exceeds max length",
			maxLength:    5,
			separator:    "-",
			prepends:     []string{"abcdef"},
			wantName:     "",
			wantFullName: "abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewNameBuilder(tt.maxLength, tt.separator)
			for _, seg := range tt.prepends {
				builder.Prepend(seg)
			}

			if got := builder.GetTrimmedName(); got != tt.wantName {
				t.Errorf("GetTrimmedName() = %q, want %q", got, tt.wantName)
			}
			if got := builder.GetName(); got != tt.wantFullName {
				t.Errorf("GetName() = %q, want %q", got, tt.wantFullName)
			}
		})
	}
}
