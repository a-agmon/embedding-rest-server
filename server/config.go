package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port          string `yaml:"port"`
	Host          string `yaml:"host"`
	EmbeddingFile string `yaml:"embedding_file"`
	ItemsFile     string `yaml:"items_file"`
	EmbeddingSize int    `yaml:"embedding_size"`
}

func LoadConfig() (*ServerConfig, error) {

	fileName := "server.config.yaml"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}
	var config ServerConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("cannot parse config file: %v", err)
	}
	return &config, nil
}
