#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"
START_DATE=`date +%Y%m%d`

${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS downloading started." \
-body="PEDS data is now started downloading. Will let you know when it's done (or failed). Date: ${START_DATE}" \
-recipient=${RECIPIENT}

echo "Start downloading latest data."
wget --tries=3 --output-document=${DATA_DIR}/raw.${START_DATE}.zip https://ped.uspto.gov/api/full-download\?format\=JSON

if [ $? -eq 0 ]; 
then
    echo "Download complete!"

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient=${RECIPIENT}

    # Remove oldest file keep total files count 3.
    ls ${DATA_DIR}/*.zip -1t | tail -n +4 | xargs rm -f
else
    echo "Downloading failed."

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS data downloading is failed." -body="PEDS bulk data is NOT downloaded. Please check." -recipient=${RECIPIENT}
fi