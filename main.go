package main

import (
    "flag"
    "os"
    "io/ioutil"
    "strconv"
)

var (
    version string

    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    storePath = flag.String("storePath", "", "Key-value store path to workflow script.")
    filePath = flag.String("filePath", "/usr/local/vamp", "Path to workflow files.")
    executionPeriod = flag.Int("executionPeriod", -1, "Period between successive executions in seconds (0 if disabled).")
    executionTimeout = flag.Int("executionTimeout", -1, "Maximum allowed execution time in seconds (0 if no timeout).")

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

    checkString(storeType, "VAMP_KEY_VALUE_STORE_TYPE", "Key-value store type not specified.")
    checkString(storePath, "VAMP_KEY_VALUE_STORE_PATH", "Key-value store root key path not specified.")
    checkString(storeConnection, "VAMP_KEY_VALUE_STORE_CONNECTION", "Key-value store connection not specified.")

    checkInt(executionPeriod, "WORKFLOW_EXECUTION_PERIOD", "Execution period must be specified and it must be equal or greater than 0.")
    checkInt(executionTimeout, "WORKFLOW_EXECUTION_TIMEOUT", "Execution timeout must be specified and it must be equal or greater than 0.")

    logger.Notice("Starting Vamp Workflow Agent")

    logger.Info("Key-value store type          : %s", *storeType)
    logger.Info("Key-value store connection    : %s", *storeConnection)
    logger.Info("Key-value store root key path : %s", *storePath)
    logger.Info("Workflow file path            : %s", *filePath)
    logger.Info("Workflow execution period     : %d", *executionPeriod)
    logger.Info("Workflow execution timeout    : %d", *executionTimeout)

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

func checkString(argument *string, environmentVariable, panic string) {
    if len(*argument) == 0 {
        *argument = os.Getenv(environmentVariable)
        if len(*argument) == 0 {
            logger.Panic(panic)
        }
    }
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
