# PEDS Parser

You should do this first in your local environment to get the sample data.
Parse PEDS raw JSON file into organized file. Using Go.

# Prepare

Install wget if don't have already. Install MySQL.
Install Go. Put this in .bashrc or .zshrc: `export GOPATH=$HOME/workspace/pedsparser`

# Local Testing Data

Unzip 2018.json from compressed data file, and **move to ./temp folder**.

Parse: `$ ./bin/parser -in=./temp/2018.json -out=./temp`. This will generate three files in the `./temp` folder, taking about 1~5 minutes.

Execute these lines, inserting into local MySQL:
- `$ ./insert_to_database.sh application ./temp/applications 2018`
- `$ ./insert_to_database.sh code ./temp/codes 2018`
- `$ ./insert_to_database.sh transaction ./temp/transactions 2018` -> This will take a long time

Cleanup:
`$ rm -rf ./temp`

Check:

Use Sequel Pro to make sure you loaded all the data.

# Usage

- Download PEDS data from PEDS website: https://ped.uspto.gov/peds/#/apiDocumentation

`./download_latest_peds.sh`

- Parse year range:

`./parse_peds_between 2000 2018`

# Data flow [Outdated]

PEDS webiste (https://ped.uspto.gov/peds/) ->
*download_peds.sh* -> ./data/raw.zip
*parse_peds.sh* -> ./temp/2018.json
(processing, injecting ^) -> ./temp/applications, ./temp/codes, ./temp/transactions
-> mysql