# Vamp Workflow Agent

[![Join the chat at https://gitter.im/magneticio/vamp](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/magneticio/vamp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Docker](https://img.shields.io/badge/docker-images-blue.svg)](https://hub.docker.com/r/magneticio/vamp-workflow-agent/tags/)
[![Build Status](https://travis-ci.org/magneticio/vamp-workflow-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-workflow-agent)
[![Download](https://api.bintray.com/packages/magnetic-io/downloads/vamp-workflow-agent/images/download.svg) ](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent/_latestVersion)

- retrieves a workflow JavaScript file using [confd](https://github.com/kelseyhightower/confd)
- launches Node.js runtime to execute the script
- gives an overview (UI) of script execution

## Usage

```
$ ./vamp-workflow-agent -help
                                       
Usage of ./vamp-workflow-agent:
  -help
        Print usage.
  -httpPort int
        HTTP port. (default 8080)
  -uiPath string
        Path to UI static content. (default "./ui/")
  -workflow string
        Path to workflow file. (default "/usr/local/vamp/workflow.js")
  -executionPeriod int
        Period between successive executions in seconds (0 if disabled).
  -executionTimeout int
        Maximum allowed execution time in seconds (0 if no timeout).
```

Some arguments are mandatory and if they are not provided, agent will try to get them from environment variables. 
For environment variable names check out [Executing Workflow](https://github.com/magneticio/vamp-workflow-agent#executing-workflow).

### Metrics

The Vamp workflow agent docker image uses Metricbeat to collect performance metrics and ship them off to Elasticsearch. 
By default the [system module](https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-module-system.html) is configured to store metrics, with the additional tags to ease filtering:

- `vamp`
- `workflow`
- name of running workflow

## Building Binary

Using `make`:
```
make vamp-workflow-agent
```

Alternatively:

```
go get -d ./...
go install
CGO_ENABLED=0 go build -v -a -installsuffix cgo
```


Released binaries can be also [downloaded](https://bintray.com/magnetic-io/downloads/vamp-workflow-agent).
 
## Building Docker Images

Building the vamp-workflow-agent Docker image includes building the Go binary, downloading the vamp-node-client and building the workflow UI.

```
make
```

Docker images after the build: `magneticio/vamp-workflow-agent:katana`

For more details on available targets see the contents of the `Makefile`.

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
- `VAMP_ELASTICSEARCH_URL <=> http://elasticsearch:9200`

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
           -e VAMP_ELASTICSEARCH_URL=http://localhost:9200 \
           magneticio/vamp-workflow-agent:katana
```

In this example JavaScript is read from `/scripts` entry.
