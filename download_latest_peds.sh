#!/bin/bash

# Load color output module.
. ./_rainbow.sh

# Remove raw.old.zip if exists.
rm -f ./data/raw.old.zip
echogreen "Removed old raw.zip file."

# Rename raw.zip to raw.old.zip
if [ -e ./data/raw.zip ]
then
    mv ./data/raw.zip ./data/raw.old.zip
    echogreen "Renamed raw.zip to raw.old.zip"
else
    echored "Not found raw.zip file."
    exit 1
fi

echogreen "Start download latest data."
wget https://ped.uspto.gov/api/full-download\?format\=JSON --output-document=./data/raw.zip --show-progress

if [ $? -eq 0 ]; then
    echogreen "Download complete!"
else
    rm -rf ./data/raw.zip
    mv ./data/raw.old.zip ./data/raw.zip
    echored "Download failed."
fi