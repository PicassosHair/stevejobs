package util

import (
	"encoding/json"
	"bytes"
	"strings"
)

// ExtractApplID gets applId from raw record.
func ExtractApplID(record *RawPatentRecords) string {
	applText := (*record)[0].PatentCaseMetadata.ApplicationNumberText.Value

	return applText
}

// extractTitle gets title text without linebreaks.
func extractTitle(record *RawPatentRecords) string {
  titleContent := (*record)[0].PatentCaseMetadata.InventionTitle.Content
	titleText := ""

	if titleContent != nil {
		titleText = titleContent[0]
	}

	// Remove line breaks
	titleTextProcessed := strings.Replace(titleText, "\n", " ", -1)
	titleTextProcessed = strings.Replace(titleTextProcessed, "^^", " ", -1)
  return titleTextProcessed
}

// ProcessApplication processes generated JSON record and generates a csv-like string for each application.
func ProcessApplication(record *RawPatentRecords) bytes.Buffer {
	var result bytes.Buffer
  metadata := (*record)[0].PatentCaseMetadata

	result.WriteString(ExtractApplID(record))
	result.WriteString("^^")

	// pedsData, _ := json.Marshal((*record)[0].PatentCaseMetadata)

	// result.WriteString(string(pedsData))
	// result.WriteString("^^")

	result.WriteString(metadata.FilingDate)
	result.WriteString("^^")

  result.WriteString(metadata.ApplicationTypeCategory)
	result.WriteString("^^")

  parties := metadata.PartyBag.ApplicantBagOrInventorBagOrOwnerBag
  
  // Follow this order: [examiner, applicant, inventor, practioner, identifier].
  var partyTexts [4]string{"", "", "", "", ""}
  for _, party := range parties {
    // Examiner
    if raw, ok := party["primaryExaminerOrAssistantExaminerOrAuthorizedOfficer"]; ok {
      var partyBag Examiner
      err := json.Unmarshal(*raw, &partyBag)
      if err == nil {
        partyTexts[0] = partyBag[0].Name.PersonNameOrOrganizationNameOrEntityName[0].PersonFullName
      }
    }
    // Applicant
    if raw, ok := party["applicant"]; ok {
      var partyBag Applicant
      err := json.Unmarshal(*raw, &partyBag)
      if err == nil {
        partyTexts[0] = partyBag[0].ContactOrPublicationContact[0].Name.PersonNameOrOrganizationNameOrEntityName[0].PersonStructuredName.LastName
      }
    }
    // Inventor
    if raw, ok := party["inventorOrDeceasedInventor"]; ok {
      var partyBag Inventor
      err := json.Unmarshal(*raw, &partyBag)
      if err == nil {
        partyTexts[0] = partyBag[0].ContactOrPublicationContact[0].Name.PersonNameOrOrganizationNameOrEntityName[0].PersonStructuredName.LastName
      }
    }
    // Practioner
    if raw, ok := party["registeredPractitioner"]; ok {
      var partyBag Practitioner
      err := json.Unmarshal(*raw, &partyBag)
      if err == nil {
        partyTexts[0] = partyBag[0].ContactOrPublicationContact[0].Name.PersonNameOrOrganizationNameOrEntityName[0].PersonStructuredName.LastName
      }
    }
    // Identifier is left as blank for now.
  }

  result.WriteString(partyTexts.join("^^"))
  result.WriteString("^^")

  result.WriteString(metadata.ApplicantFileReference)
  result.WriteString("^^")

	result.WriteString(extractTitle(record))
  
  result.WriteString("\n")
	return result
}

// ProcessCode processes generated JSON record and generate a string of transaction codes. Since the total amount of code is ~600, we will just use a map to dedup here.
func ProcessCode(record *RawPatentRecords, codeMap map[string]bool) bytes.Buffer {
	var result bytes.Buffer
	transactionData := (*record)[0].ProsecutionHistoryDataBag.ProsecutionHistoryData
	for _, event := range transactionData {
		descText := event.EventDescriptionText
		texts := strings.Split(descText, " , ")
		if len(texts) != 2 {
			continue
		}
		code := texts[1]
		if (codeMap)[code] {
			continue
		} else {
			result.WriteString(strings.Join(texts, "^^"))
			result.WriteString("^^")
			result.WriteString("info")
			result.WriteString("^^")
			result.WriteString("uspto")
			result.WriteString("^^")
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
	transactionData := (*record)[0].ProsecutionHistoryDataBag.ProsecutionHistoryData
	applID := ExtractApplID(record)

	for _, event := range transactionData {
		descText := event.EventDescriptionText
		texts := strings.Split(descText, " , ")
		if len(texts) != 2 {
			continue
		}
		result.WriteString(texts[1])
		result.WriteString("^^")
		result.WriteString(applID)
		result.WriteString("^^")
		result.WriteString(event.EventDate)
		result.WriteString("\n")
	}

	return result
}
