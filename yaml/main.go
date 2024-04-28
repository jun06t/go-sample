package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func main() {
	var fp string
	flag.StringVar(&fp, "c", "./config.yaml", "set yaml file path")
	flag.Parse()

	b, err := os.ReadFile(fp)
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
