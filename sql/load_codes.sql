-- This template script will read codes file generated by the go package and load to mysql database
DROP TABLE IF EXISTS temp_Codes;
-- Codes table like temp table without constraints.
CREATE TABLE IF NOT EXISTS `temp_Codes` (
  `code` VARCHAR(100) NOT NULL,
  `description` VARCHAR(255) DEFAULT NULL,
  `type` enum('info','warning','success') DEFAULT NULL,
  `initiator` enum('uspto','applicant') DEFAULT NULL,
  `isActionable` tinyint(1) NOT NULL DEFAULT '1'
) ENGINE=InnoDB;

LOAD DATA LOCAL INFILE @infile
INTO TABLE temp_Codes
FIELDS TERMINATED BY '^^'
LINES TERMINATED BY '\n'
(description, code, type, initiator, isActionable);

INSERT IGNORE INTO TransactionCodes (createdAt, updatedAt, code, description, type, initiator, isActionable)
SELECT NOW(), NOW(), code, description, type, initiator, isActionable FROM temp_Codes;

DROP TABLE temp_Codes;