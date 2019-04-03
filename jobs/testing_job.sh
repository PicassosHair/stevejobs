#!/bin/bash

LOG_SLACK=${APP_DIR}/jobs/log_slack.sh

echo "This is a testing job that does not do anything"
LOG_SLACK info "Test okay msg."