#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"

# Remove raw.old.zip if exists.
rm -f ${DATA_DIR}/raw.old.zip
echo "Removed old raw.zip file."

# Rename raw.zip to raw.old.zip if needed, or just download the data.
if [ -e ${DATA_DIR}/raw.zip ]
then
    mv ${DATA_DIR}/raw.zip ${DATA_DIR}/raw.old.zip
    echo "Renamed raw.zip to raw.old.zip."
else
    echo "Not found old raw.zip file."
    touch ${DATA_DIR}/raw.zip
fi

echo "Start downloading latest data."
wget --output-document=${DATA_DIR}/raw.zip --show-progress https://ped.uspto.gov/api/full-download\?format\=JSON

if [ $? -eq 0 ]; then
    echo "Download complete!"

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient=${RECIPIENT}
else
    rm -rf ${DATA_DIR}/raw.zip
    mv ${DATA_DIR}/raw.old.zip ${DATA_DIR}/raw.zip
    echo "Downloading failed. Rolled back everything."

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS data downloading is failed." -body="New bulk data is NOT downloaded. Please check." -recipient=${RECIPIENT}
fi