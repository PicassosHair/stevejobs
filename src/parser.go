package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// RawPatentRecords defines raw PEDS bulk data.
type RawPatentRecords []struct {
	PatentCaseMetadata struct {
		ApplicationNumberText struct {
			Value, ElectronicText string
		}
	}
}

func main() {
	fmt.Printf("starting...")

	file, _ := os.Open("/Users/hao/Desktop/data/1964.json")
	decoder := json.NewDecoder(file)
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if t == "patentRecord" {
			for decoder.More() {
				var m RawPatentRecords
				err := decoder.Decode(&m)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(m[0].PatentCaseMetadata.ApplicationNumberText.Value)
			}
		}
	}

	file.Close()
}
