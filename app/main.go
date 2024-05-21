package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os/exec"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"
)

// 用于RPC调用的空结构体
type MessageQueueArgs struct{}

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

func GetOssClient(endpoint, ak, sk string) (client *oss.Client, err error) {
	client, err = oss.New(endpoint, ak, sk)
	return
}

const (
	endpoint        = "oss-cn-beijing.aliyuncs.com"
	accessKeyID     = "secret"
	accessKeySecret = "secret"
	bucketName      = "secret"
	rpcUrl          = "secret"
	compileError    = 4
	groupSize       = 3
	redisUrl        = "secret"
	redisPassword   = "secret"
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	return client
}

func run(bucket *oss.Bucket, msg message) {
	redisClient := newClient()
	defer redisClient.Close()

	if CompileCpp(bucket, &msg, redisClient) {
		return
	}

	RunELF(redisClient, msg)
}

func main() {
	client, err := rpc.Dial("tcp", rpcUrl)
	if err != nil {
		log.Panic("dialing error:", err)
	}
	ossClient, err := GetOssClient(
		endpoint,
		accessKeyID,
		accessKeySecret,
	)
	if err != nil {
		log.Panic("OSS Connect error:", err)
	}
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		log.Panic("URL Connect error", err)
	}
	for {
		args := MessageQueueArgs{}
		msg := MessageQueueReply{}
		err = client.Call("MessageQueue.Get", &args, &msg)
		if err != nil {
			log.Panic("Call error:", err)
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
		// if msg.ProblemType == 1 {
		// 	runNormal(bucket, message(msg))
		// } else if msg.ProblemType == 2 {
		// 	runSpecial(bucket, message(msg))
		// } else {
		// 	runInteractive(bucket, message(msg))
		// }
		run(bucket, message(msg))

		path := ShiftPath(msg.ProblemType)
		inputPath := fmt.Sprintf("/app/%s/*input*", path)
		answerPath := fmt.Sprintf("/app/%s/*answer*", path)
		outputPath := fmt.Sprintf("/app/%s/*output*", path)
		judgePath := fmt.Sprintf("/app/%s/*judge*", path)
		userPath := fmt.Sprintf("/app/%s/*user*", path)
		cmd := exec.Command(
			"sh",
			"-c",
			fmt.Sprintf("rm -rf %s %s %s %s %s", inputPath, answerPath, outputPath, judgePath, userPath),
		)
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("命令执行出错: %s\n", err)
			return
		}
		log.Println(string(cmdOutput))
	}
}
