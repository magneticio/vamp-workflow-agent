package main

import (
    "io"
    "os"
    "flag"
    "time"
    "strings"

    "github.com/op/go-logging"

    "golang.org/x/net/context"
    "github.com/coreos/etcd/client"
    "github.com/hashicorp/consul/api"
    "github.com/samuel/go-zookeeper/zk")

var (
    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    rootPath = flag.String("rootPath", "", "Scheduled workflow key-value store root path.")

    workflowPath = flag.String("workflowPath", "/opt/vamp/workflow", "Path to workflow files.")

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

func createLogger() *logging.Logger {
    var logger = logging.MustGetLogger("vamp-workflow-agent")
    var backend = logging.NewLogBackend(io.Writer(os.Stdout), "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(
        "%{color}%{time:15:04:05.000} %{shortpkg:.4s} %{level:.4s} ==> %{message} %{color:reset}",
    ))
    logging.SetBackend(backendFormatter)
    return logger
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

    workflow := *rootPath + "/workflow"
    logger.Info("Reading workflow from         : %s", workflow)

    response, err := readFromKeyValueStore(workflow)

    if err != nil {
        logger.Panic("Can't read the workflow: ", err)
        return
    }

    logger.Info("Value: %s", response)
}

func readFromKeyValueStore(path string) ([]byte, error) {
    if *storeType == "etcd" {
        return readFromEtcd(path)
    } else if *storeType == "consul" {
        return readFromConsul(path)
    } else if *storeType == "zookeeper" {
        return readFromZooKeeper(path)
    } else {
        logger.Panic("Key-value store type not supported: ", *storeType)
        return nil, nil
    }
}

func readFromEtcd(path string) ([]byte, error) {

    logger.Info("Initializing etcd connection: %s", *storeConnection)

    servers := strings.Split(*storeConnection, ",")
    cfg := client.Config{
        Endpoints:               servers,
        Transport:               client.DefaultTransport,
        HeaderTimeoutPerRequest: 5 * time.Second,
    }
    c, err := client.New(cfg)
    if err != nil {
        logger.Fatal(err)
        return nil, err
    }
    api := client.NewKeysAPI(c)

    logger.Info("etcd getting '%s' key value.", path)
    response, err := api.Get(context.Background(), path, nil)
    if err != nil {
        return nil, err
    } else {
        return []byte(response.Node.Value), nil
    }
}

func readFromConsul(path string) ([]byte, error) {

    logger.Info("Initializing Consul connection: %s", *storeConnection)

    conf := api.DefaultConfig()
    conf.Address = *storeConnection
    client, err := api.NewClient(conf)

    if err != nil {
        return nil, err
    }

    kv := client.KV()

    logger.Info("Consul getting '%s' key value.", path)
    pair, _, err := kv.Get(strings.TrimPrefix(path, "/"), nil)
    if err != nil {
        return nil, err
    }

    return pair.Value, nil
}

func readFromZooKeeper(path string) ([]byte, error) {

    logger.Info("Initializing ZooKeeper connection: %s", *storeConnection)

    servers := strings.Split(*storeConnection, ",")
    connection, _, err := zk.Connect(servers, 10 * time.Second)
    if err != nil {
        return nil, err
    }

    logger.Info("ZooKeeper getting '%s' key value.", path)
    response, _, err := connection.Get(path)
    if err != nil {
        return nil, err
    }

    return response, nil
}
