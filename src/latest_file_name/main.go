package main

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"strings"
)

type ApiData struct {
	HelpText string `json:"helpText"`
	JsonDownloadMetadata []struct {
		FileName string `json:"fileName"`
		LastUpdated string `json:"lastUpdated"`
		UpdatedFile bool `json:"updatedFile"`
	} `json:"jsonDownloadMetadata`
}

func main() {
	isFull := flag.Bool("isFull", false, "Whether to download the full data.")

	flag.Parse()

	resp, err := http.Get("https://ped.uspto.gov/api/")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
	}
	var data ApiData
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	
	// Get full and delta download URLs.
	deltaUrl := ""
	fullUrl := ""

	for _, metadata := range data.JsonDownloadMetadata {
		if strings.HasPrefix(metadata.FileName, "pairbulk-delta") {
			deltaUrl = metadata.FileName
		}
		if strings.HasPrefix(metadata.FileName, "2000-2019") {
			fullUrl = metadata.FileName
		}
	}

	if *isFull {
		fmt.Println(fullUrl)
	} else {
		fmt.Println(deltaUrl)
	}

	os.Exit(0)
}
