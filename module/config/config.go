package config

import "github.com/BurntSushi/toml"

type Config struct {
	Meta     Meta     `toml:"meta"`
	Oss      Oss      `toml:"oss"`
	LLM      LLM      `toml:"llm"`
	Workflow Workflow `toml:"workflow"`
	Telegram Telegram `toml:"telegram"`
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

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
