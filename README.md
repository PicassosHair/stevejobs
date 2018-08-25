# PEDS Parser

Parse PEDS raw JSON file into organized file.

# Usage

./parse_peds.sh <application_out_file_path>

# Data flow

PEDS webiste (https://ped.uspto.gov/peds/) ->
*download_peds.sh* -> ./data/raw.zip
*parse_peds.sh* -> ./temp/2018.json
(processing, injecting ^) -> ./temp/applications, ./temp/codes, ./temp/transactions
-> mysql