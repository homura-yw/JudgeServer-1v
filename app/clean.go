package main

import (
	"fmt"
	"log"
	"os/exec"
)

func clear(msg message) {
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
		log.Panic(err)
		return
	}
	log.Println(string(cmdOutput))
}
