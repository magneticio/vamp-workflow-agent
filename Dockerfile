FROM magneticio/alpine-node:6.2.2

ENV VAMP_WORKFLOW_VERSION=0.9.0

ADD https://bintray.com/artifact/download/magnetic-io/downloads/vamp-workflow-agent/vamp-workflow-agent_${VAMP_WORKFLOW_VERSION}_linux_amd64.tar.gz /usr/local

RUN cd /usr/local/ && \
    tar xzvf vamp-workflow-agent_${VAMP_WORKFLOW_VERSION}_linux_amd64.tar.gz && \
    rm -Rf vamp-workflow-agent_${VAMP_WORKFLOW_VERSION}_linux_amd64.tar.gz

ENTRYPOINT ["/usr/local/vamp/vamp-workflow-agent"]