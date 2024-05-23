package main

import (
	"fmt"
	"sync"
	"time"

	"log"
	"os/exec"
	"strconv"

	"github.com/go-redis/redis"
)

func checkResult(result int, msg message, redisClient *redis.Client) {
	var output string
	defer func() {
		redisClient.Set(msg.SubmitId+":final", output, time.Second*1800)
	}()
	if result == 0 {
		log.Println("Accept")
		output = "Accept"
	} else if result == 1 {
		log.Println("Wrong Answer")
		output = "Wrong Answer"
	} else if result == 2 {
		log.Println("Time Limit Exceed")
		output = "Time Limit Exceed"
	} else if result == 3 {
		log.Println("Memory Limit Exceed")
		output = "Memory Limit Exceed"
	} else if result == 4 {
		log.Println("Complie Error")
		output = "Complie Error"
	} else if result == 5 {
		log.Println("Runtime Error")
		output = "Runtime Error"
	} else {
		output = "JudgeServer Error"
		log.Fatal("JudgeServer Error")
	}
}

func RunELF(redisClient *redis.Client, msg message) {
	path := ShiftPath(msg.ProblemType)
	path = fmt.Sprintf("/app/%s/main", path)
	if msg.IsContest == 0 {
		isAc := make([]int, msg.SubtestNum+1)
		for i := 1; i <= msg.SubtestNum; i += groupSize {
			var wg sync.WaitGroup
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				wg.Add(1)
				go func(offset int, isAc []int) {
					res, costTime, costMemory := 0, 0, 0
					cmd := exec.Command(
						path,
						strconv.Itoa(offset),
						strconv.Itoa(msg.TimeLimit),
						strconv.Itoa(msg.MemoryLimit),
					)
					output, err := cmd.CombinedOutput()
					if err != nil {
						log.Fatal(err)
						return
					}
					outputMsg := string(output)
					fmt.Sscanf(outputMsg, "%d %d %d", &res, &costTime, &costMemory)
					log.Printf("result:%v, cost_time:%v, cost_memory:%v\n", res, costTime, costMemory)
					isAc[offset] = res
					wg.Done()
				}(j, isAc)
			}
			wg.Wait()
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*1800)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*1800)
				if isAc[j] != 0 {
					checkResult(isAc[j], msg, redisClient)
					return
				}
			}
		}
		checkResult(0, msg, redisClient)
	} else {
		for i := 1; i <= msg.SubtestNum; i++ {
			res, costTime, costMemory := 0, 0, 0
			cmd := exec.Command(
				path,
				strconv.Itoa(i),
				strconv.Itoa(msg.TimeLimit),
				strconv.Itoa(msg.MemoryLimit),
			)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return
			}
			outputMsg := string(output)
			fmt.Sscanf(outputMsg, "%d %d %d", &res, &costTime, &costMemory)
			log.Printf("result:%v, cost_time:%v, cost_memory:%v\n", res, costTime, costMemory)
			redisClient.Set(msg.SubmitId, res, time.Second*1800)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*1800)
			if res != 0 {
				checkResult(res, msg, redisClient)
				return
			}
		}
		checkResult(0, msg, redisClient)
	}
}
