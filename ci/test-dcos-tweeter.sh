#!/usr/bin/env bash

# Installs and tests Tweeter on DC/OS.
# Requires dcos CLI to be installed, configured, and logged in.
#
# Usage:
# $ ci/test-tweeter.sh

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

project_dir=$(cd "$(dirname "${BASH_SOURCE}")/.." && pwd -P)
cd "${project_dir}"

# Install Cassandra
dcos package install --options=examples/dcos-minimal/pkg-cassandra.json cassandra --yes
ci/test-dcos-app-health.sh 'cassandra'

# Install Marathon-LB
dcos package install --options=examples/dcos-minimal/pkg-marathon-lb.json marathon-lb --yes
ci/test-dcos-app-health.sh 'marathon-lb'

# Install Tweeter
dcos marathon app add examples/dcos-minimal/app-tweeter.json
ci/test-dcos-app-health.sh 'tweeter'

# Test HTTP status
curl --fail --location --silent --show-error http://tweeter.acme.org/ -o /dev/null

# Test load balancing uses all instances
ci/test-dcos-tweeter-lb.sh

# Test posting and reading posts
ci/test-tweeting.sh

# Uninstall Tweeter
dcos marathon app remove tweeter

# Uninstall Marathon-LB
dcos package uninstall marathon-lb

# Uninstall Cassandra
dcos package uninstall cassandra

# Clean up after Cassandra
dcos marathon app add examples/dcos-minimal/app-janitor-cassandra.json
# App will exit automatically as TASK_KILLED or TASK_FAILED.
# There's no convenient way to wait for and verify completion without SSH access. :(
sleep 2

# TODO: poll mesos state for framework completion?
