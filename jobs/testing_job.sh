#!/bin/bash

S="${APP_DIR}/jobs/log_slack.sh"

echo "This is a testing job that does not do anything"
${APP_DIR}/jobs/log_slack.sh info "Test okay msg."

$S info "Okay."