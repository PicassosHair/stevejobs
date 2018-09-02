-- This template script will read applications file generated by the go package and load to mysql database
DROP TABLE IF EXISTS temp_Applications;
CREATE TABLE IF NOT EXISTS `temp_Applications` (
  `applId` varchar(20) NOT NULL,
  `pedsData` text,
  `title` text,
  `filingDate` varchar(20) NOT NULL
) ENGINE=InnoDB;

LOAD DATA LOCAL INFILE @infile
INTO TABLE temp_Applications
FIELDS TERMINATED BY '^^'
LINES TERMINATED BY '\n'
(applId, pedsData, title, filingDate);

INSERT INTO Applications (createdAt, updatedAt, applId, pedsData, title, filingDate)
SELECT NOW(), NOW(), applId, pedsData, title, filingDate FROM temp_Applications
ON DUPLICATE KEY UPDATE updatedAt = NOW(), pedsData = VALUES(pedsData), title = VALUES(title), filingDate = VALUES(filingDate);

DROP TABLE temp_Applications;