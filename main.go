package main

import (
    "flag"
    "os"
    "strconv"
)

var (
    version string
    workflow = flag.String("workflow", "/usr/local/vamp/workflow.js", "Path to workflow file.")
    executionPeriod = flag.Int("executionPeriod", -1, "Period between successive executions in seconds (0 if disabled).")
    executionTimeout = flag.Int("executionTimeout", -1, "Maximum allowed execution time in seconds (0 if no timeout).")
    help = flag.Bool("help", false, "Print usage.")
    logger = createLogger()
)

func main() {

    flag.Parse()

    if *help {
        flag.Usage()
        return
    }

    checkInt(executionPeriod, "VAMP_WORKFLOW_EXECUTION_PERIOD", "Execution period must be specified and it must be equal or greater than 0.")
    checkInt(executionTimeout, "VAMP_WORKFLOW_EXECUTION_TIMEOUT", "Execution timeout must be specified and it must be equal or greater than 0.")

    logger.Notice("Starting Vamp Workflow Agent")

    logger.Info("Workflow file path            : %s", *workflow)
    logger.Info("Workflow execution period     : %d", *executionPeriod)
    logger.Info("Workflow execution timeout    : %d", *executionTimeout)

    run(*workflow)
}

func checkInt(argument *int, environmentVariable, panic string) {
    if *argument < 0 {
        number, err := strconv.Atoi(os.Getenv(environmentVariable))
        if err != nil {
            logger.Panic(panic)
        }
        *argument = number
        if *argument < 0 {
            logger.Panic(panic)
        }
    }
}
