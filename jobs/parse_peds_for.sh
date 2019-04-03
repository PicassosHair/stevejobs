#!/bin/bash
# This script pass $1 = year parameter and load the data of $1.json to mysql database.

# Stop if error.
set -e

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"
YEAR=$1
START_TIME=`date +%s`
START_DATE=`date +%Y%m%d`

# Prepare work. Remove temp folder if exists.
rm -rf ${DATA_DIR}/temp
mkdir -p ${DATA_DIR}/temp

${APP_DIR}/jobs/log_slack.sh info "Parse PEDS data for year ${YEAR}."

# Check for latest raw.YYYYMMDD.zip file existance.
LATEST_RAW_ZIP=`ls ${DATA_DIR}/*.zip -t | head -n 1`

if [ -e ${LATEST_RAW_ZIP} ] 
then
  ${APP_DIR}/jobs/log_slack.sh info "Found raw data: ${LATEST_RAW_ZIP}, continue."
else
  ${APP_DIR}/jobs/log_slack.sh error "Error. Latest raw zip file doesn't exist. Stop parsing for year ${YEAR}."
  exit 1
fi


# Unzip raw.json to temp/YYYY.json
${APP_DIR}/jobs/log_slack.sh info "Unzipping raw.${START_DATE}.zip. for year ${YEAR}."
unzip -o ${LATEST_RAW_ZIP} $1.json -d ${DATA_DIR}/temp/

if [ $? -ne 0 ]; then
    ${APP_DIR}/jobs/log_slack.sh error "Unzip failed for year ${YEAR}."
    exit 1
fi

# Parse json file to temp/applications, temp/codes, temp/transactions, which are csv-like files.
${APP_DIR}/jobs/log_slack.sh info "Parsing ${YEAR}.json."
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/$YEAR.json -out=${DATA_DIR}/temp

if [ $? -ne 0 ]; then
    ${APP_DIR}/jobs/log_slack.sh error "Parsing ${YEAR}.json failed."
    exit 1
fi

# Generate raw load application SQL file.
${APP_DIR}/jobs/insert_to_database.sh application ${DATA_DIR}/temp/applications ${YEAR}

# Generate raw load codes SQL file.
${APP_DIR}/jobs/insert_to_database.sh code ${DATA_DIR}/temp/codes ${YEAR}

# Generate raw load codes SQL file.
${APP_DIR}/jobs/insert_to_database.sh transaction ${DATA_DIR}/temp/transactions ${YEAR}

${APP_DIR}/jobs/log_slack.sh success "Done parsing data for year ${YEAR}! Used $(expr `date +%s` - $START_TIME) s."