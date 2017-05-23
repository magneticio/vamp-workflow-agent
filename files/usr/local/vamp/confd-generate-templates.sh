#! /bin/bash

# Ensure we have our directories
dir_confd="/usr/local/vamp/confd/conf.d"
dir_templates="/usr/local/vamp/confd/templates"

mkdir -p "$dir_confd"
mkdir -p "$dir_templates"


# Generate config and templates for HAproxy
echo "creating confd configuration and template"
cat <<EOT > "${dir_confd}/workflow.toml"
[template]
src = "workflow.tmpl"
dest = "/usr/local/vamp/workflow.js"
keys = [ "${VAMP_KEY_VALUE_STORE_PATH}" ]
EOT
cat <<EOT > "${dir_templates}/workflow.tmpl"
{{getv "${VAMP_KEY_VALUE_STORE_PATH}"}}
EOT