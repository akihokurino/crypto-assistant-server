#!/bin/bash
set -eu

APP_ROOT=$(dirname $0)/..

SERVE_OPTS=${SERVE_OPTS:-"--storage_path=$APP_ROOT/.var/dev_appserver"}
SERVE_TARGETS=${SERVE_TARGETS:-"api batch"}
dev_appserver.py $SERVE_OPTS $(echo $SERVE_TARGETS | sed -E 's@([^ ]+)@entrypoint/\1@g')
