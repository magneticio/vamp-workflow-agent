package main

import (
    "flag"
    "bytes"
    "os/exec"
    "io/ioutil"
    "syscall"
    "os"
)

var (
    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    rootPath = flag.String("rootPath", "", "Scheduled workflow key-value store root path.")

    workflowPath = flag.String("workflowPath", "/opt/vamp", "Path to workflow files.")

    elasticsearchConnection = flag.String("elasticsearchConnection", "", "Elasticsearch connection string.")

    logo = flag.Bool("logo", true, "Show logo.")
    help = flag.Bool("help", false, "Print usage.")

    logger = createLogger()
)

func printLogo(version string) string {
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
        logger.Notice(printLogo("0.8.4"))
    }

    if *help {
        flag.Usage()
        return
    }

    check(storeType, "VAMP_KEY_VALUE_STORE_TYPE", "Key-value store type not specified.")
    check(rootPath, "VAMP_KEY_VALUE_STORE_ROOT_PATH", "Key-value store root key path not specified.")
    check(storeConnection, "VAMP_KEY_VALUE_STORE_CONNECTION", "Key-value store connection not specified.")

    logger.Notice("Starting Vamp Workflow Agent")

    logger.Info("Key-value store type          : %s", *storeType)
    logger.Info("Key-value store connection    : %s", *storeConnection)
    logger.Info("Key-value store root key path : %s", *rootPath)
    logger.Info("Workflow file path            : %s", *workflowPath)
    logger.Info("Elasticsearch connection      : %s", *elasticsearchConnection)

    workflowKey := *rootPath + "/workflow"
    logger.Info("Reading workflow from         : %s", workflowKey)

    content, err := readFromKeyValueStore(workflowKey)

    if err != nil {
        logger.Panic("Can't read the workflow: ", err)
        return
    }

    workflowFile := *workflowPath + "/workflow.js"

    err = writeWorkflowScript(workflowFile, content)

    if err != nil {
        logger.Panic("Can't write to the workflow script: ", err)
        return
    }

    err = setEnvironmentVariables()

    if err != nil {
        logger.Panic("Can't set eEnvironment variables: ", err)
        return
    }

    exitStatusCode := executeWorkflowScript(workflowFile)

    os.Exit(exitStatusCode)
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
    environmentVariables["VAMP_KEY_VALUE_STORE_ROOT_PATH"] = *rootPath
    environmentVariables["VAMP_WORKFLOW_DIRECTORY"] = *workflowPath
    environmentVariables["VAMP_ELASTICSEARCH_CONNECTION"] = *elasticsearchConnection

    for key, value := range environmentVariables {
        err := os.Setenv(key, value)
        if err != nil {
            return err
        }
    }

    return nil
}

func executeWorkflowScript(workflowFile string) int {

    logger.Info("Executing 'workflow.js' by Node.js.")
    var cmd *exec.Cmd
    cmd = exec.Command("node", workflowFile)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()

    exitStatusCode := 0

    logger.Info("Execution standard output: %s", string(stdout.Bytes()[:]))

    if err != nil {
        logger.Error("Error during execution: %s, %s", err.Error(), string(stderr.Bytes()[:]))
        if exitError, ok := err.(*exec.ExitError); ok {
            if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
                exitStatusCode = status.ExitStatus()
            }
        }
    }

    logger.Notice("Exit status code: %d", exitStatusCode)

    return exitStatusCode
}
