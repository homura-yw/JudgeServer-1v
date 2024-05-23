package main

import (
	loadutil "judgeserver/loadUtil"
	"log"
	"net/rpc"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// 用于RPC调用的空结构体
type MessageQueueArgs struct {
	Key string
}

// RPC调用的响应类型，需要导出字段
type MessageQueueReply struct {
	SubmitId    string
	TestUrl     string
	CodeUrl     string
	SubtestNum  int
	MemoryLimit int
	TimeLimit   int
	IsContest   int
	ProblemType int
}

type message struct {
	SubmitId    string
	TestUrl     string
	CodeUrl     string
	SubtestNum  int
	MemoryLimit int
	TimeLimit   int
	IsContest   int
	ProblemType int
}

const (
	compileError = 4
	configPath   = "/app/conf.yaml"
	groupSize    = 3
	NORMAL       = 1
	SPECIAL      = 2
	INTERACTIVE  = 3
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Url,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})
	return client
}

func run(connection loadutil.Loadutil, msg message, redisClient *redis.Client) {
	if CompileCpp(connection, &msg, redisClient) {
		return
	}

	RunELF(redisClient, msg)
}

var UUID string

func main() {
	redisClient := newClient()
	UUID = uuid.New().String()
	defer func() {
		redisClient.Del(UUID)
	}()
	client, err := rpc.Dial("tcp", config.RpcUrl)
	if err != nil {
		log.Panic("dialing error:", err)
	}
	go register(UUID)
	connection, err := loadutil.LoadutilFactory(config)
	if err != nil {
		log.Panic("Connect error!!")
	}

	for {
		args := MessageQueueArgs{
			Key: UUID,
		}
		msg := MessageQueueReply{}
		err = client.Call("MessageQueue.Get", &args, &msg)
		if err != nil {
			log.Panic("Call error 1:", err)
		}
		log.Printf(
			"submid_id:%v\ntest_url:%v\ncode_url:%v\nsubtest_num:%v\nmemory_limit:%v\ntime_limit:%v\nis_contest:%v\nproblem_type:%v\ndate:%v\n",
			msg.SubmitId,
			msg.TestUrl,
			msg.CodeUrl,
			msg.SubtestNum,
			msg.MemoryLimit,
			msg.TimeLimit,
			msg.IsContest,
			msg.ProblemType,
			time.Now(),
		)
		run(connection, message(msg), redisClient)

		submitId := msg.SubmitId
		logs := ""
		err = client.Call("MessageQueue.PushResult", &submitId, &logs)
		if err != nil {
			log.Panic("connect error2:", err)
		}
		log.Println(logs)
		clear(message(msg))
	}
}
