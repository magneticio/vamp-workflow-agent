package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func run(service *executionService, workflowFile string) {

	if *executionPeriod == 0 {
		status := executeWorkflowScript(service, workflowFile)
		os.Exit(status)
		return
	}

	running := false

	execute := func() {
		if running {
			return
		}
		defer func() {
			running = false
		}()
		running = true
		executeWorkflowScript(service, workflowFile)
	}

	go execute()

	ticker := time.NewTicker(time.Duration(*executionPeriod) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			go execute()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func executeWorkflowScript(service *executionService, workflowFile string) int {
	log.Println("Executing workflow by Node.js.")

	e := service.createExecution()
	cmd := exec.Command("node", workflowFile)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go processOutput(service, e, stdout, false)
	go processOutput(service, e, stderr, true)

	finished := func(err error) int {
		exitStatusCode := 0
		if err != nil {
			log.Println("Error during execution of the workflow script.")
			if exitError, ok := err.(*exec.ExitError); ok {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					exitStatusCode = status.ExitStatus()
				}
			}
		}
		service.finalizeExecution(e, exitStatusCode)
		log.Println("Workflow execution took  :", e.Finish.Sub(e.Start))
		log.Println("Workflow exit status code:", e.Status)
		return e.Status
	}

	if *executionTimeout == 0 {
		return finished(cmd.Wait())
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Duration(*executionTimeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Println("ERROR - failed to kill the workflow process:", err)
		}
		log.Println("Workflow process is killed as timeout reached.")
		return 1
	case err := <-done:
		return finished(err)
	}
}

func processOutput(service *executionService, e *execution, rd io.Reader, error bool) {
	reader := bufio.NewReader(rd)
	for {
		input, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		workflowLog := strings.TrimSuffix(input, "\n")
		service.executionLog(e, workflowLog, error)
		message := "WORKFLOW - " + workflowLog
		if error {
			log.Println("ERROR - ", message)
		} else {
			log.Println(message)
		}
	}
}
