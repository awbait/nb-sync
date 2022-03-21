package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Providers struct {
		VSphere struct {
			Enable   bool   `yaml:"enable"`
			Host     string `yaml:"host"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"vsphere"`
	} `yaml:"providers"`
}

func GetConfig() (*Config, error) {
	config := &Config{}

	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
