#!/bin/bash
# This script accept two years and diff them.
# $1 - starting year (e.g. 2013)
# $2 - ending year (e.g. 2018)

APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"
SLACK=/usr/src/app/jobs/log_slack.sh

$SLACK info "Start diff PEDS data between $1 and $2."

for (( y=$1; y<=$2; y++ ))
    do 
    ${APP_DIR}/jobs/diff_peds_for.sh $y
done

${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS parsing is done." -body="No error found." -recipient=${RECIPIENT}
$SLACK success "PEDS diffing between $1 and $2 is done."