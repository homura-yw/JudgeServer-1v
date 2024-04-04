package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os/exec"
	"sync"
	"time"
	"unsafe"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"

	//#include "/app/NormalTest/main.h"
	//#include "/app/SpecialTest/main.h"
	//#include "/app/InteractiveTest/main.h"
	// #cgo LDFLAGS: -lseccomp
	// #cgo pkg-config: libseccomp
	"C"
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
	accessKeyID     = "accessKeyID"
	accessKeySecret = "accessKeySecret"
	bucketName      = "bucketName"
	rpcUrl          = "localhost:303"
	compileError    = 4
	groupSize       = 3
	redisUrl        = "redisUrl"
	redisPassword   = "redis"
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	return client
}

func runNormal(bucket *oss.Bucket, msg message) {
	redisClient := newClient()
	defer redisClient.Close()
	err := bucket.GetObjectToFile(msg.CodeUrl, "/app/NormalTest/user.cpp")
	if err != nil {
		return
	}
	cmd := exec.Command(
		"g++",
		"-o",
		"/app/NormalTest/user",
		"/app/NormalTest/user.cpp",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("编译失败: %s\n", err)
		fmt.Printf("输出: %s\n", output)
		redisClient.Set(msg.SubmitId, compileError, time.Second*10)
		return
	}
	fmt.Println(msg.CodeUrl)
	var loadGroup sync.WaitGroup
	for i := 1; i <= msg.SubtestNum; i++ {
		loadGroup.Add(1)
		go func(i int) {
			judgeUrl := fmt.Sprintf("%s/%d/judge.cpp", msg.TestUrl, i)
			answerUrl := fmt.Sprintf("%s/%d/answer", msg.TestUrl, i)
			inputUrl := fmt.Sprintf("%s/%d/input", msg.TestUrl, i)

			judgePath := fmt.Sprintf("/app/NormalTest/judge%d.cpp", i)
			answerPath := fmt.Sprintf("/app/NormalTest/answer%d", i)
			inputPath := fmt.Sprintf("/app/NormalTest/input%d", i)

			err = bucket.GetObjectToFile(judgeUrl, judgePath)
			if err != nil {
				return
			}
			err = bucket.GetObjectToFile(answerUrl, answerPath)
			if err != nil {
				return
			}
			err = bucket.GetObjectToFile(inputUrl, inputPath)
			if err != nil {
				return
			}

			judgeExec := fmt.Sprintf("/app/NormalTest/judge%d", i)
			cmd := exec.Command(
				"g++",
				"-o",
				judgeExec,
				judgePath,
			)
			cmd.CombinedOutput()
			loadGroup.Done()
		}(i)
	}
	loadGroup.Wait()
	if msg.IsContest == 1 {
		isAc := make([]int, msg.SubtestNum+1)
		for i := 1; i <= msg.SubtestNum; i += groupSize {
			var wg sync.WaitGroup
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				wg.Add(1)
				go func(offset int, isAc []int) {
					res, costTime, costMemory := 0, 0, 0
					C.NormalTest(
						(C.int)(msg.TimeLimit),
						(C.int)(msg.MemoryLimit),
						(C.int)(j),
						(*C.int)(unsafe.Pointer(&res)),
						(*C.int)(unsafe.Pointer(&costTime)),
						(*C.int)(unsafe.Pointer(&costMemory)),
					)
					isAc[offset] = res
					wg.Done()
				}(j, isAc)
			}
			wg.Wait()
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*10)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*10)
				if isAc[j] != 0 {
					fmt.Println("Wrong Answer")
					return
				}
			}
		}
		fmt.Println("Accept")
	} else {
		for i := 1; i <= msg.SubtestNum; i++ {
			res, costTime, costMemory := 0, 0, 0
			C.NormalTest(
				(C.int)(msg.TimeLimit),
				(C.int)(msg.MemoryLimit),
				(C.int)(i),
				(*C.int)(unsafe.Pointer(&res)),
				(*C.int)(unsafe.Pointer(&costTime)),
				(*C.int)(unsafe.Pointer(&costMemory)),
			)
			redisClient.Set(msg.SubmitId, res, time.Second*10)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*10)
			if res != 0 {
				fmt.Println("Wrong Answer")
				return
			}
		}
		fmt.Println("Accept")
	}
}

func runSpecial(bucket *oss.Bucket, msg message) {
	redisClient := newClient()
	defer redisClient.Close()
	err := bucket.GetObjectToFile(msg.CodeUrl, "/app/SpecialTest/user.cpp")
	if err != nil {
		return
	}
	cmd := exec.Command(
		"g++",
		"-o",
		"/app/SpecialTest/user",
		"/app/SpecialTest/user.cpp",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("编译失败: %s\n", err)
		fmt.Printf("输出: %s\n", output)
		redisClient.Set(msg.SubmitId, compileError, time.Second*10)
		return
	}
	fmt.Println(msg.CodeUrl)
	var loadGroup sync.WaitGroup
	for i := 1; i <= msg.SubtestNum; i++ {
		loadGroup.Add(1)
		go func(i int) {
			judgeUrl := fmt.Sprintf("%s/%d/judge.cpp", msg.TestUrl, i)
			answerUrl := fmt.Sprintf("%s/%d/answer", msg.TestUrl, i)
			inputUrl := fmt.Sprintf("%s/%d/input", msg.TestUrl, i)

			judgePath := fmt.Sprintf("/app/SpecialTest/judge%d.cpp", i)
			answerPath := fmt.Sprintf("/app/SpecialTest/answer%d", i)
			inputPath := fmt.Sprintf("/app/SpecialTest/input%d", i)

			err = bucket.GetObjectToFile(judgeUrl, judgePath)
			if err != nil {
				return
			}
			err = bucket.GetObjectToFile(answerUrl, answerPath)
			if err != nil {
				return
			}
			err = bucket.GetObjectToFile(inputUrl, inputPath)
			if err != nil {
				return
			}

			judgeExec := fmt.Sprintf("/app/SpecialTest/judge%d", i)
			cmd := exec.Command(
				"g++",
				"-o",
				judgeExec,
				judgePath,
			)
			cmd.CombinedOutput()
			loadGroup.Done()
		}(i)
	}
	loadGroup.Wait()
	if msg.IsContest == 1 {
		isAc := make([]int, msg.SubtestNum+1)
		for i := 1; i <= msg.SubtestNum; i += groupSize {
			var wg sync.WaitGroup
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				wg.Add(1)
				go func(offset int, isAc []int) {
					res, costTime, costMemory := 0, 0, 0
					C.SpecialTest(
						(C.int)(msg.TimeLimit),
						(C.int)(msg.MemoryLimit),
						(C.int)(j),
						(*C.int)(unsafe.Pointer(&res)),
						(*C.int)(unsafe.Pointer(&costTime)),
						(*C.int)(unsafe.Pointer(&costMemory)),
					)
					isAc[offset] = res
					wg.Done()
				}(j, isAc)
			}
			wg.Wait()
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*10)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*10)
				if isAc[j] != 0 {
					fmt.Println("Wrong Answer")
					return
				}
			}
		}
		fmt.Println("Accept")
	} else {
		for i := 1; i <= msg.SubtestNum; i++ {
			res, costTime, costMemory := 0, 0, 0
			C.SpecialTest(
				(C.int)(msg.TimeLimit),
				(C.int)(msg.MemoryLimit),
				(C.int)(i),
				(*C.int)(unsafe.Pointer(&res)),
				(*C.int)(unsafe.Pointer(&costTime)),
				(*C.int)(unsafe.Pointer(&costMemory)),
			)
			redisClient.Set(msg.SubmitId, res, time.Second*10)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*10)
			if res != 0 {
				fmt.Println("Wrong Answer")
				return
			}
		}
		fmt.Println("Accept")
	}
}

func runInteractive(bucket *oss.Bucket, msg message) {
	redisClient := newClient()
	defer redisClient.Close()
	err := bucket.GetObjectToFile(msg.CodeUrl, "/app/InteractiveTest/user.cpp")
	if err != nil {
		return
	}
	cmd := exec.Command(
		"g++",
		"-o",
		"/app/InteractiveTest/user",
		"/app/InteractiveTest/user.cpp",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("编译失败: %s\n", err)
		fmt.Printf("输出: %s\n", output)
		redisClient.Set(msg.SubmitId, compileError, time.Second*10)
		return
	}
	var loadGroup sync.WaitGroup
	for i := 1; i <= msg.SubtestNum; i++ {
		loadGroup.Add(1)
		go func(i int) {
			judgeUrl := fmt.Sprintf("%s/%d/judge.cpp", msg.TestUrl, i)

			judgePath := fmt.Sprintf("/app/InteractiveTest/judge%d.cpp", i)

			err = bucket.GetObjectToFile(judgeUrl, judgePath)
			if err != nil {
				return
			}

			judgeExec := fmt.Sprintf("/app/InteractiveTest/judge%d", i)
			cmd := exec.Command(
				"g++",
				"-o",
				judgeExec,
				judgePath,
			)
			cmd.CombinedOutput()
			loadGroup.Done()
		}(i)
	}
	loadGroup.Wait()
	if msg.IsContest == 1 {
		isAc := make([]int, msg.SubtestNum+1)
		for i := 1; i <= msg.SubtestNum; i += groupSize {
			var wg sync.WaitGroup
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				wg.Add(1)
				go func(offset int, isAc []int) {
					res, costTime, costMemory := 0, 0, 0
					C.InteractiveTest(
						(C.int)(msg.TimeLimit),
						(C.int)(msg.MemoryLimit),
						(C.int)(j),
						(*C.int)(unsafe.Pointer(&res)),
						(*C.int)(unsafe.Pointer(&costTime)),
						(*C.int)(unsafe.Pointer(&costMemory)),
					)
					isAc[offset] = res
					wg.Done()
				}(j, isAc)
			}
			wg.Wait()
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*10)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*10)
				if isAc[j] != 0 {
					fmt.Println("Wrong Answer")
					return
				}
			}
		}
		fmt.Println("Accept")
	} else {
		for i := 1; i <= msg.SubtestNum; i++ {
			res, costTime, costMemory := 0, 0, 0
			C.InteractiveTest(
				(C.int)(msg.TimeLimit),
				(C.int)(msg.MemoryLimit),
				(C.int)(i),
				(*C.int)(unsafe.Pointer(&res)),
				(*C.int)(unsafe.Pointer(&costTime)),
				(*C.int)(unsafe.Pointer(&costMemory)),
			)
			redisClient.Set(msg.SubmitId, res, time.Second*10)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*10)
			if res != 0 {
				fmt.Println("Wrong Answer")
				return
			}
		}
		fmt.Println("Accept")
	}
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
		fmt.Println(
			msg.SubmitId,
			msg.TestUrl,
			msg.CodeUrl,
			msg.SubtestNum,
			msg.MemoryLimit,
			msg.TimeLimit,
			msg.IsContest,
			msg.ProblemType,
		)
		if msg.ProblemType == 1 {
			runNormal(bucket, message(msg))
		} else if msg.ProblemType == 2 {
			runSpecial(bucket, message(msg))
		} else {
			runInteractive(bucket, message(msg))
		}
	}
}
