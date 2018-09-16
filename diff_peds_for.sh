#!/bin/bash
# This script pass $1 = year parameter and load differences.

# Load color echo file.
. ./_rainbow.sh

# Variables.
YEAR=$1
BASEDIR=$(pwd)
START_TIME=`date +%s`

# Prep work.
# Clear all temp/*.json files if exists.
# Prepare work.
rm -rf ./temp
mkdir -p ./temp

# Unzip old year file and rename to xxxx.old.json
unzip -o ${BASEDIR}/data/raw.old.zip ${YEAR}.json -d ${BASEDIR}/temp/ ${YEAR}.json
mv ./temp/${YEAR}.json ./temp/${YEAR}.old.json

# Unzip new year
unzip -o ${BASEDIR}/data/raw.zip ${YEAR}.json -d ${BASEDIR}/temp/ ${YEAR}.json
echogreen "Unzip done."

# Generate applications, codes, and transactions by line.
echogreen "Parsing old raw json."
${BASEDIR}/bin/parser -in=${BASEDIR}/temp/${YEAR}.old.json -out=${BASEDIR}/temp

echogreen "Renaming old output to xxx.old."
mv ./temp/applications ./temp/applications.old
mv ./temp/codes ./temp/codes.old
mv ./temp/transactions ./temp/transactions.old

echogreen "Parsing new raw json."
${BASEDIR}/bin/parser -in=${BASEDIR}/temp/${YEAR}.json -out=${BASEDIR}/temp

# Run diff to generate changed.
echogreen "Diffing applications..."
diff --speed-large-files ${BASEDIR}/temp/applications ${BASEDIR}/temp/applications.old > ${BASEDIR}/temp/applications.diff

echogreen "Diffing codes..."
diff --speed-large-files ${BASEDIR}/temp/codes ${BASEDIR}/temp/codes.old > ${BASEDIR}/temp/codes.diff

echogreen "Diffing transactions.."
diff --speed-large-files ${BASEDIR}/temp/transactions ${BASEDIR}/temp/transactions.old > ${BASEDIR}/temp/transactions.diff

# Run post diff to generate a SQL-ready list.
echogreen "Post diffing..."
${BASEDIR}/bin/postdiff -in=${BASEDIR}/temp/applications.diff -out=${BASEDIR}/temp/applications.final
${BASEDIR}/bin/postdiff -in=${BASEDIR}/temp/codes.diff -out=${BASEDIR}/temp/codes.final
${BASEDIR}/bin/postdiff -in=${BASEDIR}/temp/transactions.diff -out=${BASEDIR}/temp/transactions.final

# Load to DB.
echogreen "Loading to DB..."
# Generate raw load application SQL file.
./insert_to_database.sh application ${BASEDIR}/temp/applications.final ${YEAR}

# Generate raw load codes SQL file.
./insert_to_database.sh code ${BASEDIR}/temp/codes.final ${YEAR}

# Generate raw load codes SQL file.
./insert_to_database.sh transaction ${BASEDIR}/temp/transactions.final ${YEAR}

# Cleanup
rm -rf ./temp

echogreen "Done! Used $(expr `date +%s` - $START_TIME) s"