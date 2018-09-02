package util

import (
	"bytes"
	"encoding/json"
	"strings"
)

func extractApplID(record *RawPatentRecords) string {
	applText := (*record)[0].PatentCaseMetadata["applicationNumberText"].(map[string]interface{})

	return applText["value"].(string)
}

// ProcessApplication processes generated JSON record and generates a string.
func ProcessApplication(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
	result.WriteString(extractApplID(record))

	result.WriteString("\x00")

	pedsData, _ := json.Marshal((*record)[0].PatentCaseMetadata)

	result.WriteString(string(pedsData))
	result.WriteString("\x00")

	title := (*record)[0].PatentCaseMetadata["inventionTitle"].(map[string]interface{})
	titleText := title["content"].([]interface{})

	// Remove line breaks
	titleTextWithoutBreaks := strings.Replace(titleText[0].(string), "\n", " ", -1)

	result.WriteString(titleTextWithoutBreaks)
	result.WriteString("\x00")

	filingDate := (*record)[0].PatentCaseMetadata["filingDate"].(string)
	result.WriteString(filingDate)
	result.WriteString("\n")

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
			result.WriteString(strings.Join(texts, "\x00"))
			result.WriteString("\x00")
			result.WriteString("info")
			result.WriteString("\x00")
			result.WriteString("uspto")
			result.WriteString("\x00")
			result.WriteString("1")
			result.WriteString("\n")
			(codeMap)[code] = true
		}
	}
	return result
}

// ProcessTransaction processes the record and generates a string of transactions. Separated by linebreaks.
func ProcessTransaction(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
	transactionData := (*record)[0].ProsecutionHistoryDataOrPatentTermData
	applID := extractApplID(record)

	for _, event := range transactionData {
		descText := event.CaseActionDescriptionText
		texts := strings.Split(descText, " , ")
		if len(texts) != 2 {
			continue
		}
		result.WriteString(texts[1])
		result.WriteString("\x00")
		result.WriteString(applID)
		result.WriteString("\x00")
		result.WriteString(event.RecordedDate)
		result.WriteString("\n")
	}

	return result
}
