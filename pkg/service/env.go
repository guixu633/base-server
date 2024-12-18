package service

import (
	"os"
	"strings"

	"github.com/guixu633/base-server/module/config"
)

func (s *Service) ParseEnv() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	return envMap
}

func (s *Service) GetConfig() *config.Config {
	return s.cfg
}
