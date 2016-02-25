package main

import (
    "time"
    "strings"

    "golang.org/x/net/context"
    "github.com/coreos/etcd/client"
    "github.com/hashicorp/consul/api"
    "github.com/samuel/go-zookeeper/zk"
)

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
