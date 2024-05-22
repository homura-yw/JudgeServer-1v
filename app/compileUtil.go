package main

import (
	"fmt"
	loadutil "judgeserver/loadUtil"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func ShiftPath(ProblemType int) string {
	if ProblemType == 1 {
		return "NormalTest"
	} else if ProblemType == 2 {
		return "SpecialTest"
	} else {
		return "InteractiveTest"
	}
}

func CompileCpp(connection loadutil.Loadutil, msg *message, redisClient *redis.Client) bool {
	path := ShiftPath(msg.ProblemType)
	err := connection.LoadToFile(msg.CodeUrl, "/app/"+path+"/user.cpp")
	log.Println("start Compile")
	if err != nil {
		log.Println("user load error!!!")
		return true
	}
	cmd := exec.Command(
		"g++",
		"-o",
		"/app/"+path+"/user",
		"/app/"+path+"/user.cpp",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("编译失败: %s\n", err)
		log.Printf("输出: %s\n", output)
		redisClient.Set(msg.SubmitId, compileError, time.Second*10)
		return true
	}
	var loadGroup sync.WaitGroup
	for i := 1; i <= msg.SubtestNum; i++ {
		loadGroup.Add(1)
		go func(i int) {
			judgeUrl := fmt.Sprintf("%s/%d/judge.cpp", msg.TestUrl, i)
			judgePath := fmt.Sprintf("/app/%s/judge%d.cpp", path, i)
			err = connection.LoadToFile(judgeUrl, judgePath)
			if err != nil {
				log.Println("judge load error!")
				return
			}
			if msg.ProblemType == 1 || msg.ProblemType == 2 {
				answerUrl := fmt.Sprintf("%s/%d/answer", msg.TestUrl, i)
				answerPath := fmt.Sprintf("/app/%s/answer%d", path, i)
				err = connection.LoadToFile(answerUrl, answerPath)
				if err != nil {
					log.Println("answer load error!")
					return
				}

				inputUrl := fmt.Sprintf("%s/%d/input", msg.TestUrl, i)
				inputPath := fmt.Sprintf("/app/%s/input%d", path, i)
				err = connection.LoadToFile(inputUrl, inputPath)
				if err != nil {
					log.Println("input load error!")
					return
				}
			}
			judgeExec := fmt.Sprintf("/app/%s/judge%d", path, i)
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
	log.Println("complete compile!!!")
	return false
}
