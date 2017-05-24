package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	workflow         = flag.String("workflow", "/usr/local/vamp/workflow.js", "Path to workflow file.")
	httpPort         = flag.Int("httpPort", 8080, "HTTP port.")
	uiPath           = flag.String("uiPath", "./ui/", "Path to UI static content.")
	executionPeriod  = flag.Int("executionPeriod", -1, "Period between successive executions in seconds (0 if disabled).")
	executionTimeout = flag.Int("executionTimeout", -1, "Maximum allowed execution time in seconds (0 if no timeout).")
	help             = flag.Bool("help", false, "Print usage.")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds | log.Lshortfile)
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	uiAbsPath, err := filepath.Abs(*uiPath)

	if err != nil {
		log.Fatal(err)
	}

	checkInt(executionPeriod, "VAMP_WORKFLOW_EXECUTION_PERIOD", "Execution period must be specified and it must be equal or greater than 0.")
	checkInt(executionTimeout, "VAMP_WORKFLOW_EXECUTION_TIMEOUT", "Execution timeout must be specified and it must be equal or greater than 0.")

	log.Println("HTTP port                 :", *httpPort)
	log.Println("HTTP static content path  :", uiAbsPath)
	log.Println("Workflow file path        :", *workflow)
	log.Println("Workflow execution period :", *executionPeriod)
	log.Println("Workflow execution timeout:", *executionTimeout)

	api := &Api{make(chan interface{})}

	go run(api, *workflow)
	serve(api, *httpPort, uiAbsPath)
}

func checkInt(argument *int, environmentVariable, panic string) {
	if *argument < 0 {
		number, err := strconv.Atoi(os.Getenv(environmentVariable))
		if err != nil {
			log.Fatal(panic)
		}
		*argument = number
		if *argument < 0 {
			log.Fatal(panic)
		}
	}
}
