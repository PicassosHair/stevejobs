#!/bin/bash

. ./_rainbow.sh

wget https://ped.uspto.gov/api/full-download\?format\=JSON --output-document=./data/raw.zip --show-progress

if [ $? -eq 0 ]; then
    echogreen "Download complete!"
else
    echored "Download failed."
fi