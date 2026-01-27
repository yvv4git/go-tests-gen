package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// Set environment variables that should override config file values
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("LOG_LEVEL")
	}()

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "info"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config based on test file and env override
	expected := Config{
		LogLevel: "debug",
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}

func TestConfigFromFile(t *testing.T) {
	// Ensure LOG_LEVEL is not set
	os.Unsetenv("LOG_LEVEL")

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "warn"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config from file
	expected := Config{
		LogLevel: "warn",
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}
