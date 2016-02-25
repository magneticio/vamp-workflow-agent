package main

import (
    "flag"
    "bytes"
    "os/exec"
    "io/ioutil"
    "syscall"
)

var (
    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    rootPath = flag.String("rootPath", "", "Scheduled workflow key-value store root path.")

    workflowPath = flag.String("workflowPath", "/opt/vamp", "Path to workflow files.")

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
        logger.Notice(printLogo("0.9.0"))
    }

    if *help {
        flag.Usage()
        return
    }

    if len(*storeType) == 0 {
        logger.Panic("Key-value store type not speciffed.")
        return
    }

    if len(*rootPath) == 0 {
        logger.Panic("Key-value store root key path not speciffed.")
        return
    }

    if len(*storeConnection) == 0 {
        logger.Panic("Key-value store connection not speciffed.")
        return
    }

    logger.Notice("Starting Vamp Workflow Agent")

    logger.Info("Key-value store type          : %s", *storeType)
    logger.Info("Key-value store connection    : %s", *storeConnection)
    logger.Info("Key-value store root key path : %s", *rootPath)
    logger.Info("Workflow file path            : %s", *workflowPath)

    workflowKey := *rootPath + "/workflow"
    logger.Info("Reading workflow from         : %s", workflowKey)

    content, err := readFromKeyValueStore(workflowKey)

    if err != nil {
        logger.Panic("Can't read the workflow: ", err)
        return
    }

    workflowFile := *workflowPath + "/workflow.js"
    logger.Info("Writing workflow script       : %s", workflowFile)
    err = ioutil.WriteFile(workflowFile, content, 0644)

    if err != nil {
        logger.Panic("Can't write to the workflow script: ", err)
        return
    }

    logger.Info("Executing 'workflow.js' by Node.js.")
    var cmd *exec.Cmd
    cmd = exec.Command("node", workflowFile)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err = cmd.Run()

    exitStatus := 0

    if err != nil {
        logger.Error("Error during execution: %s, %s", err.Error(), string(stderr.Bytes()[:]))
        if exitError, ok := err.(*exec.ExitError); ok {
            if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
                exitStatus = status.ExitStatus()
            }
        }

    } else {
        logger.Info("Execution standard output: %s", string(stdout.Bytes()[:]))
    }

    logger.Notice("Exit status: %d", exitStatus)
}
