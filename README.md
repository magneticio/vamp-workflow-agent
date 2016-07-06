# Vamp Workflow Agent

Vamp Workflow Agent reads workflow JavaScript file from key-value store [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/) and launches Node.js runtime to execute the file.

[![Build Status](https://travis-ci.org/magneticio/vamp-workflow-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-workflow-agent)
[ ![Download](https://api.bintray.com/packages/magnetic-io/downloads/vamp-workflow-agent/images/download.svg) ](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent/_latestVersion)

## Usage

```
$ ./vamp-workflow-agent -help
                                       
Usage of ./vamp-workflow-agent:
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
        Path to workflow files. (default "/usr/local/vamp/workflow")
```

Some arguments are mandatory and if they are not provided, agent will try to get them from environment variables. 
For environment variable names check out [Executing Workflow](https://github.com/magneticio/vamp-workflow-agent#executing-workflow).

## Building Binary

Using the `build.sh` script:
```
  ./build.sh --make
```

Alternatively:

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `CGO_ENABLED=0 go build -v -a -installsuffix cgo`

Deliverable is in `target/go` directory.

Released binaries can be also [downloaded](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent).
 
## Building Docker Images

```
$ ./build.sh -h

Usage of ./build.sh:

  -h|--help   Help.
  -l|--list   List built Docker images.
  -r|--remove Remove Docker image.
  -m|--make   Build the binary and copy it to the Docker directories.
  -b|--build  Build Docker image.
  -a|--all    Build all binaries, by default only linux:amd64.
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

By default in Docker image the following npm packages are installed (in `$workflowPath`) and available to workflow script:

- [vamp-node-client](https://github.com/magneticio/vamp-node-client)

More details: [package.json](https://github.com/magneticio/vamp-workflow-agent/blob/master/package.json)

## Docker Images

Docker Hub [repo](https://hub.docker.com/r/magneticio/vamp-workflow-agent/).

[![](https://badge.imagelayers.io/magneticio/vamp-workflow-agent:0.9.0.svg)](https://imagelayers.io/?images=magneticio/vamp-workflow-agent:0.9.0)

Example:

```
docker run magneticio/vamp-workflow-agent:0.9.0 \
           -storeType=zookeeper \
           -storeConnection=localhost:2181 \
           -rootPath=/scripts
```

In this example JavaScript is read from `/scripts/workflow` key (`$rootPath/workflow`).
