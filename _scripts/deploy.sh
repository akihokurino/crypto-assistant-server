#!/bin/bash

APP_ROOT=$(dirname $0)/..

gcloud app deploy $APP_ROOT/entrypoint/api/cron.yaml
gcloud app deploy $APP_ROOT/entrypoint/api/index.yaml
appcfg.py update --version 1 --application crypto-assistant-dev $APP_ROOT/entrypoint/api/app.yaml $APP_ROOT/entrypoint/batch/app.yaml