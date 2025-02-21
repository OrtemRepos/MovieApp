package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type serviceConfig struct {
	APIConfig apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

func loadConfig() *serviceConfig {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
