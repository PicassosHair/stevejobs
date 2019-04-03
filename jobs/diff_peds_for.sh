#!/bin/bash
# This script pass $1 = year parameter and load differences.
# This job should not be ran with parse_peds_for at the same time.

# Variables.
YEAR=$1
DATA_DIR=/data
APP_DIR=/usr/src/app

START_TIME=`date +%s`
START_DATE=`date +%Y%m%d`
LOG_SLACK=${APP_DIR}/jobs/log_slack.sh

# Validating data files.
LATEST_RAW_ZIP=`ls ${DATA_DIR}/*.zip -t | head -n 1`
SECOND_LATEST_RAW_ZIP=`ls ${DATA_DIR}/*.zip -t | head -n 2 | tail -n 1`

${APP_DIR}/bin/slack chat send "Start: Diff PEDS data for year ${YEAR}" "#jobs"

if [ -e ${LATEST_RAW_ZIP} ] && [ -e ${SECOND_LATEST_RAW_ZIP} ]
then
  if [[ ${LATEST_RAW_ZIP} -ef ${SECOND_LATEST_RAW_ZIP} ]]
  then
    echo "Error. Data files are same. Stop."
    exit 1
  fi
  echo "Found raw data files -latest: ${LATEST_RAW_ZIP}, -older: ${SECOND_LATEST_RAW_ZIP}. Continue."
else
  echo "Error. Latest raw zip file doesn't exist. Stop parsing."
  exit 1
fi

# Prep work.
# Clear all temp/*.json files if exists.
rm -rf ${DATA_DIR}/temp
mkdir -p ${DATA_DIR}/temp

# Unzip old year file and rename to xxxx.old.json.
echo "Unzip data files."
unzip -o ${SECOND_LATEST_RAW_ZIP} ${YEAR}.json -d ${DATA_DIR}/temp/ ${YEAR}.json
mv ${DATA_DIR}/temp/${YEAR}.json ${DATA_DIR}/temp/${YEAR}.old.json

# Unzip new year.
unzip -o ${LATEST_RAW_ZIP} ${YEAR}.json -d ${DATA_DIR}/temp/ ${YEAR}.json
echo "Unzip done."

# Generate applications, codes, and transactions by line.
echo "Parsing old raw json."
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/${YEAR}.old.json -out=${DATA_DIR}/temp

echo "Renaming old output to xxx.old."
mv ${DATA_DIR}/temp/applications ${DATA_DIR}/temp/applications.old
mv ${DATA_DIR}/temp/codes ${DATA_DIR}/temp/codes.old
mv ${DATA_DIR}/temp/transactions ${DATA_DIR}/temp/transactions.old

echo "Parsing new raw json."
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/${YEAR}.json -out=${DATA_DIR}/temp

# Run diff to generate changed.
echo "Diffing applications."
diff --speed-large-files ${DATA_DIR}/temp/applications ${DATA_DIR}/temp/applications.old > ${DATA_DIR}/temp/applications.diff

echo "Diffing codes."
diff --speed-large-files ${DATA_DIR}/temp/codes ${DATA_DIR}/temp/codes.old > ${DATA_DIR}/temp/codes.diff

echo "Diffing transactions"
diff --speed-large-files ${DATA_DIR}/temp/transactions ${DATA_DIR}/temp/transactions.old > ${DATA_DIR}/temp/transactions.diff

# Run post diff to generate a SQL-ready list.
echo "Post diffing."
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/applications.diff -out=${DATA_DIR}/temp/applications.final
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/codes.diff -out=${DATA_DIR}/temp/codes.final
${APP_DIR}/bin/postdiff -in=${DATA_DIR}/temp/transactions.diff -out=${DATA_DIR}/temp/transactions.final

# Load to DB.
echo "Loading to DB..."
# Generate raw load application SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh application ${DATA_DIR}/temp/applications.final ${YEAR}

# Generate raw load codes SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh code ${DATA_DIR}/temp/codes.final ${YEAR}

# Generate raw load codes SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh transaction ${DATA_DIR}/temp/transactions.final ${YEAR}

echo "Done! Used $(expr `date +%s` - $START_TIME) s."