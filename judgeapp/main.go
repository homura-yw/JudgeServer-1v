package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

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

// channel的最大容量
const (
	MAXN       = 10000
	configPath = "conf.yaml"
)

type MessageQueue struct {
	ch      chan message
	results chan string
}

func (mq *MessageQueue) Push(msg message) error {
	mq.ch <- msg
	log.Println("push successfully!!")
	return nil
}

func (mq *MessageQueue) GetResults() []string {
	results := make([]string, 0)
	for i := 1; i <= config.BufferSize; i++ {
		IsEnd := false
		select {
		case msg := <-mq.results:
			results = append(results, msg)
		default:
			IsEnd = true
		}
		if IsEnd {
			break
		}
	}
	log.Printf("load %d messages\n", len(results))
	return results
}

func (mq *MessageQueue) Get(args *MessageQueueArgs, reply *MessageQueueReply) error {
	msg := <-mq.ch
	*reply = MessageQueueReply(msg)
	val, err := redisClint.Get(args.Key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("service redirect", err)
			mq.ch <- msg
			return nil
		} else {
			log.Panic("redis connect shutdown", val)
		}
	} else {
		num, err := strconv.Atoi(val)
		if err != nil {
			log.Panic("service error")
		}
		val = strconv.Itoa(num + 1)
		redisClint.Set(args.Key, val, 0)
	}
	log.Println("load successfully!!")
	return nil
}

func (mq *MessageQueue) PushResult(args *string, reply *string) error {
	*reply = "server get result successfully"
	mq.results <- *args
	log.Println("get result successfully")
	return nil
}

var MQ MessageQueue = MessageQueue{
	ch:      make(chan message, MAXN),
	results: make(chan string, MAXN),
}

func ProblemTest(ctx *gin.Context) {
	subtestNum, err := strconv.Atoi(ctx.PostForm("subtestNum"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
		return
	}
	memoryLimit, err := strconv.Atoi(ctx.PostForm("memoryLimit"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
		return
	}
	timeLimit, err := strconv.Atoi(ctx.PostForm("timeLimit"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
		return
	}
	isContest, err := strconv.Atoi(ctx.PostForm("isContest"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
		return
	}
	problemType, err := strconv.Atoi(ctx.PostForm("problemType"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
		return
	}
	msg := message{
		SubmitId:    ctx.PostForm("submitId"),
		TestUrl:     ctx.PostForm("testUrl"),
		CodeUrl:     ctx.PostForm("codeUrl"),
		SubtestNum:  subtestNum,
		MemoryLimit: memoryLimit,
		TimeLimit:   timeLimit,
		IsContest:   isContest,
		ProblemType: problemType,
	}
	MQ.Push(msg)
	ctx.String(http.StatusOK, "ok")
}

func getCompleteSubmission(ctx *gin.Context) {
	results := MQ.GetResults()
	ctx.JSON(http.StatusOK, results)
}

func test(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
}

func main() {
	go func() {
		rpc.Register(&MQ)
		listener, err := net.Listen("tcp", ":"+config.RpcPort)
		if err != nil {
			log.Panic("listen error:", err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Panic("accept error:", err)
			}
			go rpc.ServeConn(conn)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	log.Println("service start")

	router := gin.Default()

	router.POST("/send", ProblemTest)
	router.GET("/hello", test)
	router.GET("/pull", getCompleteSubmission)

	router.Run(config.ServiceUrl)
}
