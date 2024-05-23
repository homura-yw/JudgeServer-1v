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
	log.Println("service start")
	file, err := os.Open(configPath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	yamlFile, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	yaml.Unmarshal(yamlFile, &config)
	log.Printf("load success, source:%s", config.Source)
}
