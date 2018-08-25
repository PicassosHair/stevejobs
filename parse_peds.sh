#!/bin/sh
LOAD_APPLICATION_SCRIPT_PATH="./sql/load_applications.sql"
LOG_FILE_PATH="./result.log"

# Prepare work.
mkdir -p ./temp

# Generate raw load application SQL file.
cat ${LOAD_APPLICATION_SCRIPT_PATH} > ./temp/load_application.sql
sed -i -e "s|@infile|'${1}'|g" ./temp/load_application.sql

# Load applications to database
mysql --defaults-extra-file=$HOME/config/mysql.conf --local-infile -e \
"USE ${IDSGUARD_DB_NAME};\
SOURCE ./temp/load_application.sql;" > ${LOG_FILE_PATH}

# Clean work
rm -rf ./temp