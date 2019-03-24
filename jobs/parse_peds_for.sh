#!/bin/bash
# This script pass $1 = year parameter and load the data of $1.json to mysql database.

# Stop if error.
set -e

DATA_DIR=/data
APP_DIR=/usr/src/app
RECIPIENT="liuhao1990@gmail.com,hinmeng@gmail.com"

# Put MySQL credentials into $HOME/config/mysql.conf
YEAR=$1
START_TIME=`date +%s`

# Prepare work. Remove temp folder if exists.
rm -rf ${DATA_DIR}/temp
mkdir -p ${DATA_DIR}/temp

# Unzip raw.json to temp/YYYY.json
echo "Unzipping raw.zip. for year ${YEAR}"
unzip -o ${DATA_DIR}/raw.zip $1.json -d ${DATA_DIR}/temp/

if [ $? -ne 0 ]; then
    echo "Unzip failed."
    exit 1
fi

# Parse json file to temp/applications, temp/codes, temp/transactions
echo "Parsing ${YEAR}.json."
${APP_DIR}/bin/parser -in=${DATA_DIR}/temp/$YEAR.json -out=${DATA_DIR}/temp

if [ $? -ne 0 ]; then
    echo "Parse failed."
    exit 1
fi

# Generate raw load application SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh application ${DATA_DIR}/temp/applications ${YEAR}

# Generate raw load codes SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh code ${DATA_DIR}/temp/codes ${YEAR}

# Generate raw load codes SQL file.
bash ${APP_DIR}/jobs/insert_to_database.sh transaction ${DATA_DIR}/temp/transactions ${YEAR}

# Clean up work.
rm -rf ${DATA_DIR}/temp

echo "Done! Used $(expr `date +%s` - $START_TIME) s"