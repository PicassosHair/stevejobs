package util

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

// ExtractApplID gets applId from raw record.
func ExtractApplID(record *RawPatentRecord) string {
	applText := (*record).PatentCaseMetadata.ApplicationNumberText.Value

	return applText
}

func escapeText(tt *string) string {
	r := strings.NewReplacer(
		"\n", " ", // line break.
		"^^", " ", // field break.
		"\\", "", // special chars.
		"|", "", // field line break - each field could be an array.
		"~", "") // atom field break.
	return r.Replace(*tt)
}

// extractTitle gets title text without linebreaks.
func extractTitle(record *RawPatentRecord) string {
	titleContent := (*record).PatentCaseMetadata.InventionTitle.Content
	titleText := ""

	if titleContent != nil {
		titleText = titleContent[0]
	}

	// Remove line breaks
	titleTextProcessed := escapeText(&titleText)
	return titleTextProcessed
}

// extractContact converts the contact array to a plain text with single contact.
func extractContact(contacts *[]Contact) string {
	if len(*contacts) == 0 {
		return ""
	}
	contact := (*contacts)[0]
	var result bytes.Buffer
	name := contact.Name.PersonNameOrOrganizationNameOrEntityName
	hasName := len(name) > 0
	// Full name.
	if hasName {
		result.WriteString(name[0].PersonFullName)
	}
	result.WriteString("|")

	// First name, Middle name, Last name.
	if hasName {
		result.WriteString(name[0].PersonStructuredName.FirstName)
	}
	result.WriteString("|")

	if hasName {
		result.WriteString(name[0].PersonStructuredName.MiddleName)
	}
	result.WriteString("|")

	if hasName {
		result.WriteString(name[0].PersonStructuredName.LastName)
	}
	result.WriteString("|")

	if hasName {
		result.WriteString(name[0].PersonStructuredName.NameSuffix)
	}
	result.WriteString("|")

	// Organization name.
	if hasName {
		hasOrgName := len(name[0].OrganizationStandardName.Content) > 0
		if hasOrgName {
			result.WriteString(name[0].OrganizationStandardName.Content[0])
		}
	}
	result.WriteString("|")

	// Phone number.
	if len(contact.PhoneNumberBag.PhoneNumber) > 0 {
		result.WriteString(contact.PhoneNumberBag.PhoneNumber[0].Value)
	}
	result.WriteString("|")

	// City name.
	result.WriteString(contact.CityName)
	result.WriteString("|")

	// Region.
	result.WriteString(contact.GeographicRegionName.Value)
	result.WriteString("|")

	// Region category.
	result.WriteString(contact.GeographicRegionName.GeographicRegionCategory)
	result.WriteString("|")

	// Country Code.
	result.WriteString(contact.CountryCode)

	return result.String()
}

// ProcessApplication processes generated JSON record and generates a csv-like string for each application. TODO: parse all parties, not just the first one.
func ProcessApplication(record *RawPatentRecord) bytes.Buffer {
	var result bytes.Buffer
	metadata := (*record).PatentCaseMetadata

	result.WriteString(ExtractApplID(record))
	result.WriteString("^^")

	result.WriteString(metadata.FilingDate)
  if metadata.FilingDate != nil && len(metadata.FilingDate) > 0 {
    result.WriteString(" 17:00:00")
  }
	result.WriteString("^^")

	result.WriteString(metadata.ApplicationTypeCategory)
	result.WriteString("^^")

	// Parties.
	parties := metadata.PartyBag.ApplicantBagOrInventorBagOrOwnerBag

	// Follow this order: [examiner, applicant, inventor, practitioner, identifier].
	partyTexts := [5]string{"", "", "", "", ""}
	for _, party := range parties {
		// Examiner
		if raw, ok := party["primaryExaminerOrAssistantExaminerOrAuthorizedOfficer"]; ok {
			var examiner Examiner
			err := json.Unmarshal(*raw, &examiner)
			if err == nil {
				contacts := ([]Contact)(examiner)
				partyTexts[0] = extractContact(&contacts)
			} else {
				log.Fatal("Failed parse primaryExaminerOrAssistantExaminerOrAuthorizedOfficer.")
			}
		}
		// Applicant
		if raw, ok := party["applicant"]; ok {
			var applicant Applicant
			err := json.Unmarshal(*raw, &applicant)
			if err == nil && len(applicant) > 0 {
				contacts := ([]Contact)(applicant[0].ContactOrPublicationContact)
				partyTexts[1] = extractContact(&contacts)
			} else {
				log.Fatal("Failed parse applicant.")
			}
		}
		// Inventor
		if raw, ok := party["inventorOrDeceasedInventor"]; ok {
			var inventor Inventor
			err := json.Unmarshal(*raw, &inventor)
			if err == nil && len(inventor) > 0 {
        var inventorTexts []string
        for _, contactWrapper := range inventor {
          contact := ([]Contact)(contactWrapper.ContactOrPublicationContact)
          inventorTexts = append(inventorTexts, extractContact(&contact))
        }
				partyTexts[2] = strings.Join(inventorTexts, "~")
			} else {
				log.Fatal("Failed parse inventorOrDeceasedInventor.")
			}
		}
		// Practitioner
		if raw, ok := party["registeredPractitioner"]; ok {
			var practitioner Practitioner
			err := json.Unmarshal(*raw, &practitioner)
			if err == nil && len(practitioner) > 0 {
				contacts := ([]Contact)(practitioner[0].ContactOrPublicationContact)
				partyTexts[3] = extractContact(&contacts)
			} else {
				log.Fatal("Failed parse registeredPractitioner.")
			}
		}
		// Identifier is left as blank for now.
	}

	result.WriteString(strings.Join(partyTexts[:], "^^"))
	result.WriteString("^^")
	// End parties.

	result.WriteString(metadata.GroupArtUnitNumber.Value)
	result.WriteString("^^")

	result.WriteString(metadata.ApplicationConfirmationNumber)
	result.WriteString("^^")

	result.WriteString(escapeText(&metadata.ApplicantFileReference))
	result.WriteString("^^")

	// priorityClaimBag
	if len(metadata.PriorityClaimBag.PriorityClaim) > 0 {
		result.WriteString(metadata.PriorityClaimBag.PriorityClaim[0].ApplicationNumber.ApplicationNumberText)
		result.WriteString("~")
		result.WriteString(metadata.PriorityClaimBag.PriorityClaim[0].FilingDate)
    if metadata.PriorityClaimBag.PriorityClaim[0].FilingDate != nil &&
      len(metadata.PriorityClaimBag.PriorityClaim[0].FilingDate) > 0 {
      result.WriteString(" 17:00:00")
    }
		result.WriteString("~")
		result.WriteString(metadata.PriorityClaimBag.PriorityClaim[0].IPOfficeName)
		result.WriteString("~")
		result.WriteString(metadata.PriorityClaimBag.PriorityClaim[0].SequenceNumber)
	}
	result.WriteString("^^")

	// patentClassificationBag
	if len(metadata.PatentClassificationBag.CpcClassificationBagOrIPCClassificationOrECLAClassificationBag) > 0 {
		result.WriteString(metadata.PatentClassificationBag.CpcClassificationBagOrIPCClassificationOrECLAClassificationBag[0].IPOfficeCode)
		result.WriteString("~")
		result.WriteString(metadata.PatentClassificationBag.CpcClassificationBagOrIPCClassificationOrECLAClassificationBag[0].MainNationalClassification.NationalClass)
		result.WriteString("~")
		result.WriteString(metadata.PatentClassificationBag.CpcClassificationBagOrIPCClassificationOrECLAClassificationBag[0].MainNationalClassification.NationalSubclass)
	}
	result.WriteString("^^")

	result.WriteString(metadata.BusinessEntityStatusCategory)
	result.WriteString("^^")

	result.WriteString(metadata.FirstInventorToFileIndicator)
	result.WriteString("^^")

	result.WriteString(extractTitle(record))
	result.WriteString("^^")

	result.WriteString(metadata.ApplicationStatusCategory)
	result.WriteString("^^")

	result.WriteString(metadata.ApplicationStatusDate)
  if metadata.ApplicationStatusDate != nil && len(metadata.ApplicationStatusDate) > 0 {
    result.WriteString(" 17:00:00")
  }
	result.WriteString("^^")

	result.WriteString(metadata.OfficialFileLocationCategory)
	result.WriteString("^^")

	// RelatedDocumentData

	result.WriteString(metadata.PatentPublicationIdentification.PublicationNumber)
	result.WriteString("^^")

  // TODO: create a func and make this less repeatable.
	result.WriteString(metadata.PatentPublicationIdentification.PublicationDate)
  if metadata.PatentPublicationIdentification.PublicationDate != nil &&
    len(metadata.PatentPublicationIdentification.PublicationDate) > 0 {
    result.WriteString(" 17:00:00")
  }
	result.WriteString("^^")

	result.WriteString(metadata.PatentGrantIdentification.PatentNumber)
	result.WriteString("^^")

	result.WriteString(metadata.PatentGrantIdentification.GrantDate)
  if metadata.PatentGrantIdentification.GrantDate != nil &&
    len(metadata.PatentGrantIdentification.GrantDate) > 0 {
    result.WriteString(" 17:00:00")
  }

	result.WriteString("\n")
	return result
}

// ProcessCode processes generated JSON record and generate a string of transaction codes. Since the total amount of code is ~1000, we will just use a map to dedup here.
func ProcessCode(record *RawPatentRecord, codeMap map[string]bool) bytes.Buffer {
	var result bytes.Buffer
	transactionData := (*record).ProsecutionHistoryDataBag.ProsecutionHistoryData
	for _, event := range transactionData {
		desc := event.EventDescriptionText
		code := event.EventCode
		if (codeMap)[code] {
			continue
		} else {
			result.WriteString(desc)
			result.WriteString("^^")
			result.WriteString(code)
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
func ProcessTransaction(record *RawPatentRecord) bytes.Buffer {
	var result bytes.Buffer
	transactionData := (*record).ProsecutionHistoryDataBag.ProsecutionHistoryData
	applID := ExtractApplID(record)

	for _, event := range transactionData {
		eventCode := event.EventCode

		result.WriteString(eventCode)
		result.WriteString("^^")
		result.WriteString(applID)
		result.WriteString("^^")
		result.WriteString(event.EventDate)
    if event.EventDate != nil && len(event.EventDate) > 0 {
      result.WriteString(" 17:00:00")
    }
		result.WriteString("\n")
	}

	return result
}
