package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

// ProcessApplication processes generated JSON record and generates a string.
func ProcessApplication(record *RawPatentRecords) bytes.Buffer {
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

// ProcessCode processes generated JSON record and generate a string of transaction codes. Since the total amount of code is ~600, we will just use a map to dedup here.
func ProcessCode(record *RawPatentRecords, codeMap map[string]bool) bytes.Buffer {
	var result bytes.Buffer
	transactionData := (*record)[0].ProsecutionHistoryDataOrPatentTermData
	for _, event := range transactionData {
		descText := event.CaseActionDescriptionText
		texts := strings.Split(descText, " , ")
		if len(texts) != 2 {
			continue
		}
		code := texts[1]
		if (codeMap)[code] {
			continue
		} else {
			result.WriteString(strings.Join(texts, "^"))
			result.WriteByte('\n')
			(codeMap)[code] = true
		}
	}
	return result
}
