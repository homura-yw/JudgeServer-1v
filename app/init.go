package main

import (
	"io"
	loadutil "judgeserver/loadUtil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var config loadutil.Config

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
