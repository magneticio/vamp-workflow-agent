# Vamp Workflow Agent

[![Join the chat at https://gitter.im/magneticio/vamp](https://badges.gitter.im/Join%20Chat.svg)](https://hub.docker.com/r/magneticio/vamp-workflow-agent/tags/)
[![Docker](https://img.shields.io/badge/docker-images-blue.svg)](https://img.shields.io/badge/docker-images-blue.svg)
[![Build Status](https://travis-ci.org/magneticio/vamp-workflow-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-workflow-agent)
[![Download](https://api.bintray.com/packages/magnetic-io/downloads/vamp-workflow-agent/images/download.svg) ](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent/_latestVersion)

- retrieves a workflow JavaScript file using [confd](https://github.com/kelseyhightower/confd)
- launches Node.js runtime to execute the script

## Usage

```
$ ./vamp-workflow-agent -help
                                       
Usage of ./vamp-workflow-agent:
  -help
        Print usage.
  -workflow string
        Path to workflow file. (default "/usr/local/vamp/workflow.js")
  -executionPeriod int
        Period between successive executions in seconds (0 if disabled).
  -executionTimeout int
        Maximum allowed execution time in seconds (0 if no timeout).
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

- retrieves workflow script
- saves it as `/usr/local/vamp/workflow.js`
- executes `node /usr/local/vamp/workflow.js`

Important environment variables:

- `VAMP_KEY_VALUE_STORE_TYPE <=> confd -backend`
- `VAMP_KEY_VALUE_STORE_CONNECTION <=> confd -node`
- `VAMP_KEY_VALUE_STORE_PATH <=> key used by confd`
- `VAMP_WORKFLOW_EXECUTION_PERIOD <=> $executionPeriod`
- `VAMP_WORKFLOW_EXECUTION_TIMEOUT <=> $executionTimeout`

Vamp JavaScript API [vamp-node-client](https://github.com/magneticio/vamp-node-client)

More details: [package.json](https://github.com/magneticio/vamp-workflow-agent/blob/master/package.json)

## Docker Images

Docker Hub [repo](https://hub.docker.com/r/magneticio/vamp-workflow-agent/).

Example:

```
docker run -e VAMP_KEY_VALUE_STORE_TYPE=zookeeper \
           -e VAMP_KEY_VALUE_STORE_CONNECTION=localhost:2181 \
           -e VAMP_KEY_VALUE_STORE_PATH=/scripts \
           -e VAMP_WORKFLOW_EXECUTION_PERIOD=0 \
           -e VAMP_WORKFLOW_EXECUTION_TIMEOUT=10 \
           magneticio/vamp-workflow-agent:katana
```

In this example JavaScript is read from `/scripts` entry.
