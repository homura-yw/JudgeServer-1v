package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/gin-gonic/gin"
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

// channel的最大容量
const MAXN = 10000

type MessageQueue struct {
	ch chan message
}

func (mq *MessageQueue) Push(msg message) error {
	mq.ch <- msg
	return nil
}

func (mq *MessageQueue) Get(args *MessageQueueArgs, reply *MessageQueueReply) error {
	msg := <-mq.ch
	*reply = MessageQueueReply(msg)
	return nil
}

var MQ MessageQueue = MessageQueue{
	ch: make(chan message, MAXN),
}

func ProblemTest(ctx *gin.Context) {
	subtestNum, err := strconv.Atoi(ctx.PostForm("subtestNum"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
	}
	memoryLimit, err := strconv.Atoi(ctx.PostForm("memoryLimit"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
	}
	timeLimit, err := strconv.Atoi(ctx.PostForm("timeLimit"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
	}
	isContest, err := strconv.Atoi(ctx.PostForm("isContest"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
	}
	problemType, err := strconv.Atoi(ctx.PostForm("problemType"))
	if err != nil {
		ctx.String(500, "Requst Failed:", err)
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

func test(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
}

func main() {
	go func() {
		rpc.Register(&MQ)
		listener, err := net.Listen("tcp", ":303")
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
	fmt.Println("hello world")
	router := gin.Default()
	router.POST("/send", ProblemTest)
	router.GET("/hello", test)
	router.Run("0.0.0.0:8080")
}
