package main

import (
    "time"
    "log"
    "container/ring"
)

type Execution struct {
    Id     int `json:"id"`
    Start  time.Time `json:"start"`
    Finish time.Time `json:"finish"`
    Status int `json:"status"`
    Log    []ExecutionLogMessage `json:"log"`
}

type ExecutionLogMessage struct {
    Id        int `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Message   string `json:"message"`
    Error     bool `json:"error"`
}

type ExecutionStart struct {
    Id    int `json:"id"`
    Start time.Time `json:"start"`
}

type ExecutionFinish struct {
    Id     int `json:"id"`
    Finish time.Time `json:"finish"`
    Status int `json:"status"`
}

type ExecutionLog struct {
    Execution int `json:"execution"`
    Log       ExecutionLogMessage `json:"log"`
}

var executionId int = 0
var executions = ring.New(16)
var failedExecutions = ring.New(16)

type Api struct {
    stream chan interface{}
}

type ApiReply interface {
    Reply(message interface{})
}

func (api *Api) CreateExecution() *Execution {
    executionId++
    now := time.Now()
    exe := &Execution{executionId, now, now, -1, make([]ExecutionLogMessage, 0, 64)}

    executions.Value = exe

    log.Print("New execution:", exe.Id)
    api.stream <- ExecutionStart{exe.Id, exe.Start}

    return exe
}

func (api *Api) FinalizeExecution(exe *Execution, status int) {
    exe.Status = status
    exe.Finish = time.Now()

    executions = executions.Next()

    if status > 0 {
        failedExecutions.Value = exe
        failedExecutions = failedExecutions.Next()
    }

    log.Print("Finalized execution:", exe.Id)
    api.stream <- ExecutionFinish{exe.Id, exe.Finish, exe.Status}
}

func (api *Api) ExecutionLog(exe *Execution, log string, error bool) {
    message := ExecutionLogMessage{len(exe.Log), time.Now(), log, error}
    exe.Log = append(exe.Log, message)
    api.stream <- ExecutionLog{exe.Id, message}
}

func (api *Api) Command(command string, replyTo ApiReply) {
    if command == "execution-history" {
        executions.Do(func(execution interface{}) {
            if execution != nil {
                replyTo.Reply(*execution.(*Execution))
            }
        })
        failedExecutions.Do(func(execution interface{}) {
            if execution != nil {
                replyTo.Reply(*execution.(*Execution))
            }
        })
    }
}