#!/usr/bin/env sh

echo "
██╗   ██╗ █████╗ ███╗   ███╗██████╗     ██╗    ██╗ ██████╗ ██████╗ ██╗  ██╗███████╗██╗      ██████╗ ██╗    ██╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗    ██║    ██║██╔═══██╗██╔══██╗██║ ██╔╝██╔════╝██║     ██╔═══██╗██║    ██║
██║   ██║███████║██╔████╔██║██████╔╝    ██║ █╗ ██║██║   ██║██████╔╝█████╔╝ █████╗  ██║     ██║   ██║██║ █╗ ██║
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝     ██║███╗██║██║   ██║██╔══██╗██╔═██╗ ██╔══╝  ██║     ██║   ██║██║███╗██║
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║         ╚███╔███╔╝╚██████╔╝██║  ██║██║  ██╗██║     ███████╗╚██████╔╝╚███╔███╔╝
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝          ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝
"

: "${VAMP_KEY_VALUE_STORE_PATH?not provided.}"
: "${VAMP_KEY_VALUE_STORE_TYPE?not provided.}"
: "${VAMP_KEY_VALUE_STORE_CONNECTION?not provided.}"

echo "VAMP_KEY_VALUE_STORE_TYPE       : ${VAMP_KEY_VALUE_STORE_TYPE}"
echo "VAMP_KEY_VALUE_STORE_CONNECTION : ${VAMP_KEY_VALUE_STORE_CONNECTION}"
echo "VAMP_KEY_VALUE_STORE_PATH       : ${VAMP_KEY_VALUE_STORE_PATH}"

printf "node  " && node --version
printf "npm   " && npm --version
/usr/bin/confd -version

mkdir -p /usr/local/vamp/confd/conf.d
mkdir -p /usr/local/vamp/confd/templates

echo "creating confd configuration and template"
cat <<EOT > /usr/local/vamp/confd/conf.d/workflow.toml
[template]
src = "workflow.tmpl"
dest = "/usr/local/vamp/workflow.js"
keys = [ "${VAMP_KEY_VALUE_STORE_PATH}" ]
EOT
cat <<EOT > /usr/local/vamp/confd/templates/workflow.tmpl
{{getv "${VAMP_KEY_VALUE_STORE_PATH}"}}
EOT

echo "running confd to retrieve workflow script"
/usr/bin/confd -onetime -backend \
        ${VAMP_KEY_VALUE_STORE_TYPE} \
        -node ${VAMP_KEY_VALUE_STORE_CONNECTION} \
        -confdir /usr/local/vamp/confd \
        || { exit 1; }

echo "running vamp workflow agent"
/usr/local/vamp/vamp-workflow-agent --workflow="/usr/local/vamp/workflow.js"
