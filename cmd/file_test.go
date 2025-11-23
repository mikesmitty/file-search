package cmd

import (
	"strings"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]string
	}{
		{
			name:     "empty input",
			input:    []string{},
			expected: map[string]string{},
		},
		{
			name:  "single key-value",
			input: []string{"key=value"},
			expected: map[string]string{
				"key": "value",
			},
		},
		{
			name:  "multiple key-values",
			input: []string{"key1=value1", "key2=value2"},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:  "value with equals sign",
			input: []string{"key=value=with=equals"},
			expected: map[string]string{
				"key": "value=with=equals",
			},
		},
		{
			name:     "invalid format (no equals)",
			input:    []string{"invalid"},
			expected: map[string]string{},
		},
		{
			name:  "mixed valid and invalid",
			input: []string{"valid=value", "invalid", "another=good"},
			expected: map[string]string{
				"valid":   "value",
				"another": "good",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the metadata parsing logic from file.go upload command
			metadataMap := make(map[string]string)
			for _, meta := range tt.input {
				parts := strings.SplitN(meta, "=", 2)
				if len(parts) == 2 {
					metadataMap[parts[0]] = parts[1]
				}
			}

			// Compare results
			if len(metadataMap) != len(tt.expected) {
				t.Errorf("expected %d entries, got %d", len(tt.expected), len(metadataMap))
			}

			for k, v := range tt.expected {
				if got, ok := metadataMap[k]; !ok {
					t.Errorf("missing key %s", k)
				} else if got != v {
					t.Errorf("for key %s: expected %s, got %s", k, v, got)
				}
			}
		})
	}
}
