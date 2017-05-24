package main

import (
	"container/ring"
	"log"
	"time"
)

type execution struct {
	ID     int                   `json:"id"`
	Start  time.Time             `json:"start"`
	Finish time.Time             `json:"finish"`
	Status int                   `json:"status"`
	Log    []executionLogMessage `json:"log"`
}

type executionLogMessage struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Error     bool      `json:"error"`
}

type executionStart struct {
	ID    int       `json:"id"`
	Start time.Time `json:"start"`
}

type executionFinish struct {
	ID     int       `json:"id"`
	Finish time.Time `json:"finish"`
	Status int       `json:"status"`
}

type executionLog struct {
	Execution int                 `json:"execution"`
	Log       executionLogMessage `json:"log"`
}

type executionService struct {
	executionCount               int
	executions, failedExecutions *ring.Ring
	stream                       chan interface{}
}

type executionCallbackHandler interface {
	Reply(message interface{})
}

func newExecutionService() *executionService {
	return &executionService{
		executionCount:   0,
		executions:       ring.New(16),
		failedExecutions: ring.New(16),
		stream:           make(chan interface{}),
	}
}

func (es *executionService) createExecution() *execution {
	es.executionCount++
	now := time.Now()
	e := &execution{es.executionCount, now, now, -1, make([]executionLogMessage, 0, 64)}

	es.executions.Value = e

	log.Print("New execution:", e.ID)
	es.stream <- executionStart{e.ID, e.Start}

	return e
}

func (es *executionService) finalizeExecution(e *execution, status int) {
	e.Status = status
	e.Finish = time.Now()

	es.executions = es.executions.Next()

	if status > 0 {
		es.failedExecutions.Value = e
		es.failedExecutions = es.failedExecutions.Next()
	}

	log.Print("Finalized execution:", e.ID)
	es.stream <- executionFinish{e.ID, e.Finish, e.Status}
}

func (es *executionService) executionLog(e *execution, log string, error bool) {
	message := executionLogMessage{len(e.Log), time.Now(), log, error}
	e.Log = append(e.Log, message)
	es.stream <- executionLog{e.ID, message}
}

func (es *executionService) command(command string, replyTo executionCallbackHandler) {
	if command == "execution-history" {
		es.executions.Do(func(e interface{}) {
			if e != nil {
				replyTo.Reply(*e.(*execution))
			}
		})
		es.failedExecutions.Do(func(e interface{}) {
			if e != nil {
				replyTo.Reply(*e.(*execution))
			}
		})
	}
}
