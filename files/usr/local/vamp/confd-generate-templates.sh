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

cat <<EOT > "${dir_confd}/metricbeat.toml"
[template]
src = "metricbeat.tmpl"
dest = "/usr/local/metricbeat/metricbeat.yml"
EOT

cat <<EOT >> "${dir_templates}/metricbeat.tmpl"
metricbeat.modules:
- module: system
  metricsets:
    - cpu         # CPU stats
    - load        # System Load stats
    - filesystem  # Per filesystem stats
    - fsstat      # File system summary stats
    - memory      # Memory stats
    - network     # Network stats
    - process     # Per process stats
  enabled: true
  period: 1s
  processes: ['.*']
  tags: ["vamp","workflow","${VAMP_KEY_VALUE_STORE_PATH##*/}"]

output.elasticsearch:
  hosts: ["$VAMP_ELASTICSEARCH_URL"]
  index: "vamp-vwa-%{+yyyy-MM-dd}"
  template.path: /usr/local/metricbeat/metricbeat.template.json

path.home: /usr/local/metricbeat
path.config: \${path.home}
path.data: \${path.home}/data
path.logs: /var/log
EOT
