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

	Netbox struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Token string `yaml:"token"`
	} `yaml:"netbox"`

	Settings struct {
		DataCenters struct {
			Exclude []string `yaml:"exclude"`
			Include []string `yaml:"include"`
		} `yaml:"datacenters"`
	} `yaml:"settings"`
}

func GetConfig() *Config {
	config := &Config{}

	file, err := os.Open("config.yaml")
	if err != nil {
		return nil
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil
	}

	return config
}
