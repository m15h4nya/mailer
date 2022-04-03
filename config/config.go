package config

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	ListenAddr string
	DB         string
}

func Load(path string) (*Config, error) {

	file, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{
		ListenAddr: file.Get("service.address").(string),
		DB:         file.Get("service.sqlite.db").(string),
	}

	return config, nil
}
