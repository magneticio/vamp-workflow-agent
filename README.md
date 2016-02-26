# Vamp Workflow Agent

Vamp Workflow Agent reads workflow JavaScript file from key-value store [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/) and launches Node.js runtime to execute the file.

[![Build Status](https://travis-ci.org/magneticio/vamp-workflow-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-workflow-agent)

## Usage

```
$ ./vamp-workflow-agent -help
                                       
Usage of ./vamp-workflow-agent:
  -elasticsearchConnection string
        Elasticsearch connection string.
  -help
        Print usage.
  -logo
        Show logo. (default true)
  -rootPath string
        Scheduled workflow key-value store root path.
  -storeConnection string
        Key-value store connection string.
  -storeType string
        zookeeper, consul or etcd.
  -workflowPath string
        Path to workflow files. (default "/opt/vamp/workflow")
```

Some arguments are mandatory and if they are not provided, agent will try to get them from environment variables. 
For environment variable names check out [Executing Workflow](#executing-workflow).

## Building Binary

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `CGO_ENABLED=0 go build -v -a -installsuffix cgo`

Alternatively using the `build.sh` script:
```
  ./build.sh --make
```
Deliverable is in `target/go` directory.

Released binaries can be [downloaded](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent).
 
## Building Docker Images

```
$ ./build.sh -h

Usage of ./build.sh:

  -h|--help   Help.
  -l|--list   List built Docker images.
  -r|--remove Remove Docker image.
  -m|--make   Build the binary and copy it to the Docker directories.
  -b|--build  Build Docker image.

```

## Executing Workflow

Vamp Workflow Agent:

- retrieves workflow script file from: `$rootPath/workflow`
- saves it as `workflow.js` in directory `$workflowPath`
- sets environment variables
- executes `node $workflowPath/workflow.js`

Set environment variables:

- `VAMP_KEY_VALUE_STORE_TYPE=$storeType`
- `VAMP_KEY_VALUE_STORE_CONNECTION=$storeConnection`
- `VAMP_KEY_VALUE_STORE_ROOT_PATH=$rootPath`
- `VAMP_WORKFLOW_DIRECTORY=$workflowPath`
- `VAMP_ELASTICSEARCH_CONNECTION=$elasticsearchConnection`

By default in Docker image the following npm packages are installed (in `$workflowPath`) and available to workflow script:

- [underscore](https://github.com/jashkenas/underscore)
- [lodash](https://github.com/elastic/elasticsearch-js)
- [highland](https://github.com/caolan/highland)
- [elasticsearch](https://github.com/elastic/elasticsearch-js)
- [node-zookeeper-client](https://github.com/alexguan/node-zookeeper-client)
- [consul](https://github.com/silas/node-consul)
- [node-etcd](https://github.com/stianeikeland/node-etcd)

More details: [package.json](https://github.com/magneticio/vamp-workflow-agent/blob/master/package.json)

## Docker Images

Docker Hub [repo](https://hub.docker.com/r/magneticio/vamp-workflow-agent/).

[![](https://badge.imagelayers.io/magneticio/vamp-workflow-agent:0.9.0.svg)](https://imagelayers.io/?images=magneticio/vamp-workflow-agent:0.9.0)

Example:

```
docker run magneticio/vamp-workflow-agent:0.9.0 \
           -elasticsearchConnection=localhost:9200 \
           -storeType=zookeeper \
           -storeConnection=localhost:2181 \
           -rootPath=/scripts
```

In this example JavaScript is read from `/scripts/workflow` key (`$rootPath/workflow`).
