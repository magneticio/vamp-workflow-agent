package main

import (
    "log"
    "io"
    "syscall"
    "os/exec"
    "bufio"
    "strings"
    "time"
)

func run(workflowFile string) {

    if *executionPeriod == 0 {
        executeWorkflowScript(workflowFile)
        return
    }

    running := false

    execute := func() {
        if (running) {
            return
        }
        defer func() {
            running = false
        }()
        running = true
        executeWorkflowScript(workflowFile)
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

func executeWorkflowScript(workflowFile string) {
    logger.Notice("Executing workflow by Node.js.")
    start := time.Now()
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

    go processOutput(stdout, false)
    go processOutput(stderr, true)

    finished := func(err error) {
        exitStatusCode := 0
        if err != nil {
            logger.Error("Error during execution of the workflow script.")
            if exitError, ok := err.(*exec.ExitError); ok {
                if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
                    exitStatusCode = status.ExitStatus()
                }
            }
        }
        elapsed := time.Since(start)
        logger.Notice("Workflow execution took: %s", elapsed)
        logger.Notice("Workflow exit status code: %d", exitStatusCode)
    }

    if *executionTimeout == 0 {
        finished(cmd.Wait())
        return
    }

    done := make(chan error, 1)
    go func() {
        done <- cmd.Wait()
    }()
    select {
    case <-time.After(time.Duration(*executionTimeout) * time.Second):
        if err := cmd.Process.Kill(); err != nil {
            logger.Error("Failed to kill the workflow process: ", err)
        }
        logger.Notice("Workflow process is killed as timeout reached.")
    case err := <-done:
        finished(err)
    }
}

func processOutput(rd io.Reader, error bool) {
    reader := bufio.NewReader(rd)
    for {
        input, err := reader.ReadString('\n')
        if err != nil || err == io.EOF {
            break
        }
        message := "WORKFLOW " + strings.TrimSuffix(input, "\n")
        if error {
            logger.Error(message);
        } else {
            logger.Info(message);
        }
    }
}