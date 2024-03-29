#!/bin/bash
# Send a message to Slack channel #jobs.
# Args:
#   $1: Message type, one of "info", "success", "warning", "error".
#   $2: Message content.

APP_DIR=/usr/src/app
INFO_COLOR="#3498DB"
SUCCESS_COLOR="#28B463"
WARNING_COLOR="#F39C12"
ERROR_COLOR="#C70039"
COLOR=${INFO_COLOR}

if [ "$1" == "success" ]; then
  COLOR=${SUCCESS_COLOR}
fi
if [ "$1" == "warning" ]; then
  COLOR=${WARNING_COLOR}
fi
if [ "$1" == "error" ]; then
  COLOR=${ERROR_COLOR}
fi

LOG_PREFIX="[LOG SLACK]"

if [ "${ENV}" == "development" ]; then
  LOG_PREFIX="[DEV]${LOG_PREFIX}"
fi

echo "${LOG_PREFIX} ${2}"
${APP_DIR}/bin/slack chat send \
  --text "${2}" \
  --channel "#jobs" \
  --author "stevejobs-bot" \
  --author-icon "https://github.com/fuermosi777/stevejobs/raw/master/bot.png" \
  --color ${COLOR} \
  &>/dev/null
