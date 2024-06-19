package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Cli   bool   `yaml:"cli"`
	Name  string `yaml:"name"`
	Group string `yaml:"group"`
	Exec  struct {
		Stop   string `yaml:"stop"`
		Start  string `yaml:"start"`
		Reload string `yaml:"reload"`
	}
	Deploy struct {
		Serial int `yaml:"serial"`
		Probe  *struct {
			Type string `yaml:"type"`
			Host string `yaml:"host"`
			Port string `yaml:"port"`
			Path string `yaml:"path"`
		} `yaml:"probe"`
		Roles []string `yaml:"roles"`
	} `yaml:"deploy"`
	Limits struct {
		Mem    string `yaml:"mem"`
		NoFile int    `yaml:"no_file"`
	} `yaml:"limits"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Security    struct {
		Owner string `yaml:"owner"`
		Group string `yaml:"group"`
	} `yaml:"security"`
	Resources    []string `json:"resources"`
	Capabilities []string `yaml:"capabilities"`
	Environments string   `yaml:"environments"`
}

func ParseSpec(path string) (*Spec, error) {
	var service struct {
		Spec Spec `yaml:"service"`
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &service); err != nil {
		return nil, err
	}
	return specWithDefaults(&service.Spec), nil
}

func specWithDefaults(spec *Spec) *Spec {
	if len(spec.Security.Owner) == 0 {
		spec.Security.Owner = "kinescope"
	}
	if len(spec.Security.Group) == 0 {
		spec.Security.Group = "kinescope"
	}
	if len(spec.Limits.Mem) == 0 {
		spec.Limits.Mem = "100M"
	}
	if spec.Limits.NoFile <= 0 {
		spec.Limits.NoFile = 1024
	}
	if len(spec.Exec.Stop) == 0 {
		spec.Exec.Stop = "/bin/kill -s SIGINT $MAINPID"
	}
	if len(spec.Exec.Start) == 0 {
		spec.Exec.Start = "/usr/bin/" + spec.Name
	}
	return spec
}
