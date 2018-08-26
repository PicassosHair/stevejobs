-- This template script will read codes file generated by the go package and load to mysql database

CREATE TEMPORARY TABLE IF NOT EXISTS temp_Codes LIKE TransactionCodes;
TRUNCATE TABLE temp_Codes;

LOAD DATA LOCAL INFILE @infile
INTO TABLE temp_Codes
FIELDS TERMINATED BY '^'
LINES TERMINATED BY '$\n'
(createdAt, updatedAt, description, code, type, initiator, isActionable);

INSERT IGNORE INTO TransactionCodes (createdAt, updatedAt, code, description, type, initiator, isActionable)
SELECT createdAt, updatedAt, code, description, type, initiator, isActionable FROM temp_Codes;

DROP TEMPORARY TABLE temp_Codes;