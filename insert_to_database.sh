#!/bin/bash
# This script read a pre-parsed csv-like file, and pump to the database.
# $1 - must be "application", "code", or "transaction".
# $2 - parsed file location.
# $3 - year. (e.g. 2017)

BASEDIR=/root/pedsparser

# Load color output module.
. ${BASEDIR}/_rainbow.sh

BASEDIR=$(pwd)
TABLE_NAME=$1
PARSED_FILE_PATH=$2
YEAR=$3

echogreen "Loading ${TABLE_NAME}s to the database..."
cat "${BASEDIR}/sql/load_${TABLE_NAME}s.sql" > ${BASEDIR}/temp/load_${TABLE_NAME}.sql
sed -i -e "s|@infile|'${PARSED_FILE_PATH}'|g; s|@year|${YEAR}|g;" ${BASEDIR}/temp/load_${TABLE_NAME}.sql

# Load sql to database
mysql --defaults-extra-file=$HOME/config/mysql.conf --local-infile -e \
"SOURCE ${BASEDIR}/temp/load_${TABLE_NAME}.sql;"

if [ $? -ne 0 ]; then
    echored "Loading to DB failed."
    exit 1
fi