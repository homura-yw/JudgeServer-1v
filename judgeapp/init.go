package main

import (
	"io"
	"log"
	"os"

	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceUrl string `yaml:"service_url"`
	RpcPort    string `yaml:"rpc_port"`
	BufferSize int    `yaml:"buffer_size"`
	Redis      struct {
		Url      string `yaml:"url"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Url,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})
	return client
}

var config Config
var redisClint *redis.Client

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
	redisClint = newClient()
}
