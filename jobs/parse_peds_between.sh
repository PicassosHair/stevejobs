#!/bin/bash
# This script accept two years and parse the PEDS data between them.
# $1 - starting year (e.g. 2013)
# $2 - ending year (e.g. 2018)

APP_DIR=/usr/src/app

echo "Start parsing PEDS data between $1 ... $2";

for (( y=$1; y<=$2; y++ ))
    do 
    echo "Start parsing year for $y."
    bash ${APP_DIR}/jobs/parse_peds_for.sh $y
done