#!/usr/bin/env bash

if [ "$#" -ne 1 ]; then
    SCRIPTNAME=$(basename "$0")
    echo "Usage: ${SCRIPTNAME} <version>"
    exit 1
fi

VERSION=$1

: ${BINTRAY_USER:?"No BINTRAY_USER set"}
: ${BINTRAY_API_KEY:?"No BINTRAY_API_KEY set"}
: ${VERSION:?"Not set"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

DELIVERABLE=vamp-workflow-agent_${VERSION}_linux_amd64.tar.gz

echo "Uploading ${DELIVERABLE} to Bintray"

curl -v -T ${DIR}/target/docker/vamp.tar.gz \
 -u${BINTRAY_USER}:${BINTRAY_API_KEY} \
 -H "X-Bintray-Package:vamp-workflow-agent" \
 -H "X-Bintray-Version:$VERSION" \
 -H "X-Bintray-Publish:1" \
 https://api.bintray.com/content/magnetic-io/downloads/vamp-workflow-agen/${DELIVERABLE}
