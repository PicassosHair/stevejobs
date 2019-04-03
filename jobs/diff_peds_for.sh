#!/bin/bash
# This script pass $1 = year parameter and load differences.
# This job should not be ran with parse_peds_for at the same time.

# Variables.
YEAR=$1
DATA_DIR=/data
APP_DIR=/usr/src/app
SLACK=/usr/src/app/jobs/log_slack.sh

START_TIME=`date +%s`
START_DATE=`date +%Y%m%d`

# Validating data files.
LATEST_RAW_ZIP=`ls ${DATA_DIR}/*.zip -t | head -n 1`
SECOND_LATEST_RAW_ZIP=`ls ${DATA_DIR}/*.zip -t | head -n 2 | tail -n 1`

$SLACK info "Start diffing PEDS data for year ${YEAR}."

if [ -e ${LATEST_RAW_ZIP} ] && [ -e ${SECOND_LATEST_RAW_ZIP} ]
then
  if [[ ${LATEST_RAW_ZIP} -ef ${SECOND_LATEST_RAW_ZIP} ]]
  then
    $SLACK error "Error. Data files are same. Stop."
    exit 1
  fi
  $SLACK info "Found raw data files. Latest: ${LATEST_RAW_ZIP}, previous: ${SECOND_LATEST_RAW_ZIP}. Continue."
else
  $SLACK error "Error. Latest raw zip file doesn't exist. Stop parsing."
  exit 1
fi

# Prep work.
# Clear all temp/*.json files if exists.
rm -rf ${DATA_DIR}/temp
mkdir -p ${DATA_DIR}/temp

# Unzip old year file and rename to xxxx.old.json.
unzip -o ${SECOND_LATEST_RAW_ZIP} ${YEAR}.json -d ${DATA_DIR}/temp/ ${YEAR}.json
mv ${DATA_DIR}/temp/${YEAR}.json ${DATA_DIR}/temp/${YEAR}.old.json

# Unzip new year.
unzip -o ${LATEST_RAW_ZIP} ${YEAR}.json -d ${DATA_DIR}/temp/ ${YEAR}.json
$SLACK info "Unzip done. Latest: ${LATEST_RAW_ZIP}, previous: ${SECOND_LATEST_RAW_ZIP}."

# Generate applications, codes, and transactions by line.
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/${YEAR}.old.json -out=${DATA_DIR}/temp

# Renaming old output to xxx.old.
mv ${DATA_DIR}/temp/applications ${DATA_DIR}/temp/applications.old
mv ${DATA_DIR}/temp/codes ${DATA_DIR}/temp/codes.old
mv ${DATA_DIR}/temp/transactions ${DATA_DIR}/temp/transactions.old

# Parsing new raw json.
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/${YEAR}.json -out=${DATA_DIR}/temp

# Run diff to generate changed.
$SLACK info "Diffing applications for year ${YEAR}."
diff --speed-large-files ${DATA_DIR}/temp/applications ${DATA_DIR}/temp/applications.old > ${DATA_DIR}/temp/applications.diff

$SLACK info "Diffing codes for year ${YEAR}."
diff --speed-large-files ${DATA_DIR}/temp/codes ${DATA_DIR}/temp/codes.old > ${DATA_DIR}/temp/codes.diff

$SLACK info "Diffing transactions for year ${YEAR}."
diff --speed-large-files ${DATA_DIR}/temp/transactions ${DATA_DIR}/temp/transactions.old > ${DATA_DIR}/temp/transactions.diff

# Run post diff to generate a SQL-ready list.
$SLACK info "Generate *.final for diffs."
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/applications.diff -out=${DATA_DIR}/temp/applications.final
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/codes.diff -out=${DATA_DIR}/temp/codes.final
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/transactions.diff -out=${DATA_DIR}/temp/transactions.final

# Load to DB.
$SLACK info "Loading diffs into database."
# Generate raw load application SQL file.
${APP_DIR}/jobs/insert_to_database.sh application ${DATA_DIR}/temp/applications.final ${YEAR}

# Generate raw load codes SQL file.
${APP_DIR}/jobs/insert_to_database.sh code ${DATA_DIR}/temp/codes.final ${YEAR}

# Generate raw load codes SQL file.
${APP_DIR}/jobs/insert_to_database.sh transaction ${DATA_DIR}/temp/transactions.final ${YEAR}

$SLACK success "Done! Used $(expr `date +%s` - $START_TIME) s."