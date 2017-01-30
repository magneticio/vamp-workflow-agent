FROM magneticio/node:7.4-alpine

# https://github.com/peterbourgon/runsvinit
ENV RUNSVINIT_URL=https://github.com/peterbourgon/runsvinit/releases/download/v2.0.0/runsvinit-linux-amd64.tgz

ENV CONFD_URL=https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64

ENV METRICBEAT_VER=5.1.2
ENV METRICBEAT_URL=https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-${METRICBEAT_VER}-linux-x86_64.tar.gz

RUN set -xe \
    && apk add --no-cache \
      bash \
      curl \
      runit \
    && curl --location --silent --show-error $RUNSVINIT_URL --output - | tar zxf - -C /sbin \
    && chown 0:0 /sbin/runsvinit \
    && chmod 0775 /sbin/runsvinit \
    \
    && curl --location --silent --show-error --output /usr/bin/confd $CONFD_URL \
    && chmod 0755 /usr/bin/confd \
    \
    && curl --location --silent --show-error $METRICBEAT_URL --output - | tar zxf - -C /tmp \
    && mv /tmp/metricbeat-${METRICBEAT_VER}-linux-x86_64/metricbeat /usr/local/bin/ \
    && rm -rf /tmp/metricbeat-${METRICBEAT_VER}-linux-x86_64

ADD vamp-workflow-agent_katana_linux_amd64.tar.gz /usr/local
ADD files/ /

ENTRYPOINT ["/sbin/runsvinit"]
