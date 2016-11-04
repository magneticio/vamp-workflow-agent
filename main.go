package main

import (
    "flag"
    "os"
    "log"
    "io"
    "syscall"
    "os/exec"
    "io/ioutil"
    "bufio"
    "strings"
)

var (
    version string

    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    storePath = flag.String("storePath", "", "Key-value store path to workflow script.")
    filePath = flag.String("filePath", "/usr/local/vamp", "Path to workflow files.")
    executionPeriod = flag.String("executionPeriod", "0", "Period between successive executions in seconds (0 if disabled).")
    executionTimeout = flag.String("executionTimeout", "5", "Maximum allowed execution time in seconds.")

    logo = flag.Bool("logo", true, "Show logo.")
    help = flag.Bool("help", false, "Print usage.")

    logger = createLogger()
)

func printLogo() string {
    return `
██╗   ██╗ █████╗ ███╗   ███╗██████╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗
██║   ██║███████║██╔████╔██║██████╔╝
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝
                       workflow agent
                       version ` + version + `
                       by magnetic.io
                                      `
}

func main() {

    flag.Parse()

    if *logo {
        logger.Notice(printLogo())
    }

    if *help {
        flag.Usage()
        return
    }

    check(storeType, "VAMP_KEY_VALUE_STORE_TYPE", "Key-value store type not specified.")
    check(storePath, "VAMP_KEY_VALUE_STORE_PATH", "Key-value store root key path not specified.")
    check(storeConnection, "VAMP_KEY_VALUE_STORE_CONNECTION", "Key-value store connection not specified.")
    check(executionPeriod, "WORKFLOW_EXEUTION_PERIOD", "Execution period not specified.")
    check(executionTimeout, "WORKFLOW_EXEUTION_TIMEOUT", "Execution timeout not specified.")

    logger.Notice("Starting Vamp Workflow Agent")

    logger.Info("Key-value store type          : %s", *storeType)
    logger.Info("Key-value store connection    : %s", *storeConnection)
    logger.Info("Key-value store root key path : %s", *storePath)
    logger.Info("Workflow file path            : %s", *filePath)
    logger.Info("Workflow execution period     : %s", *executionPeriod)
    logger.Info("Workflow execution timeout    : %s", *executionTimeout)

    workflowKey := *storePath
    logger.Info("Reading workflow from         : %s", workflowKey)

    content, err := readFromKeyValueStore(workflowKey)

    if err != nil {
        logger.Panic("Can't read the workflow: ", err)
        return
    }

    workflowFile := *filePath + "/workflow.js"

    err = writeWorkflowScript(workflowFile, content)

    if err != nil {
        logger.Panic("Can't write to the workflow script: ", err)
        return
    }

    err = setEnvironmentVariables()

    if err != nil {
        logger.Panic("Can't set environment variables: ", err)
        return
    }

    run(workflowFile)
}

func check(argument *string, environmentVariable, panic string) {
    if len(*argument) == 0 {
        *argument = os.Getenv(environmentVariable)
        if len(*argument) == 0 {
            logger.Panic(panic)
        }
    }
}

func writeWorkflowScript(workflowFile string, content []byte) error {
    logger.Info("Writing workflow script       : %s", workflowFile)
    return ioutil.WriteFile(workflowFile, content, 0644)
}

func setEnvironmentVariables() error {

    environmentVariables := make(map[string]string)

    environmentVariables["VAMP_KEY_VALUE_STORE_TYPE"] = *storeType
    environmentVariables["VAMP_KEY_VALUE_STORE_CONNECTION"] = *storeConnection
    environmentVariables["VAMP_KEY_VALUE_STORE_PATH"] = *storePath
    environmentVariables["VAMP_WORKFLOW_DIRECTORY"] = *filePath

    for key, value := range environmentVariables {
        err := os.Setenv(key, value)
        if err != nil {
            return err
        }
    }

    return nil
}

func run(workflowFile string) {

    exitStatusCode := executeWorkflowScript(workflowFile)

    os.Exit(exitStatusCode)
}

func executeWorkflowScript(workflowFile string) int {

    logger.Info("Executing 'workflow.js' by Node.js.")
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

    go processOutput(bufio.NewReader(stdout), false)
    go processOutput(bufio.NewReader(stderr), true)

    exitStatusCode := 0

    err = cmd.Wait()
    if err != nil {
        logger.Error("Error during execution of the workflow script.")
        if exitError, ok := err.(*exec.ExitError); ok {
            if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
                exitStatusCode = status.ExitStatus()
            }
        }
    }

    logger.Notice("Workflow exit status code: %d", exitStatusCode)

    return exitStatusCode
}

func processOutput(reader *bufio.Reader, error bool) {
    for {
        input, err := reader.ReadString('\n')
        if err != nil && err == io.EOF {
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