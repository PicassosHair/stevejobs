package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
	"util"
)

const (
	loggingThreshold = 1000
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(2)

	inFilePath := flag.String("in", "", "The raw json file to parse.")
	outPath := flag.String("out", "", "The raw json file to parse.")
	debugMode := flag.Bool("debug", false, "Turn on the debug mode, require Application Number as input.")
	applID := flag.String("applId", "", "The application number to locate debug info.")

	flag.Parse()

	if !*debugMode && (*inFilePath == "" || *outPath == "") {
		panic("Error arguments. Use -in and -out.")
	}

	if *debugMode && *applID == "" {
		panic("Debug mode needs a application number, use -applId=12345678")
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

				if !*debugMode {
					applicationText := util.ProcessApplication(&rawRecord)
					writeApplicationsFile.WriteString(applicationText.String())

					codeText := util.ProcessCode(&rawRecord, codeSet)
					writeCodesFile.WriteString(codeText.String())

					transactionText := util.ProcessTransaction(&rawRecord)
					writeTransactionFile.WriteString(transactionText.String())
				} else {
					applIDText := util.ExtractApplID(&rawRecord)
					if *applID == applIDText {
						tcText := util.ProcessTransaction(&rawRecord)
						fmt.Println(tcText.String())
					}
				}
			}
			writeApplicationsFile.Sync()
			writeCodesFile.Sync()
			writeTransactionFile.Sync()

			// Log every 100 applications
			count++
			if count%loggingThreshold == 0 {
				fmt.Println(count, "...")
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("Done, used %v seconds.\n", duration.Seconds())
}
