package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// RawPatentRecords defines raw PEDS bulk data.
type RawPatentRecords []struct {
	PatentCaseMetadata                     map[string]interface{}
	ProsecutionHistoryDataOrPatentTermData []struct {
		RecordDate, CaseActionDescriptionText string
	}
}

func processApplication(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
	currentTime := time.Now()
	formatedTime := currentTime.Format("2006-01-02")
	result.WriteString(formatedTime)
	result.WriteString("^")
	result.WriteString(formatedTime)
	result.WriteString("^")

	applText := (*record)[0].PatentCaseMetadata["applicationNumberText"].(map[string]interface{})

	result.WriteString(applText["value"].(string))
	result.WriteString("^")

	pedsData, _ := json.Marshal((*record)[0].PatentCaseMetadata)

	result.WriteString(string(pedsData))
	result.WriteString("^")

	title := (*record)[0].PatentCaseMetadata["inventionTitle"].(map[string]interface{})
	titleText := title["content"].([]interface{})

	result.WriteString(titleText[0].(string))
	result.WriteByte('\n')

	return result
}

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
	_, err = db.Query(`
	DROP INDEX applId_index ON temp_Applications;
	`)
	checkErr(err)

	// Load data into temp Applications table.
	_, err = db.Query(`	
	LOAD DATA INFILE '/Users/hao/Desktop/data/out.txt'      
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

	outPath := "/Users/hao/Desktop/data/out.txt"

	readFile, err := os.Open("/Users/hao/Desktop/data/2017.json")

	checkErr(err)

	writeFile, err := os.Create(outPath)

	checkErr(err)

	defer readFile.Close()
	defer writeFile.Close()

	decoder := json.NewDecoder(readFile)

	fmt.Println("Parsing JSON...")
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if t == "patentRecord" {
			for decoder.More() {
				var rawRecord RawPatentRecords
				err := decoder.Decode(&rawRecord)
				if err != nil {
					log.Fatal(err)
				}

				applicationText := processApplication(&rawRecord)
				writeFile.WriteString(applicationText.String())
			}
			writeFile.Sync()
		}
	}

	fmt.Println("Loading to Database...")
	loadApplicationsToDatabase()

	fmt.Println("Remove out.txt file...")
	err = os.Remove(outPath)
	checkErr(err)

	duration := time.Since(startTime)
	fmt.Printf("Done, used %v seconds.\n", duration.Seconds())
}
