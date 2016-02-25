# Vamp Workflow Agent

Vamp Workflow Agent reads workflow Javascript file from key-value store [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/) and launches Node.js runtime to execute the file.

[![Build Status](https://travis-ci.org/magneticio/vamp-workflow-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-workflow-agent)

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
        Path to workflow files. (default "/opt/vamp/workflow")
```

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

By default in Docker image the following npm modules are installed (in `$workflowPath`) and available to workflow script:

- [underscore](https://github.com/jashkenas/underscore)
- [lodash](https://github.com/elastic/elasticsearch-js)
- [elasticsearch](https://github.com/elastic/elasticsearch-js)
- [node-zookeeper-client](https://github.com/alexguan/node-zookeeper-client)
- [consul](https://github.com/silas/node-consul)
- [node-etcd](https://github.com/stianeikeland/node-etcd)

More details [package.json](https://github.com/magneticio/vamp-workflow-agent/blob/master/package.json)
