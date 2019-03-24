#!/bin/bash
# This script read a pre-parsed csv-like file, and pump to the database.
# $1 - must be "application", "code", or "transaction".
# $2 - parsed file location.
# $3 - year. (e.g. 2017)

DATA_DIR=/data
APP_DIR=/usr/src/app

TABLE_NAME=$1
PARSED_FILE_PATH=$2
YEAR=$3

echo "Loading ${TABLE_NAME}s to the database."
# Copy SQL template file to a temp location.
cat "${APP_DIR}/sql/load_${TABLE_NAME}s.sql" > ${DATA_DIR}/temp/load_${TABLE_NAME}.sql
# Replace some keywords in the template file and create a pumping script.
sed -i -e "s|@infile|'${PARSED_FILE_PATH}'|g; s|@year|${YEAR}|g;" ${DATA_DIR}/temp/load_${TABLE_NAME}.sql

# Load sql to database. Use file to locate auth as no safe way to put the password in cmd.
mysql --defaults-extra-file=${APP_DIR}/mysql.conf \
--enforce-gtid-consistency \
--local-infile \
-e "SOURCE ${DATA_DIR}/temp/load_${TABLE_NAME}.sql;"

if [ $? -ne 0 ]; then
    echo "Loading to DB failed."
    exit 1
fi