FROM magneticio/alpine-node:6.2.2

ADD https://bintray.com/artifact/download/magnetic-io/downloads/vamp-workflow-agent/vamp-workflow-agent_${VAMP_WORKFLOW_VERSION}_linux_amd64.tar.gz /usr/local

ENTRYPOINT ["/usr/local/vamp/vamp-workflow-agent"]
