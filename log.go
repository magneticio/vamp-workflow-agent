package main

import (
    "io"
    "os"
    "github.com/op/go-logging"
)

func createLogger() *logging.Logger {
    var logger = logging.MustGetLogger("vamp-workflow-agent")
    var backend = logging.NewLogBackend(io.Writer(os.Stdout), "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(
        "%{color}%{time:15:04:05.000} %{shortpkg:.4s} %{level:.4s} ==> %{message} %{color:reset}",
    ))
    logging.SetBackend(backendFormatter)
    return logger
}
