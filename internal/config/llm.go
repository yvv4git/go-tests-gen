package config

type LLM struct {
	OpenAI      OpenAI  `toml:"openai"`
	Temperature float64 `toml:"temperature"`
	MaxTokens   int     `toml:"max_tokens"`
}

type OpenAI struct {
	URL   string `toml:"url" env:"LLM_OPENAI_URL"`
	Token string `toml:"token" env:"LLM_OPENAI_TOKEN"`
	Model string `toml:"model"`
}
