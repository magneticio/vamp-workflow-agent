package main

import (
	"testing"
)

func Test_newExecutionService(t *testing.T) {
	got := newExecutionService()
	if got.executionCount != 0 {
		t.Errorf("Unexpected count initialized: %d", got.executionCount)
	}

	if got.executions.Len() != 16 {
		t.Errorf("Unexpected executions buffer initialized: %d", got.executions.Len())
	}
}

func Test_executionService_createExecution(t *testing.T) {
	service := newExecutionService()

	var got *execution
	go func() {
		got = service.createExecution()
	}()

	assertMessage(t, service.stream, "executionStart")
	if got.ID != 1 {
		t.Errorf("ID: expected 1, got %d", got.ID)
	}

	go func() {
		got = service.createExecution()
	}()

	assertMessage(t, service.stream, "executionStart")
	if got.ID != 2 {
		t.Errorf("ID: expected 2, got %d", got.ID)
	}
}

func assertMessage(t *testing.T, stream <-chan interface{}, messageType string) {
	resp := <-stream
	switch tResp := resp.(type) {
	case executionStart:
		if messageType != "executionStart" {
			t.Errorf("Expected executionStart, got %T", tResp)
		}
	}
}
