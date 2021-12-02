package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
)

func LoadConfig(cfg interface{}, fileNames ...string) {
	if len(fileNames) == 0 {
		fileNames = []string{".env", ".env.local"}
	}

	var valid []string
	for _, f := range fileNames {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		valid = append(valid, f)
	}
	_ = godotenv.Overload(valid...)
	_ = env.Parse(cfg)
	return
}
