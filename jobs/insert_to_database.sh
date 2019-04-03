#!/bin/bash
# This script read a pre-parsed csv-like file, and pump to the database.
# $1 - Enum. Must be one of "application", "code", or "transaction".
# $2 - Parsed file location.
# $3 - Year. (e.g. 2017)

DATA_DIR=/data
APP_DIR=/usr/src/app
START_TIME=`date +%s`
TABLE_NAME=$1
PARSED_FILE_PATH=$2
YEAR=$3

echo "Loading ${TABLE_NAME}s to the database."
${APP_DIR}/bin/slack chat send "Start: Load data into database for ${TABLE_NAME}, year ${YEAR}" "#jobs"

# Copy SQL template file to a temp location.
cat "${APP_DIR}/sql/load_${TABLE_NAME}s.sql" > ${DATA_DIR}/temp/load_${TABLE_NAME}.sql
# Replace some keywords in the template file and create a pumping script.
# In this case, replace "@infile" with the csv-like file, and replace "@year" with YEAR.
sed -i -e "s|@infile|'${PARSED_FILE_PATH}'|g; s|@year|${YEAR}|g;" ${DATA_DIR}/temp/load_${TABLE_NAME}.sql

# Load sql to database. Use file to locate auth as no safe way to put the password in cmd.
mysql --defaults-extra-file=${APP_DIR}/mysql.conf \
--local-infile \
-e "SOURCE ${DATA_DIR}/temp/load_${TABLE_NAME}.sql;"

if [ $? -ne 0 ]; then
    echo "Failed to pump data into the database."
    exit 1
else
    echo "Successfully pump into database for ${TABLE_NAME} on ${YEAR}. Used $(expr `date +%s` - $START_TIME) s."
    ${APP_DIR}/bin/slack chat send "Success: Pump into database for ${TABLE_NAME} on ${YEAR}. Used $(expr `date +%s` - $START_TIME) s." "#jobs"
fi