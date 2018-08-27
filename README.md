# PEDS Parser

Parse PEDS raw JSON file into organized file.

# Need

wget
mysql
go


# Usage

1. Download PEDS data from PEDS website: https://ped.uspto.gov/peds/#/apiDocumentation

`./download_latest_peds.sh`

2. Parse year range:

`./parse_peds_between 2000 2018`

# Data flow

PEDS webiste (https://ped.uspto.gov/peds/) ->
*download_peds.sh* -> ./data/raw.zip
*parse_peds.sh* -> ./temp/2018.json
(processing, injecting ^) -> ./temp/applications, ./temp/codes, ./temp/transactions
-> mysql