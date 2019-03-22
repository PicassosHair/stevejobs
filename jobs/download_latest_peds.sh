#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=${STORAGE_DIR}/peds

# Remove raw.old.zip if exists.
rm -f DATA_DIR/raw.old.zip
echo "Removed old raw.zip file."

# Rename raw.zip to raw.old.zip if needed, or just download the data.
if [ -e ${STORAGE_DIR}/raw.zip ]
then
    mv ${STORAGE_DIR}/raw.zip ${STORAGE_DIR}/raw.old.zip
    echo "Renamed raw.zip to raw.old.zip."
else
    echo "Not found old raw.zip file."
fi

echo "Start downloading latest data."
wget https://ped.uspto.gov/api/full-download\?format\=JSON --output-document=${STORAGE_DIR}/raw.zip --show-progress

if [ $? -eq 0 ]; then
    echo "Download complete!"

    ${WORK_DIR}/bin/mail -sender="mailman@pathub.io" -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient="liuhao1990@gmai.com,hinmeng@gmail.com"
else
    rm -rf ${STORAGE_DIR}/raw.zip
    mv ${STORAGE_DIR}/raw.old.zip ${STORAGE_DIR}/raw.zip
    echo "Downloading failed. Rolled back everything."

    ${WORK_DIR}/bin/mail -sender="mailman@pathub.io" -subject="[PatHub Backend] PEDS data downloading is failed." -body="New bulk data is NOT downloaded. Please check." -recipient="liuhao1990@gmai.com,hinmeng@gmail.com"
fi