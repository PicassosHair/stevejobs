package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

func extractApplID(record *RawPatentRecords) string {
	applText := (*record)[0].PatentCaseMetadata["applicationNumberText"].(map[string]interface{})

	return applText["value"].(string)
}

// ProcessApplication processes generated JSON record and generates a string.
func ProcessApplication(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
	currentTime := time.Now()
	formatedTime := currentTime.Format("2006-01-02")
	result.WriteString(formatedTime)
	result.WriteString("^")
	result.WriteString(formatedTime)
	result.WriteString("^")
	result.WriteString(extractApplID(record))

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
	currentTime := time.Now()
	formatedTime := currentTime.Format("2006-01-02")
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
			result.WriteString(formatedTime)
			result.WriteString("^")
			result.WriteString(formatedTime)
			result.WriteString("^")
			result.WriteString(strings.Join(texts, "^"))
			result.WriteString("^")
			result.WriteString("info")
			result.WriteString("^")
			result.WriteString("uspto")
			result.WriteString("^")
			result.WriteString("1")
			result.WriteByte('\n')
			(codeMap)[code] = true
		}
	}
	return result
}

// ProcessTransaction processes the record and generates a string of transactions. Separated by linebreaks.
func ProcessTransaction(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
	currentTime := time.Now()
	formatedTime := currentTime.Format("2006-01-02")
	transactionData := (*record)[0].ProsecutionHistoryDataOrPatentTermData
	applID := extractApplID(record)

	for _, event := range transactionData {
		descText := event.CaseActionDescriptionText
		texts := strings.Split(descText, " , ")
		if len(texts) != 2 {
			continue
		}
		result.WriteString(formatedTime)
		result.WriteString("^")
		result.WriteString(formatedTime)
		result.WriteString("^")
		result.WriteString(texts[1])
		result.WriteString("^")
		result.WriteString(applID)
		result.WriteString("^")
		result.WriteString(event.RecordedDate)
		result.WriteByte('\n')
	}

	return result
}
