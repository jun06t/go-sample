package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type EnvVars struct {
	ConfigPath string `split_words:"true" default:"./config.yaml"`
}

func main() {
	var s EnvVars
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	b, err := os.ReadFile(s.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	cfg := Config{}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", cfg)
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

type ServerParams struct {
	Endpoint string        `yaml:"endpoint"`
	Timeout  time.Duration `yaml:"timeout" default:"10s"`
	Retry    int           `yaml:"retry" default:"3"`
}

type Config struct {
	Params ServerParams `yaml:"params"`
	Key    string       `yaml:"key"`
}
