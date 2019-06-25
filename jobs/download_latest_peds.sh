#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/. There are two options:
#   $1 - isFull: true or false

DATA_DIR=/data
APP_DIR=/usr/src/app
START_DATE=$(date +%Y%m%d)
SLACK=/usr/src/app/jobs/log_slack.sh

# Log delta or full.
if [[ $# -eq 0 ]]; then
    $SLACK error "Missing isFull option. Enter true or false."
    exit 1
elif [ "$1" = "true" ]
then
    $SLACK info "Download full dataset."
elif [ "$1" = "false" ]
then
    $SLACK info "Download delta dataset."
else
    $SLACK error "Invalid isFull option. Stop."
    exit 1
fi

LATEST_FILE_NAME=$(${APP_DIR}/bin/latest_file_name --isFull=$1)

PEDS_BULK_URL="https://ped.uspto.gov/api/full-download?fileName=${LATEST_FILE_NAME}"

$SLACK info "Start downloading latest data. Start date: ${START_DATE}. Downloding from ${PEDS_BULK_URL}. Saving to a temp file."
wget --tries=3 --output-document=${DATA_DIR}/raw.temp.zip ${PEDS_BULK_URL}

if [ $? -eq 0 ]; then
    $SLACK success "Download complete. Saved file to ${DATA_DIR}/raw.${START_DATE}.zip"

    # Rename the temp file to date file.
    mv ${DATA_DIR}/raw.temp.zip ${DATA_DIR}/raw.${START_DATE}.zip

    # Touch the new file to make sure its mtime is up to date.
    touch -m ${DATA_DIR}/raw.${START_DATE}.zip

    # Remove zip files older than 3 days.
    find ${DATA_DIR} -mtime +3 -name '*.zip' -delete
else
    $SLACK error "Downloading failed."

    # Remove incomplete temp file.
    rm -rf ${DATA_DIR}/raw.temp.zip

fi
