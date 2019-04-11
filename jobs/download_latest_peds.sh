#!/bin/bash
# Download the latest data from https://ped.uspto.gov/peds/.

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"
START_DATE=`date +%Y%m%d`
SLACK=/usr/src/app/jobs/log_slack.sh
PEDS_BULK_URL="https://ped.uspto.gov/api/full-download?format=JSON"

${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS downloading started." \
-body="PEDS data is now started downloading. Will let you know when it's done (or failed). Date: ${START_DATE}" \
-recipient=${RECIPIENT}

$SLACK info "Start downloading latest data. Date: ${START_DATE}. Save to a temp file."
wget --tries=3 --output-document=${DATA_DIR}/raw.temp.zip ${PEDS_BULK_URL}

if [ $? -eq 0 ]; 
then
    $SLACK success "Download complete!"

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS is downloaded." -body="New bulk data is downloaded." -recipient=${RECIPIENT}

    # Rename the temp file to date file.
    mv ${DATA_DIR}/raw.temp.zip ${DATA_DIR}/raw.${START_DATE}.zip

    # Remove zip files older than 3 days.
    find ${DATA_DIR} -mtime +3 -name '*.zip' -delete
else
    $SLACK error "Downloading failed."

    # Remove incomplete temp file.
    rm -rf ${DATA_DIR}/raw.temp.zip

    ${APP_DIR}/bin/mail -subject="[PatHub Backend] PEDS data downloading is failed." -body="PEDS bulk data is NOT downloaded. Please check." -recipient=${RECIPIENT}
fi