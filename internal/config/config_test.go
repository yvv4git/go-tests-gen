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
[llm]
  temperature = 0.7
  max_tokens = 1000
  [llm.openai]
    url = "https://api.openai.com/v1"
    token = "test-token"
    model = "gpt-4"
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
		LLM: LLM{
			Temperature: 0.7,
			MaxTokens:   1000,
			OpenAI: OpenAI{
				URL:   "https://api.openai.com/v1",
				Token: "test-token",
				Model: "gpt-4",
			},
		},
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
[llm]
  temperature = 0.5
  max_tokens = 500
  [llm.openai]
    url = "https://api.example.com/v1"
    token = "file-token"
    model = "gpt-3.5-turbo"
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
		LLM: LLM{
			Temperature: 0.5,
			MaxTokens:   500,
			OpenAI: OpenAI{
				URL:   "https://api.example.com/v1",
				Token: "file-token",
				Model: "gpt-3.5-turbo",
			},
		},
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}

func TestLLMTokenFromEnv(t *testing.T) {
	// Set LLM token from environment
	os.Setenv("LLM_OPENAI_TOKEN", "env-test-token")
	defer func() {
		os.Unsetenv("LLM_OPENAI_TOKEN")
	}()

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "info"
[llm]
  temperature = 0.7
  max_tokens = 1000
  [llm.openai]
    url = "https://api.openai.com/v1"
    token = "file-token"
    model = "gpt-4"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config with token from env
	expected := Config{
		LogLevel: "info",
		LLM: LLM{
			Temperature: 0.7,
			MaxTokens:   1000,
			OpenAI: OpenAI{
				URL:   "https://api.openai.com/v1",
				Token: "env-test-token",
				Model: "gpt-4",
			},
		},
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}

func TestLLMModelFromEnv(t *testing.T) {
	// Set LLM model from environment
	os.Setenv("LLM_OPENAI_MODEL", "env-gpt-4o")
	defer func() {
		os.Unsetenv("LLM_OPENAI_MODEL")
	}()

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "info"
[llm]
  temperature = 0.7
  max_tokens = 1000
  [llm.openai]
    url = "https://api.openai.com/v1"
    token = "file-token"
    model = "gpt-4"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config with model from env
	expected := Config{
		LogLevel: "info",
		LLM: LLM{
			Temperature: 0.7,
			MaxTokens:   1000,
			OpenAI: OpenAI{
				URL:   "https://api.openai.com/v1",
				Token: "file-token",
				Model: "env-gpt-4o",
			},
		},
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}

func TestLLMURLFromEnv(t *testing.T) {
	// Set LLM URL from environment
	os.Setenv("LLM_OPENAI_URL", "https://env-url.example.com/v1")
	defer func() {
		os.Unsetenv("LLM_OPENAI_URL")
	}()

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "info"
[llm]
  temperature = 0.7
  max_tokens = 1000
  [llm.openai]
    url = "https://api.openai.com/v1"
    token = "file-token"
    model = "gpt-4"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config with URL from env
	expected := Config{
		LogLevel: "info",
		LLM: LLM{
			Temperature: 0.7,
			MaxTokens:   1000,
			OpenAI: OpenAI{
				URL:   "https://env-url.example.com/v1",
				Token: "file-token",
				Model: "gpt-4",
			},
		},
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}
