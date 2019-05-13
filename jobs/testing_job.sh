#!/bin/bash

APP_DIR=/usr/src/app

echo "This is a testing job that does not do anything"
${APP_DIR}/jobs/log_slack.sh info "Test okay msg."

echo "Testing latest_file_name"
LATEST_FILE_NAME=$(${APP_DIR}/bin/latest_file_name)
echo "${LATEST_FILE_NAME}"