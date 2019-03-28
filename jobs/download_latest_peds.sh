#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"

${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS downloading started." \
-body="PEDS data is now started downloading. Will let you know when it's done (or failed)." \
-recipient=${RECIPIENT}

# Remove raw.old.zip if exists.
if [ -e ${DATA_DIR}/raw.old.zip ]
  rm -f ${DATA_DIR}/raw.old.zip
  echo "Removed old raw.zip file."
fi

# Rename raw.zip to raw.old.zip if needed, or just download the data.
if [ -e ${DATA_DIR}/raw.zip ] then
  echo "Renamed raw.zip to raw.old.zip."
  mv ${DATA_DIR}/raw.zip ${DATA_DIR}/raw.old.zip
else
  echo "Not found old raw.zip file. Create a fake raw.zip for placeholder."
  touch ${DATA_DIR}/raw.zip
fi

echo "Start downloading latest data."
wget --output-document=${DATA_DIR}/raw.zip https://ped.uspto.gov/api/full-download\?format\=JSON

if [ $? -eq 0 ]; then
    echo "Download complete!"

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient=${RECIPIENT}
else
    echo "Downloading failed."

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS data downloading is failed." -body="PEDS bulk data is NOT downloaded. Please check." -recipient=${RECIPIENT}
fi