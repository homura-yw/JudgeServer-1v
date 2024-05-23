package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceUrl string `yaml:"service_url"`
	RpcPort    string `yaml:"rpc_port"`
	BufferSize int    `yaml:"buffer_size"`
}

var config Config

func init() {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	yamlFile, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(yamlFile, &config)
}
