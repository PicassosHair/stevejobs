#!/bin/bash
# This script accept two years and parse the PEDS data between them.
# $1 - starting year (e.g. 2013)
# $2 - ending year (e.g. 2018)

APP_DIR=/usr/src/app
SLACK=/usr/src/app/jobs/log_slack.sh

$SLACK info "Start parsing PEDS data between year ${1} and ${2}."

for (( y=$1; y<=$2; y++ ))
    do 
    ${APP_DIR}/jobs/parse_peds_for.sh $y
done

$SLACK success "Done parsing PEDS data between year ${1} and ${2}."