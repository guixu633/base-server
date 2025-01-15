package config

import "github.com/BurntSushi/toml"

type Config struct {
	Meta          Meta          `toml:"meta"`
	Oss           Oss           `toml:"oss"`
	LLM           LLM           `toml:"llm"`
	Workflow      Workflow      `toml:"workflow"`
	Telegram      Telegram      `toml:"telegram"`
	Qdrant        Qdrant        `toml:"qdrant"`
	Embedding     Embedding     `toml:"embedding"`
	CryptoArticle CryptoArticle `toml:"crypto_article"`
}

type Meta struct {
	Env string `toml:"env"`
}

type Oss struct {
	Endpoint  string `toml:"endpoint"`
	Bucket    string `toml:"bucket"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

type LLM struct {
	OpenaiApiKey   string `toml:"openai_api_key"`
	DeepseekApiKey string `toml:"deepseek_api_key"`
}

type Workflow struct {
	Url string `toml:"url"`
}

type Telegram struct {
	ApiToken string `toml:"api_token"`
}

type Qdrant struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	ApiKey    string `toml:"api_key"`
	Dimension int    `toml:"dimension"`
}

type Embedding struct {
	Url       string `toml:"url"`
	ApiKey    string `toml:"api_key"`
	Model     string `toml:"model"`
	Dimension int    `toml:"dimension"`
}

type CryptoArticle struct {
	Collection     string  `toml:"collection"`
	ScoreThreshold float32 `toml:"score_threshold"`
	Limit          uint64  `toml:"limit"`
	HoursLimit     uint64  `toml:"hours_limit"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
