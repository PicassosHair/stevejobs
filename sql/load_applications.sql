-- This template script will read applications file generated by the go package and load to mysql database
DROP TABLE IF EXISTS temp_Applications;
-- The temp table has no constraints whatsoever, and is a subset of the Applications table.
CREATE TABLE IF NOT EXISTS `temp_Applications` (
  `applId` VARCHAR(30) NOT NULL,
  `filingDate` VARCHAR(20) NOT NULL,
  `typeCategory` VARCHAR(30),
  `examiner` TEXT,
  `applicant` TEXT,
  `inventor` TEXT,
  `practitioner` TEXT,
  `identifier` TEXT,
  `groupArtUnitNumber` VARCHAR(20),
  `confirmationNumber` VARCHAR(20),
  `applicantFileReference` VARCHAR(30),
  `priorityClaim` TEXT,
  `patentClassification` TEXT,
  `businessEntityStatusCategory` VARCHAR(20),
  `firstInventorToFileIndicator` VARCHAR(10),
  `title` text,
  `statusCategory` VARCHAR(120),
  `statusDate` VARCHAR(20),
  `officialFileLocationCategory` VARCHAR(20),
  `relatedDocumentData` TEXT,
  `publicationNumber` VARCHAR(30),
  `publicationDate` VARCHAR(30),
  `patentNumber` VARCHAR(30),
  `grantDate` VARCHAR(30),
  `assignment` TEXT
) ENGINE=InnoDB;

-- Pump the csv-like file to the temp table.
-- The sequence matters.
LOAD DATA LOCAL INFILE @infile
INTO TABLE temp_Applications
FIELDS TERMINATED BY '^^'
LINES TERMINATED BY '\n' (
  applId,
  filingDate,
  typeCategory,
  examiner,
  applicant,
  inventor,
  practitioner,
  identifier,
  groupArtUnitNumber,
  confirmationNumber,
  applicantFileReference,
  priorityClaim,
  patentClassification,
  businessEntityStatusCategory,
  firstInventorToFileIndicator,
  title,
  statusCategory,
  statusDate,
  officialFileLocationCategory,
  relatedDocumentData,
  publicationNumber,
  publicationDate,
  patentNumber,
  grantDate,
  assignment);

INSERT INTO Applications (
  createdAt,
  updatedAt,
  applId,
  filingDate,
  typeCategory,
  examiner,
  applicant,
  inventor,
  practitioner,
  identifier,
  groupArtUnitNumber,
  confirmationNumber,
  applicantFileReference,
  priorityClaim,
  patentClassification,
  businessEntityStatusCategory,
  firstInventorToFileIndicator,
  title,
  statusCategory,
  statusDate,
  officialFileLocationCategory,
  relatedDocumentData,
  publicationNumber,
  publicationDate,
  patentNumber,
  grantDate,
  assignment)
SELECT 
  NOW(),
  NOW(),
  applId,
  CASE 
    WHEN filingDate = '' THEN NULL
    WHEN filingDate IS NULL THEN NULL
    ELSE filingDate
  END AS filingDate,
  typeCategory,
  examiner,
  applicant,
  inventor,
  practitioner,
  identifier,
  groupArtUnitNumber,
  confirmationNumber,
  applicantFileReference,
  priorityClaim,
  patentClassification,
  businessEntityStatusCategory,
  firstInventorToFileIndicator,
  CASE 
    WHEN title IS NULL THEN ''
    ELSE title
  END AS title,
  statusCategory,
  CASE 
    WHEN statusDate = '' THEN NULL
    WHEN statusDate IS NULL THEN NULL
    ELSE statusDate
  END AS statusDate,
  officialFileLocationCategory,
  relatedDocumentData,
  publicationNumber,
  CASE 
    WHEN publicationDate = '' THEN NULL
    WHEN publicationDate IS NULL THEN NULL
    ELSE publicationDate
  END AS publicationDate,
  patentNumber,
  CASE 
    WHEN grantDate = '' THEN NULL
    WHEN grantDate IS NULL THEN NULL
    ELSE grantDate
  END AS grantDate,
  assignment
FROM temp_Applications
ON DUPLICATE KEY UPDATE 
  updatedAt = NOW(),
  title = VALUES(title),
  filingDate = VALUES(filingDate),
  typeCategory = VALUES(typeCategory),
  examiner = VALUES(examiner),
  applicant = VALUES(applicant),
  inventor = VALUES(inventor),
  practitioner = VALUES(practitioner),
  identifier = VALUES(identifier),
  groupArtUnitNumber = VALUES(groupArtUnitNumber),
  confirmationNumber = VALUES(confirmationNumber),
  applicantFileReference = VALUES(applicantFileReference),
  priorityClaim = VALUES(priorityClaim),
  patentClassification = VALUES(patentClassification),
  businessEntityStatusCategory = VALUES(businessEntityStatusCategory),
  firstInventorToFileIndicator = VALUES(firstInventorToFileIndicator),
  title = VALUES(title),
  statusCategory = VALUES(statusCategory),
  statusDate = VALUES(statusDate),
  officialFileLocationCategory = VALUES(officialFileLocationCategory),
  relatedDocumentData = VALUES(relatedDocumentData),
  publicationNumber = VALUES(publicationNumber),
  publicationDate = VALUES(publicationDate),
  patentNumber = VALUES(patentNumber),
  grantDate = VALUES(grantDate),
  assignment = VALUES(assignment);

DROP TABLE temp_Applications;