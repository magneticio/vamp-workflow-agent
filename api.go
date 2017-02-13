package main

import (
    "time"
    "log"
)

type Execution struct {
    Id     int `json:"id"`
    Start  time.Time `json:"start"`
    Finish time.Time `json:"finish"`
    Status int `json:"status"`
    Log    []ExecutionLogMessage `json:"log"`
}

type ExecutionLogMessage struct {
    Id      int `json:"id"`
    Message string `json:"message"`
    Error   bool `json:"error"`
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

type Api struct {
    stream chan interface{}
}

func (api *Api) CreateExecution() Execution {
    executionId++
    now := time.Now()
    exe := Execution{executionId, now, now, -1, make([]ExecutionLogMessage, 0, 64)}

    log.Print("New execution:", exe.Id)
    api.stream <- ExecutionStart{exe.Id, exe.Start}

    return exe
}

func (api *Api) FinalizeExecution(exe *Execution, status int) {
    (*exe).Status = status
    (*exe).Finish = time.Now()

    log.Print("Finalized execution:", (*exe).Id)
    api.stream <- ExecutionFinish{(*exe).Id, (*exe).Finish, (*exe).Status}
}

func (api *Api) ExecutionLog(exe *Execution, log string, error bool) {
    message := ExecutionLogMessage{len((*exe).Log), log, error}
    (*exe).Log = append((*exe).Log, message)
    api.stream <- ExecutionLog{(*exe).Id, message}
}
