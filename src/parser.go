package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
	"util"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// func loadApplicationsToDatabase() {
// 	db, err := sql.Open("mysql", "root:@/idsguard_dev")
// 	checkErr(err)

// 	// Create a new temp table.
// 	_, err = db.Query(`
// 	CREATE TABLE IF NOT EXISTS temp_Applications LIKE Applications;
// 	`)
// 	checkErr(err)

// 	// Drop indices to speed up.
// 	// _, err = db.Query(`
// 	// DROP INDEX applId_index ON temp_Applications;
// 	// `)
// 	// checkErr(err)

// 	// Load data into temp Applications table.
// 	_, err = db.Query(`
// 	LOAD DATA INFILE '/Users/hao/Desktop/data/out_applications'
// 	INTO TABLE temp_Applications
// 	FIELDS TERMINATED BY '^'
// 	LINES TERMINATED BY '\n'
// 	(createdAt, updatedAt, applId, pedsData, title);
// 	`)
// 	checkErr(err)

// 	// Sync temp table with real table.
// 	_, err = db.Query(`
// 	INSERT INTO Applications
// 	SELECT * FROM temp_Applications
// 	ON DUPLICATE KEY UPDATE updatedAt = VALUES(updatedAt), pedsData = VALUES(pedsData), title = VALUES(title);
// 	`)

// 	_, err = db.Query(`
// 	DROP TABLE temp_Applications;
// 	`)
// 	checkErr(err)

// }

func main() {

	inFilePath := flag.String("in", "", "The raw json file to parse.")
	outPath := flag.String("out", "", "The raw json file to parse.")

	flag.Parse()

	if *inFilePath == "" || *outPath == "" {
		panic("Error arguments. Use -in and -out.")
	}

	startTime := time.Now()

	readFile, err := os.Open(*inFilePath)

	checkErr(err)

	os.MkdirAll(*outPath, os.ModePerm)
	writeApplicationsFile, err := os.Create(path.Join(*outPath, "applications"))
	writeCodesFile, err := os.Create(path.Join(*outPath, "codes"))
	writeTransactionFile, err := os.Create(path.Join(*outPath, "transactions"))

	checkErr(err)

	defer readFile.Close()
	defer writeApplicationsFile.Close()
	defer writeCodesFile.Close()
	defer writeTransactionFile.Close()

	decoder := json.NewDecoder(readFile)

	fmt.Println("Parsing JSON...")

	// Used for dedup codes.
	codeSet := map[string]bool{}
	count := 0

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if t == "patentRecord" {
			for decoder.More() {
				var rawRecord util.RawPatentRecords
				err := decoder.Decode(&rawRecord)
				if err != nil {
					log.Fatal(err)
				}

				applicationText := util.ProcessApplication(&rawRecord)
				writeApplicationsFile.WriteString(applicationText.String())

				codeText := util.ProcessCode(&rawRecord, codeSet)
				writeCodesFile.WriteString(codeText.String())

				transactionText := util.ProcessTransaction(&rawRecord)
				writeTransactionFile.WriteString(transactionText.String())
			}
			writeApplicationsFile.Sync()
			writeCodesFile.Sync()
			writeTransactionFile.Sync()

			// Log every 100 applications
			count++
			if count%100 == 0 {
				fmt.Println(count, "...")
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("Done, used %v seconds.\n", duration.Seconds())
}
