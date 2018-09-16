#!/bin/bash
# This script pass $1 = year parameter and load the data of $1.json to mysql database.

# Stop if error.
set -e

# Load color echo file.
. ./_rainbow.sh

# Put MySQL credentials into $HOME/config/mysql.conf
BASEDIR=$(pwd)
YEAR=$1
START_TIME=`date +%s`

# Prepare work.
rm -rf ./temp
mkdir -p ./temp

# Unzip data/raw.json to temp/YYYY.json
echogreen "Unzipping raw.zip..."
unzip -o ${BASEDIR}/data/raw.zip $1.json -d ${BASEDIR}/temp/

if [ $? -ne 0 ]; then
    echored "Unzip failed."
    exit 1
fi

# Parse json file to temp/applications, temp/codes, temp/transactions
echogreen "Parsing ${YEAR}.json..."
${BASEDIR}/bin/parser -in=${BASEDIR}/temp/$YEAR.json -out=${BASEDIR}/temp

if [ $? -ne 0 ]; then
    echored "Parse failed."
    exit 1
fi

# Generate raw load application SQL file.
${BASEDIR}/insert_to_database.sh application ${BASEDIR}/temp/applications ${YEAR}

# Generate raw load codes SQL file.
${BASEDIR}/insert_to_database.sh code ${BASEDIR}/temp/codes ${YEAR}

# Generate raw load codes SQL file.
${BASEDIR}/insert_to_database.sh transaction ${BASEDIR}/temp/transactions ${YEAR}

# Clean work
rm -rf ./temp

echogreen "Done! Used $(expr `date +%s` - $START_TIME) s"