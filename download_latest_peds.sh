#!/bin/bash

BASEDIR=/root/pedsparser

# Load color output module.
. ${BASEDIR}/_rainbow.sh

# Remove raw.old.zip if exists.
rm -f ${BASEDIR}/data/raw.old.zip
echogreen "Removed old raw.zip file."

# Rename raw.zip to raw.old.zip
if [ -e ${BASEDIR}/data/raw.zip ]
then
    mv ${BASEDIR}/data/raw.zip ${BASEDIR}/data/raw.old.zip
    echogreen "Renamed raw.zip to raw.old.zip"
else
    echoyellow "Not found raw.zip file."
    echogreen "Downloading a new raw.zip file..."

    wget https://ped.uspto.gov/api/full-download\?format\=JSON --output-document=${BASEDIR}/data/raw.zip --show-progress
    if [ $? -eq 0 ]; then
        echogreen "Download complete!"
    else
        echored "Download failed."
    fi

    exit 1
fi

echogreen "Start download latest data."
wget https://ped.uspto.gov/api/full-download\?format\=JSON --output-document=${BASEDIR}/data/raw.zip --show-progress

if [ $? -eq 0 ]; then
    echogreen "Download complete!"
else
    rm -rf ${BASEDIR}/data/raw.zip
    mv ${BASEDIR}/data/raw.old.zip ${BASEDIR}/data/raw.zip
    echored "Download failed."
fi