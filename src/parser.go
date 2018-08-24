package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"util"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadApplicationsToDatabase() {
	db, err := sql.Open("mysql", "root:@/idsguard_dev")
	checkErr(err)

	// Create a new temp table.
	_, err = db.Query(`
	CREATE TABLE IF NOT EXISTS temp_Applications LIKE Applications;
	`)
	checkErr(err)

	// Drop indices to speed up.
	// _, err = db.Query(`
	// DROP INDEX applId_index ON temp_Applications;
	// `)
	// checkErr(err)

	// Load data into temp Applications table.
	_, err = db.Query(`	
	LOAD DATA INFILE '/Users/hao/Desktop/data/out_applications'      
	INTO TABLE temp_Applications
	FIELDS TERMINATED BY '^'
	LINES TERMINATED BY '\n'
	(createdAt, updatedAt, applId, pedsData, title);
	`)
	checkErr(err)

	// Sync temp table with real table.
	_, err = db.Query(`
	INSERT INTO Applications
	SELECT * FROM temp_Applications
	ON DUPLICATE KEY UPDATE updatedAt = VALUES(updatedAt), pedsData = VALUES(pedsData), title = VALUES(title);
	`)

	_, err = db.Query(`
	DROP TABLE temp_Applications;
	`)
	checkErr(err)

}

func main() {
	startTime := time.Now()

	outPath := "/Users/hao/Desktop/data/out"

	readFile, err := os.Open("/Users/hao/Desktop/data/1964.json")

	checkErr(err)

	writeApplicationsFile, err := os.Create(outPath + "_applications")
	writeCodesFile, err := os.Create(outPath + "_codes")

	checkErr(err)

	defer readFile.Close()
	defer writeApplicationsFile.Close()
	defer writeCodesFile.Close()

	decoder := json.NewDecoder(readFile)

	fmt.Println("Parsing JSON...")

	// Used for dedup codes.
	codeSet := map[string]bool{}

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
			}
			writeApplicationsFile.Sync()
		}
	}

	fmt.Println("Loading to Database...")
	loadApplicationsToDatabase()

	fmt.Println("Remove out.txt file...")
	err = os.Remove(outPath + "_applications")
	checkErr(err)

	duration := time.Since(startTime)
	fmt.Printf("Done, used %v seconds.\n", duration.Seconds())
}
