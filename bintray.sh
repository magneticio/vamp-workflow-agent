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

for FILE in `find ${DIR}/target/docker | grep .tar.gz`
do
  DELIVERABLE=`basename "${FILE}"`

  if curl --output /dev/null --silent --head --fail "https://bintray.com/artifact/download/magnetic-io/downloads/vamp-workflow-agent/${DELIVERABLE}"; then
    echo "${DELIVERABLE} already uploaded"
  else
    echo "Uploading ${DELIVERABLE} to Bintray"
    curl -v -T ${DIR}/target/docker/${DELIVERABLE} \
     -u${BINTRAY_USER}:${BINTRAY_API_KEY} \
     -H "X-Bintray-Package:vamp-workflow-agent" \
     -H "X-Bintray-Version:$VERSION" \
     -H "X-Bintray-Publish:1" \
     https://api.bintray.com/content/magnetic-io/downloads/vamp-workflow-agent/${DELIVERABLE}
  fi
done
