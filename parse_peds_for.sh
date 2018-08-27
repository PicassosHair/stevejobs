#!/bin/bash
# This script pass $1 = year parameter and load the data of $1.json to mysql database.

# Put MySQL credentials into $HOME/config/mysql.conf

BASEDIR=$(pwd)
LOG_FILE_PATH="./result.log"

# Generate a SQL file, then load the file to mysql.
function load_to_db ()
{
    echo "[INFO] Loading ${1}s to the database..."
    cat "${BASEDIR}/sql/load_${1}s.sql" > ./temp/load_${1}.sql
    sed -i -e "s|@infile|'${BASEDIR}/temp/${1}s'|g" ./temp/load_${1}.sql

    # Load sql to database
    mysql --defaults-extra-file=$HOME/config/mysql.conf --local-infile -e \
    "USE ${IDSGUARD_DB_NAME};\
    SOURCE ./temp/load_${1}.sql;" > ${LOG_FILE_PATH}
}

# Prepare work.
mkdir -p ./temp

# Unzip data/raw.json to temp/YYYY.json
echo "[INFO] Unzipping raw.zip..."
unzip -o ${BASEDIR}/data/raw.zip $1.json -d ${BASEDIR}/temp/

# Parse json file to temp/applications, temp/codes, temp/transactions
echo "[INFO] Parsing ${1}.json..."
${BASEDIR}/src/parser -in=${BASEDIR}/temp/$1.json -out=${BASEDIR}/temp

# Generate raw load application SQL file.
load_to_db "application"

# Generate raw load codes SQL file.
load_to_db "code"

# Generate raw load codes SQL file.
load_to_db "transaction"

# Clean work
rm -rf ./temp

echo "[INFO] Done!"