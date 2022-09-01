package main

import (
    "os"
    "fmt"
    "flag"
    "gopkg.in/yaml.v2"
)

type ServerConfig struct {
    Port int `yaml:"port"`
    Backend string `yaml:"backend"`
}

type BackendConfig struct {
    Name string `yaml:"name"`
    Driver string `yaml:"driver"`
    Output string `yaml:"output"`
}

type Config struct {
    Listen []*ServerConfig
    Backend []*BackendConfig
    Queries struct {
        Insert []string
        Select []string
    }
}

var config Config

func (cfg *Config) String() string {
        return fmt.Sprint(*cfg)
}

func (cfg *Config) Set(value string) error {
    // Open config file
    file, err := os.Open(value)
    if err != nil {
            return err
    }
    defer file.Close()

    // Init new YAML decode
    d := yaml.NewDecoder(file)

    // Start YAML decoding from file
    if err := d.Decode(cfg); err != nil {
            return err
    }

    return nil
}

func init() {
    flag.Var(&config, "config", "YAML configuration file")
}

