package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/go-redis/redis"

	//#include "/app/NormalTest/main.h"
	//#include "/app/SpecialTest/main.h"
	//#include "/app/InteractiveTest/main.h"
	// #cgo LDFLAGS: -lseccomp
	// #cgo pkg-config: libseccomp
	"C"
)
import (
	"log"
	"os/exec"
	"strconv"
)

func checkResult(result int) {
	if result == 0 {
		log.Println("Accept")
	} else if result == 1 {
		log.Println("Wrong Answer")
	} else if result == 2 {
		log.Println("Time Limit Exceed")
	} else if result == 3 {
		log.Println("Memory Limit Exceed")
	} else if result == 4 {
		log.Println("Complie Error")
	} else if result == 5 {
		log.Println("Runtime Error")
	} else {
		log.Panic("JudgeServer Error")
	}
}

func RunCGO(redisClient *redis.Client, msg message) {
	if msg.IsContest == 0 {
		isAc := make([]int, msg.SubtestNum+1)
		for i := 1; i <= msg.SubtestNum; i += groupSize {
			var wg sync.WaitGroup
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				wg.Add(1)
				go func(offset int, isAc []int) {
					res, costTime, costMemory := 0, 0, 0
					if msg.ProblemType == 1 {
						C.NormalTest(
							(C.int)(msg.TimeLimit),
							(C.int)(msg.MemoryLimit),
							(C.int)(offset),
							(*C.int)(unsafe.Pointer(&res)),
							(*C.int)(unsafe.Pointer(&costTime)),
							(*C.int)(unsafe.Pointer(&costMemory)),
						)
					} else if msg.ProblemType == 2 {
						C.SpecialTest(
							(C.int)(msg.TimeLimit),
							(C.int)(msg.MemoryLimit),
							(C.int)(offset),
							(*C.int)(unsafe.Pointer(&res)),
							(*C.int)(unsafe.Pointer(&costTime)),
							(*C.int)(unsafe.Pointer(&costMemory)),
						)
					} else {
						C.InteractiveTest(
							(C.int)(msg.TimeLimit),
							(C.int)(msg.MemoryLimit),
							(C.int)(offset),
							(*C.int)(unsafe.Pointer(&res)),
							(*C.int)(unsafe.Pointer(&costTime)),
							(*C.int)(unsafe.Pointer(&costMemory)),
						)
					}
					isAc[offset] = res
					wg.Done()
				}(j, isAc)
			}
			wg.Wait()
			for j := i; j <= msg.SubtestNum && j < i+groupSize; j++ {
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*10)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*10)
				if isAc[j] != 0 {
					checkResult(isAc[j])
					return
				}
			}
		}
		checkResult(0)
	} else {
		for i := 1; i <= msg.SubtestNum; i++ {
			res, costTime, costMemory := 0, 0, 0
			if msg.ProblemType == 1 {
				C.NormalTest(
					(C.int)(msg.TimeLimit),
					(C.int)(msg.MemoryLimit),
					(C.int)(i),
					(*C.int)(unsafe.Pointer(&res)),
					(*C.int)(unsafe.Pointer(&costTime)),
					(*C.int)(unsafe.Pointer(&costMemory)),
				)
			} else if msg.ProblemType == 2 {
				C.SpecialTest(
					(C.int)(msg.TimeLimit),
					(C.int)(msg.MemoryLimit),
					(C.int)(i),
					(*C.int)(unsafe.Pointer(&res)),
					(*C.int)(unsafe.Pointer(&costTime)),
					(*C.int)(unsafe.Pointer(&costMemory)),
				)
			} else {
				C.InteractiveTest(
					(C.int)(msg.TimeLimit),
					(C.int)(msg.MemoryLimit),
					(C.int)(i),
					(*C.int)(unsafe.Pointer(&res)),
					(*C.int)(unsafe.Pointer(&costTime)),
					(*C.int)(unsafe.Pointer(&costMemory)),
				)
			}
			redisClient.Set(msg.SubmitId, res, time.Second*10)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*10)
			if res != 0 {
				checkResult(res)
				return
			}
		}
		checkResult(0)
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
						log.Panic(err)
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
				redisClient.Set(msg.SubmitId, isAc[j], time.Second*10)
				redisClient.Set(msg.SubmitId+"num", j, time.Second*10)
				if isAc[j] != 0 {
					checkResult(isAc[j])
					return
				}
			}
		}
		checkResult(0)
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
			redisClient.Set(msg.SubmitId, res, time.Second*10)
			redisClient.Set(msg.SubmitId+"num", i, time.Second*10)
			if res != 0 {
				checkResult(res)
				return
			}
		}
		checkResult(0)
	}
}
