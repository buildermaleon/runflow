package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear all env vars
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("REDIS_URL")
	
	cfg := Load()
	
	if cfg.Port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", cfg.Port)
	}
	
	if cfg.DatabaseURL == "" {
		t.Error("Expected non-empty DatabaseURL")
	}
	
	if cfg.RedisURL == "" {
		t.Error("Expected non-empty RedisURL")
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("REDIS_URL", "redis://test:6379")
	
	cfg := Load()
	
	if cfg.Port != "3000" {
		t.Errorf("Expected port '3000', got '%s'", cfg.Port)
	}
	
	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/test" {
		t.Errorf("Expected custom DatabaseURL, got '%s'", cfg.DatabaseURL)
	}
	
	if cfg.RedisURL != "redis://test:6379" {
		t.Errorf("Expected custom RedisURL, got '%s'", cfg.RedisURL)
	}
	
	// Cleanup
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("REDIS_URL")
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		def      string
		expected string
	}{
		{"TEST_KEY", "test_value", "default", "test_value"},
		{"TEST_EMPTY", "", "default", "default"},
		{"TEST_UNSET", "", "default", "default"},
	}
	
	for _, tt := range tests {
		if tt.value != "" {
			os.Setenv(tt.key, tt.value)
		} else {
			os.Unsetenv(tt.key)
		}
		
		result := getEnv(tt.key, tt.def)
		if result != tt.expected {
			t.Errorf("Expected '%s', got '%s'", tt.expected, result)
		}
		
		os.Unsetenv(tt.key)
	}
}
