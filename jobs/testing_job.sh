#!/bin/bash

APP_DIR=/usr/src/app
S="${APP_DIR}/jobs/log_slack.sh"
SLAC=/usr/src/app/jobs/log_slack.sh

echo "This is a testing job that does not do anything"
${APP_DIR}/jobs/log_slack.sh info "Test okay msg."

$S info "Okay."
$SLAC error "Erro."