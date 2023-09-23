package main

import "time"

type Config struct {
	Static  []StaticConfig `yaml:"static"`
	Dynamic DynamicConfig  `yaml:"dynamic"`
}

type StaticConfig struct {
	Group    string         `yaml:"group"`
	Interval time.Duration  `yaml:"interval"`
	Servers  []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

type DynamicConfig struct {
	Tokens []string `yaml:"tokens"`
}
