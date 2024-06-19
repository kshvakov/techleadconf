package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Vars         map[string]string `yaml:"vars"`
	User         string            `yaml:"user"`
	Become       bool              `yaml:"become"`
	BecomeMethod string            `yaml:"become_method"`
	GatherFacts  bool              `yaml:"gather_facts"`
	Monitoring   struct {
		Hosts     string `yaml:"hosts"`
		FileSDDir string `yaml:"file_sd_dir"`
	} `yaml:"monitoring"`
	Resources map[string]string `yaml:"resources"`
}

func ParseConf(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
