#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=/data/peds

# Remove raw.old.zip if exists.
rm -f DATA_DIR/raw.old.zip
echo "Removed old raw.zip file."

# Rename raw.zip to raw.old.zip if needed, or just download the data.
if [ -e /data/raw.zip ]
then
    mv /data/raw.zip /data/raw.old.zip
    echo "Renamed raw.zip to raw.old.zip."
else
    echo "Not found old raw.zip file."
    touch /data/raw.zip
fi

echo "Start downloading latest data."
wget --output-document=/data/raw.zip --show-progress https://ped.uspto.gov/api/full-download\?format\=JSON

if [ $? -eq 0 ]; then
    echo "Download complete!"

    /usr/src/app/bin/mail -sender="mailman@pathub.io" -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient="liuhao1990@gmail.com,hinmeng@gmail.com"
else
    rm -rf /data/raw.zip
    mv /data/raw.old.zip /data/raw.zip
    echo "Downloading failed. Rolled back everything."

    # /usr/src/app/bin/mail -sender="mailman@pathub.io" -subject="[PatHub Backend] PEDS data downloading is failed." -body="New bulk data is NOT downloaded. Please check." -recipient="liuhao1990@gmail.com,hinmeng@gmail.com"
fi