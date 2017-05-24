#! /bin/bash

# check required environment variables
: "${VAMP_KEY_VALUE_STORE_PATH:?VAMP_KEY_VALUE_STORE_PATH required.}"
: "${VAMP_KEY_VALUE_STORE_TYPE:?VAMP_KEY_VALUE_STORE_TYPE required.}"
: "${VAMP_KEY_VALUE_STORE_CONNECTION:?VAMP_KEY_VALUE_STORE_CONNECTION required.}"

mkdir -p /usr/local/vamp/confd/{conf.d,templates}

echo "Generating confd template resource and template for workflow.."
# Generate template resource for workflow
cat <<EOF > /usr/local/vamp/confd/conf.d/workflow.toml
[template]
src = "workflow.tmpl"
dest = "/usr/local/vamp/workflow.js"
keys = [ "${VAMP_KEY_VALUE_STORE_PATH}" ]
EOF
# Generate template for workflow
cat <<EOF > /usr/local/vamp/confd/templates/workflow.tmpl
{{getv "${VAMP_KEY_VALUE_STORE_PATH}"}}
EOF

echo "Fetching workflow from ${VAMP_KEY_VALUE_STORE_TYPE}.."
/usr/bin/confd \
  -onetime=true \
  -backend "${VAMP_KEY_VALUE_STORE_TYPE}" \
  -node ${VAMP_KEY_VALUE_STORE_CONNECTION//[,]/" -node "} \
  -log-level=warn \
  -confdir /usr/local/vamp/confd

if [ ! -f "$VAMP_WORKFLOW_PATH" ]; then
  echo "Failed to fetch workflow from $VAMP_KEY_VALUE_STORE_TYPE"
  exit 1
fi

if [ "$VAMP_WORKFLOW_AGENT_LOGO" = "TRUE" ] || [ "$VAMP_WORKFLOW_AGENT_LOGO" = "1" ]; then
echo "
██╗   ██╗ █████╗ ███╗   ███╗██████╗     ██╗    ██╗ ██████╗ ██████╗ ██╗  ██╗███████╗██╗      ██████╗ ██╗    ██╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗    ██║    ██║██╔═══██╗██╔══██╗██║ ██╔╝██╔════╝██║     ██╔═══██╗██║    ██║
██║   ██║███████║██╔████╔██║██████╔╝    ██║ █╗ ██║██║   ██║██████╔╝█████╔╝ █████╗  ██║     ██║   ██║██║ █╗ ██║
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝     ██║███╗██║██║   ██║██╔══██╗██╔═██╗ ██╔══╝  ██║     ██║   ██║██║███╗██║
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║         ╚███╔███╔╝╚██████╔╝██║  ██║██║  ██╗██║     ███████╗╚██████╔╝╚███╔███╔╝
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝          ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝
"
fi

echo "VAMP_KEY_VALUE_STORE_TYPE       : ${VAMP_KEY_VALUE_STORE_TYPE}"
echo "VAMP_KEY_VALUE_STORE_CONNECTION : ${VAMP_KEY_VALUE_STORE_CONNECTION}"
echo "VAMP_KEY_VALUE_STORE_PATH       : ${VAMP_KEY_VALUE_STORE_PATH}"

printf "node  " && node --version
printf "npm   " && npm --version

exec 2>&1
exec /usr/local/vamp/vamp-workflow-agent \
        --workflow="$VAMP_WORKFLOW_PATH" \
        --httpPort="$VAMP_WORKFLOW_HTTP_PORT" \
        --uiPath="$VAMP_WORKFLOW_UI_PATH"
