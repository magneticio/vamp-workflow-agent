FROM node:7.3-alpine

ADD vamp-workflow-agent_katana_linux_amd64.tar.gz /usr/local
ADD https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64 /usr/bin/confd

RUN chmod u+x /usr/bin/confd /usr/local/vamp/vamp-workflow-agent.sh

ENTRYPOINT ["/usr/local/vamp/vamp-workflow-agent.sh"]
