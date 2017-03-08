package main

import (
    "log"
    "io"
    "syscall"
    "os/exec"
    "bufio"
    "strings"
    "time"
    "os"
)

func run(api *Api, workflowFile string) {

    if *executionPeriod == 0 {
        status := executeWorkflowScript(api, workflowFile)
        os.Exit(status)
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
        executeWorkflowScript(api, workflowFile)
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

func executeWorkflowScript(api *Api, workflowFile string) int {
    log.Println("Executing workflow by Node.js.")

    exe := api.CreateExecution()
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

    go processOutput(api, exe, stdout, false)
    go processOutput(api, exe, stderr, true)

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
        api.FinalizeExecution(exe, exitStatusCode)
        log.Println("Workflow execution took  :", exe.Finish.Sub(exe.Start))
        log.Println("Workflow exit status code:", exe.Status)
        return exe.Status
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

func processOutput(api *Api, exe *Execution, rd io.Reader, error bool) {
    reader := bufio.NewReader(rd)
    for {
        input, err := reader.ReadString('\n')
        if err != nil || err == io.EOF {
            break
        }
        workflowLog := strings.TrimSuffix(input, "\n")
        api.ExecutionLog(exe, workflowLog, error)
        message := "WORKFLOW - " + workflowLog
        if error {
            log.Println("ERROR - ", message);
        } else {
            log.Println(message);
        }
    }
}
