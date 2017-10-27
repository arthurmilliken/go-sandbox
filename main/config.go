package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	DefaultPath string
	Workers     int
}

func LoadConfig() *Config {
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}
	return &conf
}
