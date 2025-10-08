package greet_test

import (
	"errors"
	"testing"

	"github.com/brunoluiz/go-lab/services/hello-world/internal/service/greet"
)

func TestNew(t *testing.T) {
	greeter := greet.New()
	if greeter == nil {
		t.Fatal("Expected non-nil greeter")
	}
}

func TestGreeter_Hello(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected string
		wantErr  bool
		errType  error
	}{
		{
			name:     "English greeting",
			lang:     "en",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "Portuguese greeting",
			lang:     "pt",
			expected: "olá",
			wantErr:  false,
		},
		{
			name:     "English with region doesn't match",
			lang:     "en-US",
			expected: "",
			wantErr:  true,
			errType:  greet.ErrNotImplemented,
		},
		{
			name:     "Portuguese with region doesn't match",
			lang:     "pt-BR",
			expected: "",
			wantErr:  true,
			errType:  greet.ErrNotImplemented,
		},
		{
			name:     "Unsupported language returns error",
			lang:     "de",
			expected: "",
			wantErr:  true,
			errType:  greet.ErrNotImplemented,
		},
		{
			name:     "Spanish language returns error",
			lang:     "es",
			expected: "",
			wantErr:  true,
			errType:  greet.ErrNotImplemented,
		},
		{
			name:     "Invalid language string returns parse error",
			lang:     "invalid-lang-code",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty language string returns parse error",
			lang:     "",
			expected: "",
			wantErr:  true,
		},
	}

	greeter := greet.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := greeter.Hello(tt.lang)

			if result != tt.expected {
				t.Errorf("Expected result %q, got %q", tt.expected, result)
			}

			if tt.wantErr && err == nil {
				t.Error("Expected an error, but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			// Check specific error type if specified
			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Expected error %v, got %v", tt.errType, err)
			}
		})
	}
}
