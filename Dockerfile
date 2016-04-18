FROM alpine:3.3

ENV VAMP_WORKFLOW_VERSION=0.8.5

# Node.js installation is based on: https://github.com/mhart/alpine-node

ENV NODE_VERSION=v5.7.0 NPM_VERSION=3
ENV CONFIG_FLAGS="--fully-static --without-npm" DEL_PKGS="libgcc libstdc++" RM_DIRS=/usr/include

RUN set -ex && \
    apk add --no-cache curl make gcc g++ binutils-gold python linux-headers paxctl libgcc libstdc++ && \
    curl -sSL https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}.tar.gz | tar -xz && \
    cd /node-${NODE_VERSION} && \
    ./configure --prefix=/usr ${CONFIG_FLAGS} && \
    make -j$(grep -c ^processor /proc/cpuinfo 2>/dev/null || 1) && \
    make install && \
    paxctl -cm /usr/bin/node && \
    cd / && \
    if [ -x /usr/bin/npm ]; then \
      npm install -g npm@${NPM_VERSION} && \
      find /usr/lib/node_modules/npm -name test -o -name .bin -type d | xargs rm -rf; \
    fi && \
    apk del curl make gcc g++ binutils-gold python linux-headers paxctl ${DEL_PKGS} && \
    rm -rf /etc/ssl /node-${NODE_VERSION} ${RM_DIRS} \
      /usr/share/man /tmp/* /var/cache/apk/* /root/.npm /root/.node-gyp \
      /usr/lib/node_modules/npm/man /usr/lib/node_modules/npm/doc /usr/lib/node_modules/npm/html

ADD https://bintray.com/artifact/download/magnetic-io/downloads/vamp-workflow-agent/vamp-workflow-agent_${VAMP_WORKFLOW_VERSION}_linux_amd64.tar.gz /opt

ENTRYPOINT ["/opt/vamp/vamp-workflow-agent"]
