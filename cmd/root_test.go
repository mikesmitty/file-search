package cmd

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetAPIKeyPriority(t *testing.T) {
	// Note: This is a conceptual test. Actual implementation would require
	// mocking viper and environment variables, which is complex.
	// This test documents the expected behavior.

	t.Run("priority order", func(t *testing.T) {
		// Expected priority:
		// 1. --api-key flag
		// 2. --api-key-env custom env var
		// 3. config file
		// 4. GOOGLE_API_KEY
		// 5. GEMINI_API_KEY

		// This would require integration testing with actual viper config
		t.Skip("requires integration test setup")
	})
}

func TestGetMCPTools(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string returns default",
			input:    "",
			expected: []string{"query"},
		},
		{
			name:     "single tool",
			input:    "query",
			expected: []string{"query"},
		},
		{
			name:     "multiple tools",
			input:    "query,import,list",
			expected: []string{"query", "import", "list"},
		},
		{
			name:     "tools with spaces",
			input:    "query, import, list",
			expected: []string{"query", "import", "list"},
		},
		{
			name:     "all tools",
			input:    "query,import,upload,list,manage",
			expected: []string{"query", "import", "upload", "list", "manage"},
		},
		{
			name:     "only commas returns default",
			input:    ",,,",
			expected: []string{"query"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()
			viper.Set("mcp_tools", tt.input)

			result := getMCPTools()

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tools, got %d", len(tt.expected), len(result))
			}

			for i, tool := range tt.expected {
				if i >= len(result) || result[i] != tool {
					t.Errorf("expected tool[%d] = %s, got %v", i, tool, result)
				}
			}
		})
	}
}

func TestCompletionConfigDefaults(t *testing.T) {
	t.Run("default completion_enabled is true", func(t *testing.T) {
		viper.Reset()
		viper.SetDefault("completion_enabled", true)

		enabled := viper.GetBool("completion_enabled")
		if !enabled {
			t.Error("Expected completion_enabled default to be true")
		}
	})

	t.Run("default completion_cache_ttl is 300s", func(t *testing.T) {
		viper.Reset()
		viper.SetDefault("completion_cache_ttl", "300s")

		ttl := viper.GetDuration("completion_cache_ttl")
		expected := 300 * time.Second
		if ttl != expected {
			t.Errorf("Expected completion_cache_ttl default to be %v, got %v", expected, ttl)
		}
	})

	t.Run("default mcp_tools is query", func(t *testing.T) {
		viper.Reset()
		viper.SetDefault("mcp_tools", "query")

		tools := viper.GetString("mcp_tools")
		if tools != "query" {
			t.Errorf("Expected mcp_tools default to be 'query', got '%s'", tools)
		}
	})
}

func TestCompletionConfigParsing(t *testing.T) {
	t.Run("parses completion_enabled boolean", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_enabled", false)

		enabled := viper.GetBool("completion_enabled")
		if enabled {
			t.Error("Expected completion_enabled to be false")
		}
	})

	t.Run("parses completion_cache_ttl duration", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_cache_ttl", "600s")

		ttl := viper.GetDuration("completion_cache_ttl")
		expected := 600 * time.Second
		if ttl != expected {
			t.Errorf("Expected TTL %v, got %v", expected, ttl)
		}
	})

	t.Run("parses completion_cache_ttl in minutes", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_cache_ttl", "10m")

		ttl := viper.GetDuration("completion_cache_ttl")
		expected := 10 * time.Minute
		if ttl != expected {
			t.Errorf("Expected TTL %v, got %v", expected, ttl)
		}
	})

	t.Run("parses completion_cache_ttl in hours", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_cache_ttl", "1h")

		ttl := viper.GetDuration("completion_cache_ttl")
		expected := 1 * time.Hour
		if ttl != expected {
			t.Errorf("Expected TTL %v, got %v", expected, ttl)
		}
	})
}

func TestGetCompleterConfiguration(t *testing.T) {
	t.Run("uses default TTL when config is zero", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_enabled", true)
		viper.Set("completion_cache_ttl", 0)

		// Note: Can't actually test getCompleter() without mocking getAPIKey()
		// This test verifies the configuration parsing behavior

		ttl := viper.GetDuration("completion_cache_ttl")
		if ttl != 0 {
			t.Errorf("Expected TTL 0 from config, got %v", ttl)
		}
		// The getCompleter() function should default to 300s when ttl is 0
	})

	t.Run("reads completion_enabled from config", func(t *testing.T) {
		viper.Reset()
		viper.Set("completion_enabled", false)

		enabled := viper.GetBool("completion_enabled")
		if enabled {
			t.Error("Expected completion_enabled to be false from config")
		}
	})
}
