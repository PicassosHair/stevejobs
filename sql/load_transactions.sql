-- This template script will read transactions file generated by the go package and load to mysql database
DROP TABLE IF EXISTS temp_Transactions;
CREATE TABLE IF NOT EXISTS `temp_Transactions` (
  `code` VARCHAR(100) NOT NULL,
  `applId` VARCHAR(10) NOT NULL,
  `recordDate` datetime NOT NULL
) ENGINE=InnoDB;

LOAD DATA LOCAL INFILE @infile
  INTO TABLE temp_Transactions
  FIELDS TERMINATED BY '^^'
  LINES TERMINATED BY '\n'
  (code, applId, recordDate);

-- Temp table, a subset of applications by year, for later join.
-- This can accelerate the process as we only include applications for this year.
DROP TABLE IF EXISTS temp_Applications_ByYear;
CREATE TABLE IF NOT EXISTS `temp_Applications_ByYear` (
    `id` INT(11) NOT NULL,
    `applId` varchar(15) NOT NULL,
    UNIQUE KEY `applId_index` (`applId`)
) ENGINE=InnoDB;

INSERT INTO temp_Applications_ByYear (id, applId)
  SELECT id, applId
  FROM Applications
  WHERE YEAR(filingDate) = @year;

-- Get applicationId.
ALTER TABLE `temp_Transactions` ADD INDEX `applId_index` (`applId`);

DROP TABLE IF EXISTS temp_Transactions_WithAppl;
CREATE TABLE IF NOT EXISTS `temp_Transactions_WithAppl` (
  `code` VARCHAR(100) NOT NULL,
  `applId` VARCHAR(10) NOT NULL,
  `recordDate` datetime NOT NULL,
  `applicationId` varchar(15) NOT NULL
) ENGINE=InnoDB;
INSERT INTO `temp_Transactions_WithAppl` (code, applId, recordDate, applicationId)
  SELECT `code`, tt.`applId`, `recordDate`, tay.id
  FROM temp_Transactions tt
  LEFT JOIN temp_Applications_ByYear tay
  ON tt.applId = tay.applId;

-- Get transactionCodeId
ALTER TABLE `temp_Transactions_WithAppl` ADD INDEX `code_index` (`code`);

DROP TABLE IF EXISTS temp_Transactions_Final;
CREATE TABLE IF NOT EXISTS temp_Transactions_Final (
  `recordDate` datetime NOT NULL,
  `applicationId` varchar(15) NOT NULL,
  `transactionCodeId` varchar(15) NOT NULL
) ENGINE=InnoDB;
INSERT INTO `temp_Transactions_Final` (recordDate, applicationId, transactionCodeId)
  SELECT `recordDate`, tt.applicationId, tc.id
  FROM temp_Transactions_WithAppl tt
  LEFT JOIN TransactionCodes tc
  ON tc.code = tt.code;

-- Final insertion.
INSERT IGNORE INTO Transactions
   (createdAt, updatedAt, transactionCodeId, applicationId, recordDate)
  SELECT NOW(), NOW(), transactionCodeId, applicationId, recordDate
  FROM temp_Transactions_Final;

DROP TABLE IF EXISTS temp_Transactions_Final;
DROP TABLE IF EXISTS temp_Transactions_WithAppl;
DROP TABLE IF EXISTS temp_Applications_ByYear;
DROP TABLE IF EXISTS temp_Transactions;